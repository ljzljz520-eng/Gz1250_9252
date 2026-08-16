# Studio People Service

Deterministic Go HTTP backend for administrators maintaining photographers, photo editors, makeup artists, and customer-service staff. It loads `fixtures/people.yaml` and stores records in memory.

## Standard commands

```bash
GOTOOLCHAIN=local CGO_ENABLED=0 go build ./...
GOTOOLCHAIN=local go test -count=1 ./...
```

## Run

```bash
GOTOOLCHAIN=local CGO_ENABLED=0 go run ./cmd/studio -addr :8080 -fixture fixtures/people.yaml
```

Use `GET /health` for readiness. Staff are listed with `GET /admin/people`, created with `POST /admin/people`, updated with `PATCH /admin/people/{id}`, and removed with `DELETE /admin/people/{id}`. Staff responses include phone, email, role, status, and `created_at`; failures return a stable `code` field.

## Persistence configuration

`people.Service` accepts the `people.Repository` interface, and `cmd/studio/main.go` currently wires `people.MemoryRepository`. To use PostgreSQL or another durable store, implement that interface and replace the repository construction in the command. Pass the DSN, connection-pool limits, and migration options into that implementation through explicit environment variables or command flags; no database setting is needed for the in-memory default. The HTTP handler and service API remain unchanged.

## Docker validation

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh studiopeople linux/arm64
./build_benzhi_docker.sh studiopeople linux/amd64
docker run -it studiopeople:latest
```

## Known initial failures

See `BUG_REPRO.md` for the exact command and output captured during packaging.
