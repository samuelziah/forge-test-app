# forge-test-app

A small HTTP service in Go, packaged with a signed release tag so you can
verify the published bytes are exactly what I claim.

This is a **test repo** for the Forge platform's publish-and-verify
workflow. The program itself is intentionally small but real: it binds to
a port, serves JSON over HTTP, and has unit tests using `httptest`.

## Endpoints

| Method | Path       | Returns                                    |
|--------|------------|--------------------------------------------|
| GET    | /version   | name, version, Go runtime version, uptime  |
| GET    | /health    | status, total request count                |
| POST   | /echo      | echoes the JSON body you sent              |

## Build

    go build -o forge ./cmd/forge

## Run

    ./forge -addr :8788

In another terminal:

    curl -s http://127.0.0.1:8788/version
    curl -s http://127.0.0.1:8788/health
    curl -s -X POST -H 'Content-Type: application/json' \
        -d '{"hello":"world"}' http://127.0.0.1:8788/echo

## Test

    go test ./...

## Verify the published artifact

The `v1.0` tag is **GPG-signed**.

    git clone https://github.com/samuelziah/forge-test-app.git
    cd forge-test-app
    git rev-parse HEAD
    git verify-tag v1.0
    git fsck --full --strict
    go test ./...

GPG public key fingerprint:

    78AD 6BA6 872D 3685 AA15 0D8A AAA1 6A7B BF1A 8C8B

## Repo layout

    cmd/forge/main.go         # CLI entry point
    internal/forge/app.go     # HTTP handlers
    internal/forge/app_test.go    # httptest-based unit tests
    README.md
    .gitignore

## License

Public domain.
