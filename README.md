# tf-iamgen 🛡️

**Automatically discover and generate AWS IAM permissions required to execute a Terraform project.**

Helps DevOps and SecOps teams achieve least-privilege IAM configurations without guesswork.

## 🎯 Quick Start

```bash
# Analyze Terraform project for required IAM actions
tf-iamgen analyze ./terraform

# Generate IAM policy
tf-iamgen generate --output policy.json

# (Future) Learn from CloudTrail
tf-iamgen learn --trail mytrail
```

## ✨ Key Features (Phase 1 - MVP)

- ✅ **Static Analysis**: Parse Terraform files and identify AWS resources
- ✅ **IAM Mapping**: Map Terraform resources to AWS IAM actions via local YAML database
- ✅ **Policy Generation**: Generate least-privilege IAM policy JSON
- ✅ **Service Grouping**: Output grouped by AWS service
- ✅ **Caching**: Fast execution with intelligent caching
- ✅ **Examples**: Bundled example projects to get started

## 🏗️ Project Structure

```
tf-iamgen/
├── cmd/                    # CLI commands
├── internal/               # Core business logic
│   ├── parser/            # Terraform HCL parser
│   ├── mapping/           # Resource-to-IAM action mappings
│   ├── policy/            # Policy generation logic
│   └── cloudtrail/        # CloudTrail integration (Phase 2+)
├── mappings/              # YAML/JSON mapping database
├── examples/              # Example Terraform projects
├── tests/                 # Unit and integration tests
├── ui/                    # Web UI (future)
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
└── main.go                # Application entry point
```

## 🚀 Roadmap

| Phase | Features | Status |
|-------|----------|--------|
| **1: MVP** | Static analysis, policy generation, CLI | 🔨 In Progress |
| **2: CloudTrail Learning** | Dynamic analysis, least-privilege refinement, dashboard | 📋 Planned |
| **3: CI/CD Integration** | GitHub Actions, GitLab CI, access analyzer | 📋 Planned |
| **4: Enterprise** | Multi-account, SaaS portal, ML optimization | 📋 Planned |

## 📋 Requirements

- Go 1.21+
- Terraform 1.0+
- AWS credentials (for Phase 2+)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./docs/CONTRIBUTING.md) for guidelines.

## 📄 License

MIT License - See [LICENSE](./LICENSE)

---

**Making infrastructure secure by default.** 🚀
