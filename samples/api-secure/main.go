package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiResponse[T any] struct {
	Data  *T        `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
}

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
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
	jwtSecret  = []byte(getEnv("JWT_SECRET", "dev-secret"))
	reqCounter = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "api_requests_total", Help: "total api requests"}, []string{"path", "method", "code"})
)

func getEnv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

func main() {
	prometheus.MustRegister(reqCounter)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/items", withSecurity(handleItems))
	mux.HandleFunc("/items/", withSecurity(handleItemByID))
	mux.Handle("/metrics", promhttp.Handler())

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
		// request id
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = generateID()
		}
		w.Header().Set("X-Request-ID", rid)
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
		// limit body size (1MB)
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, nil, &apiError{Code: "internal_error", Message: "internal server error"})
			}
			log.Printf("rid=%s %s %s %dms", rid, r.Method, r.URL.Path, time.Since(start).Milliseconds())
		}()
		sw := &statusWriter{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(sw, r)
		reqCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(sw.code)).Inc()
	})
}

func withSecurity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorized := false
		// Prefer JWT if provided
		if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			tok := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if t, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("alg")
				}
				return jwtSecret, nil
			}); err == nil && t.Valid {
				authorized = true
			}
		}
		if !authorized {
			// Fallback API key
			key := r.Header.Get("X-API-Key")
			if key == apiKey {
				authorized = true
			}
		}
		if !authorized {
			writeJSON(w, http.StatusUnauthorized, nil, &apiError{Code: "unauthorized", Message: "invalid credentials"})
			return
		}
		// naive rate limit per key (token bucket size 5 per minute)
		k := r.Header.Get("X-API-Key")
		if k == "" {
			k = "jwt"
		}
		if !allow(k) {
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
		// ETag support
		b, _ := json.Marshal(it)
		etag := computeETag(b)
		w.Header().Set("ETag", etag)
		if inm := r.Header.Get("If-None-Match"); inm == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
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

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}

func computeETag(b []byte) string {
	h := sha256.Sum256(b)
	return `W/"` + hex.EncodeToString(h[:]) + `"`
}
