# Contributing to tf-iamgen

Thank you for your interest in contributing to tf-iamgen! 🙌

This document provides guidelines and instructions for contributing to the project.

## 🎯 Code of Conduct

Be respectful, inclusive, and collaborative. We're building something great together!

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make

### Development Setup

```bash
# Clone the repository
git clone https://github.com/honeybadger/tf-iamgen.git
cd tf-iamgen

# Install dependencies
make install-deps

# Run tests
make test

# Build the application
make build
```

## 📋 Types of Contributions

### 🐛 Bug Reports

Found a bug? Create an issue with:
- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS

### ✨ Feature Requests

Have an idea? Open an issue with:
- Use case and value
- Proposed solution (if any)
- Examples or mockups

### 📝 Adding Resource Mappings

See [MAPPING_FORMAT.md](./MAPPING_FORMAT.md) for:
- YAML structure guidelines
- How to discover required IAM actions
- Validation and testing

### 💻 Code Changes

Follow the workflow below.

## 🔄 Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/description
# or
git checkout -b fix/issue-number
```

### 2. Make Changes

- Write clean, idiomatic Go code
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Add tests for new functionality
- Update documentation as needed

### 3. Run Quality Checks

```bash
# Format code
make fmt

# Run linters
make lint

# Run tests
make test
```

### 4. Commit with Clear Messages

```bash
git commit -m "feat: add S3 bucket encryption mapping"
# or
git commit -m "fix: handle missing resource type gracefully"
```

**Commit message format:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `refactor:` Code refactoring
- `test:` Test additions
- `chore:` Build, CI, dependencies

### 5. Push and Create Pull Request

```bash
git push origin feature/description
```

Create a PR on GitHub with:
- Clear description of changes
- Reference to related issues (#123)
- Screenshot/output if UI-related
- Test evidence for bug fixes

## 📐 Code Style

### Go Conventions

```go
// Package documentation (required)
// Package parser provides Terraform HCL parsing functionality.
package parser

// Exported function (capitalized, with docs)
// ParseFile parses a Terraform HCL file and returns discovered resources.
func ParseFile(path string) ([]Resource, error) {
    // Implementation
}

// Private function (lowercase, no export)
func parseResource(node *hclsyntax.Block) (*Resource, error) {
    // Implementation
}

// Constant grouping
const (
    DefaultCacheSize = 1000
    MaxRetries       = 3
)

// Variable grouping
var (
    ErrInvalidPath = errors.New("invalid file path")
    ErrParseError  = errors.New("failed to parse HCL")
)
```

### File Organization

1. Package documentation
2. Imports
3. Constants
4. Variables
5. Types/Structs
6. Exported functions
7. Private functions

## 🧪 Testing

### Write Tests For

- New public functions
- Bug fixes (regression test)
- Complex logic paths

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := FunctionUnderTest(input)
    
    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Run Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# With coverage
go test -cover ./...

# Specific package
go test -v ./internal/parser
```

## 📚 Documentation

### Update Documentation When

- Adding new commands or flags
- Changing existing behavior
- Adding new resource types
- Significant feature additions

### Documentation Files

- `README.md`: Quick start and overview
- `docs/ARCHITECTURE.md`: System design
- `docs/MAPPING_FORMAT.md`: Resource mapping guidelines
- `docs/API.md`: CLI commands and flags
- `docs/SECURITY.md`: Security considerations

## 🏗️ Project Structure

Respect the established structure:

```
cmd/              - CLI commands only
internal/         - Core business logic (private)
mappings/         - YAML mapping database
examples/         - Example Terraform projects
tests/            - Test files
docs/             - Documentation
scripts/          - Build scripts
```

## ✅ Review Checklist

Before submitting a PR, ensure:

- [ ] Code follows Go conventions
- [ ] All tests pass (`make test`)
- [ ] Linters pass (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] No unrelated changes included

## 🎓 Learning Resources

### Go
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Terraform & IAM
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS IAM Documentation](https://docs.aws.amazon.com/iam/)

### HCL
- [HCL2 Spec](https://github.com/hashicorp/hcl/blob/v2main/hclsyntax/spec.md)

## 🤝 Getting Help

- **Questions?** Open a discussion on GitHub
- **Stuck?** Comment on an issue or PR
- **Ideas?** Start a discussion or issue

## 📝 License

By contributing, you agree your code is licensed under the MIT License.

---

**Thank you for contributing to make infrastructure security better!** ✨
