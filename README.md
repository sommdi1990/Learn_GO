## Learn_GO

نمونه‌کدها و مستندات آموزشی زبان Go بر اساس پروژه‌های واقعی موجود در مسیر `/home/soroush/go`.

### ویکی (مستندات)

- فهرست کامل مباحث در `Learn_GO.wiki/Home.md`
- موضوعات: ساختار پروژه، ماژول‌ها، کانفیگ، لاگینگ، خطاها، HTTP/API، دیتابیس و مهاجرت‌ها، همزمانی،
  Worker/Scheduler/RateLimiter، تست و ابزارها، Docker/Compose، SMTP/DMARC، CoreDNS، Kafka/Redpanda، Redis، ClickHouse،
  Git hooks، طراحی API و امنیت.

### پیش‌نیازها

- Go 1.22+
- Docker و Docker Compose (برای برخی نمونه‌ها مانند Redis/ClickHouse)
- ابزار اختیاری: `golangci-lint`

### نمونه‌کدها (Samples)

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

برای توضیحات بیشتر هر نمونه، به صفحهٔ «نمونه‌کدها» در ویکی مراجعه کنید: `Learn_GO.wiki/99-نمونه-کدها.md`.
