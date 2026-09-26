# Simple Commerce

Small Go service for comparing a PostgreSQL catalog baseline (`v1`) with a
best-effort Redis cache variant (`v2`). Both versions expose the same response
contract and use the same deterministic 10,000-product seed.

## Scope

- Customer register, login, refresh, logout, and authenticated `GET /api/v1/customer/me`.
- Read-only categories and product list/detail endpoints.
- PostgreSQL is required. Redis is optional: v2 falls back to PostgreSQL when unavailable.
- Cart, seller, checkout, order, transaction, email verification, queues, and telemetry stacks are intentionally out of scope.

## Run

Copy `.env.example` to `.env`, replace every placeholder secret, then run:

```sh
docker compose up --build
```

`db.sql` is a fresh-database initializer. It does not migrate an existing
PostgreSQL volume. Remove/reset local volumes only when their data is disposable.

Useful checks:

```sh
go test ./...
go vet ./...
docker compose config --quiet
k6 run -e BASE_URL=http://localhost:8080 tests/k6/smoke.js
```

The API contract is in `api.yaml`; the cleanup and benchmark sequence is in
`docs/simple-commerce-revamp-plan.md`.
