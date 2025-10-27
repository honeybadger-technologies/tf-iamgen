# GitHub Actions Workflows

This directory contains GitHub Actions workflows that automatically run tests, lint code, and build releases for tf-iamgen.

## Workflows

### 1. **test.yml** - Tests
**Trigger:** On push to `main`/`develop` or pull request

**What it does:**
- ✅ Runs unit tests on **multiple Go versions** (1.21, 1.22)
- ✅ Tests on **multiple OS** (Ubuntu, macOS)
- ✅ Verifies dependencies with `go mod verify`
- ✅ Runs benchmarks for performance tracking
- ✅ Uploads coverage reports to Codecov

**Required for:**
- All pull requests
- Every commit to main/develop

**Status badge:**
```markdown
[![Tests](https://github.com/honeybadger/tf-iamgen/actions/workflows/test.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/test.yml)
```

### 2. **lint.yml** - Code Quality
**Trigger:** On push to `main`/`develop` or pull request

**What it does:**
- ✅ Checks code formatting with `gofmt`
- ✅ Runs `go vet` for common mistakes
- ✅ Runs **golangci-lint** with multiple linters
- ✅ Checks for inefficient assignments
- ✅ Checks for error handling issues
- ✅ Scans for security issues with **Gosec**
- ✅ Uploads security findings to GitHub Security tab

**Linters included:**
- gofmt (formatting)
- go vet (correctness)
- golangci-lint (comprehensive)
- ineffassign (unused assignments)
- errcheck (unhandled errors)
- gosec (security)

**Status badge:**
```markdown
[![Lint](https://github.com/honeybadger/tf-iamgen/actions/workflows/lint.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/lint.yml)
```

### 3. **build.yml** - Build & Release
**Trigger:** On push to `main`, tags `v*`, or pull request to `main`/`develop`

**What it does:**
- ✅ Builds binaries for **multiple platforms:**
  - Linux (x86_64, ARM64)
  - macOS (x86_64, ARM64)
  - Windows (x86_64)
- ✅ Injects version info from git tags
- ✅ Uploads artifacts to GitHub Actions
- ✅ Creates GitHub release with binaries (on tag push)
- ✅ Supports pre-releases for alpha/beta versions
- ✅ Builds Docker image (on tag push)

**Release process:**
1. Tag a commit: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. Workflow automatically creates GitHub release with binaries

**Supported platforms:**
| OS | Architecture | Binary |
|-------|---|---------|
| Linux | x86_64 | tf-iamgen |
| Linux | ARM64 | tf-iamgen |
| macOS | x86_64 | tf-iamgen |
| macOS | ARM64 | tf-iamgen |
| Windows | x86_64 | tf-iamgen.exe |

**Status badge:**
```markdown
[![Build](https://github.com/honeybadger/tf-iamgen/actions/workflows/build.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/build.yml)
```

## Matrix Strategy

### Go Versions Tested:
- **1.21** (minimum supported)
- **1.22** (latest)

### Operating Systems:
- **Ubuntu Latest** (default build environment)
- **macOS Latest** (for macOS compatibility)

### Build Targets:
- Linux amd64, arm64
- macOS amd64, arm64
- Windows amd64

## Required Checks (Branch Protection)

To ensure code quality, we recommend enabling these required checks on `main`:

```
- ✅ Tests / Unit Tests (ubuntu-latest, 1.21)
- ✅ Tests / Unit Tests (ubuntu-latest, 1.22)
- ✅ Tests / Integration Tests
- ✅ Lint / Code Quality
- ✅ Lint / Security Scan
```

## Continuous Integration Setup

### 1. Enable GitHub Actions
Settings → Actions → All workflows and reusable workflows enabled

### 2. Set Branch Protection Rules
Settings → Branches → Add rule for `main`:
- Require status checks to pass before merging
- Select the workflows listed above
- Require branches to be up to date before merging

### 3. Configure Codecov (Optional)
1. Go to https://codecov.io
2. Add repository
3. Copy token (if private repo)
4. Add to GitHub Secrets: `CODECOV_TOKEN`

### 4. Configure Security Scanning
- CodeQL is already enabled in workflows
- GitHub Security tab shows found issues

## Manual Trigger

You can manually trigger workflows from GitHub UI:

1. Go to **Actions** tab
2. Select workflow
3. Click **Run workflow**
4. Choose branch/ref
5. Click **Run workflow**

## Viewing Results

### Test Results
Actions → Tests → Click run to see logs

### Lint Results
Actions → Lint → Inline annotations on PR

### Build Results
Actions → Build → Download artifacts or check release

### Coverage
- View on Codecov: https://codecov.io/gh/honeybadger/tf-iamgen
- Badge: `[![codecov](https://codecov.io/gh/honeybadger/tf-iamgen/branch/main/graph/badge.svg)](https://codecov.io/gh/honeybadger/tf-iamgen)`

## Troubleshooting

### Workflow Failing
1. Check **Logs** in GitHub Actions tab
2. Look for error messages
3. Common issues:
   - Missing Go version (update `.go-version`)
   - Dependencies not installed (update `go.mod`)
   - Lint failures (run `make fmt` and `make lint` locally)
   - Tests failing (run `make test` locally)

### Slow Tests
- Workflows run in parallel across matrix
- Check if tests are waiting for network/files
- Profile with `make benchmark`

### Build Not Creating Release
- Must be on `main` branch
- Must be a tag starting with `v` (e.g., `v1.0.0`)
- Must pass all tests and lint checks

## Best Practices

1. **Always run tests locally before pushing**
   ```bash
   make all  # format, lint, test, build
   ```

2. **Keep dependencies updated**
   ```bash
   go get -u ./...
   go mod tidy
   ```

3. **Write meaningful commit messages**
   ```
   feat: add xyz
   fix: resolve abc
   docs: update readme
   ```

4. **Create pull requests for review**
   - Workflows run automatically
   - Fix any failures before merging

5. **Tag releases properly**
   ```bash
   git tag v1.0.0  # semantic versioning
   git push origin v1.0.0
   ```

## Status Badges

Add to README.md:

```markdown
## CI/CD Status

[![Tests](https://github.com/honeybadger/tf-iamgen/actions/workflows/test.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/test.yml)
[![Lint](https://github.com/honeybadger/tf-iamgen/actions/workflows/lint.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/lint.yml)
[![Build](https://github.com/honeybadger/tf-iamgen/actions/workflows/build.yml/badge.svg)](https://github.com/honeybadger/tf-iamgen/actions/workflows/build.yml)
```

---

**Need help?** Check the individual workflow files or GitHub Actions documentation.
