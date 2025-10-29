# Build Instructions & Troubleshooting

## Quick Start

### Building

```bash
# Using make
make build

# Or direct Go build
go build -o build/bin/tf-iamgen .

# With version information
go build -ldflags "-X main.Version=v0.1.0" -o build/bin/tf-iamgen .
```

### Running

```bash
# Run locally
./build/bin/tf-iamgen version

# Or with make
make analyze
make generate
```

## macOS Specific Issues

### Problem: `dyld: missing LC_UUID load command`

**Symptoms:**
```
dyld[18535]: missing LC_UUID load command in /path/to/tf-iamgen
dyld[18535]: missing LC_UUID load command
```

**Root Cause:**
- Go's CGO is enabled by default on macOS
- On some macOS systems/Go versions, CGO-enabled binaries have LC_UUID issues
- This is not a code issue, but an environment configuration issue

**Solution:**
Disable CGO for static binary build:

```bash
# Method 1: Using make (recommended)
make build

# Method 2: Direct build
CGO_ENABLED=0 go build -o build/bin/tf-iamgen .
```

The Makefile already includes `CGO_ENABLED=0` in the build target.

**Why This Works:**
- Creates a pure Go static binary
- Removes dependency on C libraries
- Avoids macOS dyld issues
- No performance impact for CLI tools

### Verification

After building, verify the binary works:

```bash
./build/bin/tf-iamgen version
# Should output:
# tf-iamgen version dev
# Built at: unknown
# Git commit: unknown
```

## Build Targets

```bash
make help          # Show all available targets
make build         # Build the application
make clean         # Remove build artifacts
make test          # Run all tests
make lint          # Run code quality checks
make install-deps  # Install dependencies
make analyze       # Run analyze on current dir
make generate      # Generate IAM policy
make all           # Run: clean, install-deps, lint, test, build
```

## Development

### Running Tests

```bash
# All tests
make test

# Specific package
go test -v ./internal/policy/

# With coverage
go test -cover ./...

# Benchmark
go test -bench=. -benchmem ./internal/policy/
```

### Code Formatting

```bash
# Format all code
make fmt

# Or directly
gofmt -s -w .

# Check without modifying
gofmt -s -l .
```

### Linting

```bash
# Run linting
make lint

# Or directly
go fmt ./...
go vet ./...
golangci-lint run ./...  # if installed
```

## Troubleshooting

### Build Fails with "no Go files"

**Solution:**
```bash
cd /Users/your/path/to/tf-iamgen
make clean
make build
```

### Tests Fail

**Check Go version:**
```bash
go version
# Should be 1.21 or higher
```

**Clean and rebuild:**
```bash
go clean -cache -testcache
make test
```

### Import errors

**Ensure dependencies are installed:**
```bash
go mod tidy
go mod download
go mod verify
make test
```

## Cross-Platform Builds

For building on different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o tf-iamgen-linux ...

# Windows
GOOS=windows GOARCH=amd64 go build -o tf-iamgen.exe ...

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o tf-iamgen-darwin-amd64 ...

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o tf-iamgen-darwin-arm64 ...
```

**Always use `CGO_ENABLED=0` for CLI distributions:**

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tf-iamgen-linux ...
```

## Docker Build

```bash
# Build Docker image
docker build -t honeybadger/tf-iamgen .

# Run in Docker
docker run honeybadger/tf-iamgen analyze /tf

# Or with volume mount
docker run -v /path/to/terraform:/tf honeybadger/tf-iamgen generate /tf --output policy.json
```

## Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `GO` | Go command | `go` |
| `GOFLAGS` | Go build flags | `-v` |
| `CGO_ENABLED` | Enable C interop | Set to `0` in Makefile |
| `GOOS` | Target OS | Current OS |
| `GOARCH` | Target architecture | Current arch |

## Performance

Binary size after build:
- With CGO: ~15-20 MB (varies)
- Without CGO: ~12 MB

Memory usage:
- Minimal - typically < 50 MB
- Caching improves performance for large projects

## Release Build

For production releases:

```bash
make clean
CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=v0.1.0" -o build/bin/tf-iamgen .
# Strip binary for smaller size
strip build/bin/tf-iamgen
```

## Support

For build issues, check:

1. Go version: `go version` (1.21+)
2. Dependencies: `go mod verify`
3. Formatting: `gofmt -s -l .`
4. Clean rebuild: `make clean build`

For macOS LC_UUID issues specifically, ensure:
- `CGO_ENABLED=0` is used
- Clean build directory: `rm -rf build/bin`
- Fresh module cache: `go clean -cache`
