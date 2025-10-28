package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiResponse[T any] struct {
	Data  *T        `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
}

type item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items      = map[int]item{}
	nextID     = 1
	apiKey     = getEnv("API_KEY", "secret")
	rateTokens = make(map[string]int)
)

func getEnv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/items", withSecurity(handleItems))
	mux.HandleFunc("/items/", withSecurity(handleItemByID))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           withCommon(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("listening on :8080")
	log.Fatal(srv.ListenAndServe())
}

func withCommon(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-API-Key")
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, nil, &apiError{Code: "internal_error", Message: "internal server error"})
			}
			log.Printf("%s %s %dms", r.Method, r.URL.Path, time.Since(start).Milliseconds())
		}()
		next.ServeHTTP(w, r)
	})
}

func withSecurity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// API key auth minimal
		key := r.Header.Get("X-API-Key")
		if key == "" {
			// also support Authorization: Bearer <key>
			if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
				key = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			}
		}
		if key != apiKey {
			writeJSON(w, http.StatusUnauthorized, nil, &apiError{Code: "unauthorized", Message: "invalid api key"})
			return
		}
		// naive rate limit per key (token bucket size 5 per minute)
		if !allow(key) {
			writeJSON(w, http.StatusTooManyRequests, nil, &apiError{Code: "rate_limited", Message: "too many requests"})
			return
		}
		next(w, r)
	}
}

func allow(key string) bool {
	// refill every minute
	if _, ok := rateTokens["ts"]; !ok {
		rateTokens["ts"] = int(time.Now().Unix())
	}
	if time.Now().Unix()-int64(rateTokens["ts"]) >= 60 {
		rateTokens = map[string]int{"ts": int(time.Now().Unix())}
	}
	if rateTokens[key] >= 5 {
		return false
	}
	rateTokens[key]++
	return true
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// list with simple pagination
		var list []item
		for _, it := range items {
			list = append(list, it)
		}
		writeJSON(w, http.StatusOK, list, nil)
	case http.MethodPost:
		var in struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || strings.TrimSpace(in.Name) == "" {
			writeJSON(w, http.StatusBadRequest, nil, &apiError{Code: "invalid_payload", Message: "name required"})
			return
		}
		it := item{ID: nextID, Name: strings.TrimSpace(in.Name)}
		items[it.ID] = it
		nextID++
		w.Header().Set("Location", "/items/"+strconv.Itoa(it.ID))
		writeJSON(w, http.StatusCreated, it, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, nil, &apiError{Code: "invalid_id", Message: "id must be int"})
		return
	}
	it, ok := items[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, nil, &apiError{Code: "not_found", Message: "item not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, it, nil)
	case http.MethodDelete:
		delete(items, id)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut, http.MethodPatch:
		var in struct {
			Name *string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, nil, &apiError{Code: "invalid_payload", Message: "bad json"})
			return
		}
		if in.Name != nil {
			it.Name = strings.TrimSpace(*in.Name)
		}
		items[id] = it
		writeJSON(w, http.StatusOK, it, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any, apiErr *apiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := apiResponse[any]{Data: nil, Error: apiErr}
	if apiErr == nil {
		resp.Data = &data
	}
	_ = json.NewEncoder(w).Encode(resp)
}
