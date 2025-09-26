# VaultOne Public Edge API (Local Dev)

## Quickstart

- Requirements: Go 1.24+, Docker Desktop (optional for compose), curl
- Start API (no auth):
  - `make dev`
  - `curl http://localhost:8080/healthz`
- Start with local auth:
  - `LOCAL_AUTH_ENABLED=true LOCAL_API_KEY=dev-please-change make dev`
  - Use header `X-API-Key: dev-please-change`

## Endpoints (MVP)
- `POST /v1/accounts`
- `GET /v1/accounts/{id}/balance`
- `POST /v1/transfers`
- `GET /v1/transfers/{id}`
- `POST /v1/webhooks/payment_event`

See `docs/openapi.yaml` for schemas.

## Docker Compose
- `make up` (requires Docker running) brings up API + Postgres + Redpanda.

## Idempotency
- Provide `Idempotency-Key` on POSTs; identical payload → cached response; different payload → 409.

## Telemetry
- OpenTelemetry stdout exporter enabled; spans printed to terminal.

## Development
- `make tidy` — tidy modules
- `make test` — run tests
- `make build` — build binary
- `make docker-build` — build container image
- `make gen` — generate gRPC stubs (requires buf installed)

## Demo
- `scripts/demo.sh` executes a basic E2E flow (local, stubbed).

## Notes
- No external IdP; optional local API key.
