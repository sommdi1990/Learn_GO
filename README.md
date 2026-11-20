# Go examples for protocols, handlers and concurrency

Independent Go programs isolate topics that are hard to see inside a larger service. Examples include assembling DNS messages, handling HTTP requests, working with channels and controlling worker throughput.

The DNS exercise splits query construction from response parsing. The HTTP test example puts the handler beside its test, making it a small entry point for reviewing request and response behavior without following the database-backed samples.

- [query_pack.go](samples/dnsmessage/query_pack.go)
- [response_parse.go](samples/dnsmessage/response_parse.go)
- [handler.go](samples/httptest-demo/handler.go)
- [handler_test.go](samples/httptest-demo/handler_test.go)

## نمونه‌های مستقل

1) API امن و CRUD: `samples/api-secure`

```bash
cd samples/api-secure
API_KEY=secret go run .
# نمونه درخواست‌ها
curl -H 'X-API-Key: secret' http://localhost:8080/health
curl -H 'X-API-Key: secret' -H 'Content-Type: application/json' -d '{"name":"A"}' http://localhost:8080/items
curl -H 'X-API-Key: secret' http://localhost:8080/items/1
curl -X DELETE -H 'X-API-Key: secret' http://localhost:8080/items/1 -i
```

2) ClickHouse: `samples/clickhouse-demo`

```bash
cd samples/clickhouse-demo
docker compose up -d
go run .
```

3) Redis: `samples/redis-demo`

```bash
cd samples/redis-demo
docker compose up -d
go run .
```

4) وب‌سرور ساده: `samples/http-basic`
```bash
cd samples/http-basic && go run .
```

5) کانفیگ از ENV: `samples/config-env`

```bash
cd samples/config-env && go run .
```

6) خاتمهٔ ایمن: `samples/graceful-shutdown`

```bash
cd samples/graceful-shutdown && go run .
```

7) محدودسازی نرخ ساده: `samples/worker-ratelimiter`

```bash
cd samples/worker-ratelimiter && go run .
```

8) لینت: `samples/lint-demo`

```bash
cd samples/lint-demo
make lint
```

9) Git hooks (pre-commit): `samples/git-hooks-demo`

```bash
cd samples/git-hooks-demo
make install-hooks
git add . && git commit -m "test hooks"
```
