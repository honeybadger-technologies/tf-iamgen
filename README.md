# tf-iamgen 🛡️

**Automatically discover and generate AWS IAM permissions required to execute a Terraform project.**

Helps DevOps and SecOps teams achieve least-privilege IAM configurations without guesswork.

[![Build Status](https://github.com/honeybadger-technologies/tf-iamgen/workflows/Build/badge.svg)](https://github.com/honeybadger-technologies/tf-iamgen/actions/workflows/build.yml)
[![Tests Status](https://github.com/honeybadger-technologies/tf-iamgen/workflows/Tests/badge.svg)](https://github.com/honeybadger-technologies/tf-iamgen/actions/workflows/test.yml)
[![Lint Status](https://github.com/honeybadger-technologies/tf-iamgen/workflows/Lint/badge.svg)](https://github.com/honeybadger-technologies/tf-iamgen/actions/workflows/lint.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/honeybadger-technologies/tf-iamgen)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/honeybadger-technologies/tf-iamgen)](https://goreportcard.com/report/github.com/honeybadger-technologies/tf-iamgen)

## 🎯 Quick Start

```bash
# Analyze Terraform project for required IAM actions
tf-iamgen analyze ./terraform --coverage

# Generate least-privilege IAM policy
tf-iamgen generate ./terraform --output policy.json

# View policy in different organizations
tf-iamgen generate ./terraform --group-by service --output policy.json

# (Future) Learn from CloudTrail for runtime analysis
tf-iamgen learn --trail mytrail
```

## ✨ Key Features (Phase 1 - MVP)

- ✅ **Static Analysis**: Parse Terraform files and identify AWS resources
- ✅ **IAM Mapping**: Map Terraform resources to AWS IAM actions via local YAML database
- ✅ **Policy Generation**: Generate least-privilege IAM policy JSON
- ✅ **Least-Privilege by Design**: Only includes necessary permissions
- ✅ **Service Grouping**: Output grouped by AWS service for organization
- ✅ **Coverage Analysis**: Shows which resource types are mapped
- ✅ **Policy Validation**: Detects and warns about overly broad permissions
- ✅ **Provider-Aware**: Extensible architecture for provider operation tracking
- ✅ **Caching**: Fast execution with intelligent caching
- ✅ **Examples**: Bundled example projects to get started

## 📖 Usage Examples

### Analyze Resources

```bash
# Basic analysis
$ tf-iamgen analyze ./terraform

# With IAM mapping coverage report
$ tf-iamgen analyze ./terraform --coverage

# Output
Found 8 resources in ./terraform
Parsed 3 files

Discovered Resources:
  - aws_s3_bucket.bucket (main.tf:5)
  - aws_instance.web (main.tf:15)
  - aws_iam_role.app (iam.tf:3)

========================================================
IAM Mapping Coverage Analysis
========================================================

Resource Type Coverage:
  ✓ aws_s3_bucket
  ✓ aws_instance
  ✓ aws_iam_role

Coverage: 3/3 resource types mapped

Generated Policy Preview:
  Total Actions: 24
  Services: [ec2 iam s3]
  Statements: 3
```

### Generate Policy

```bash
# Generate to stdout (default)
$ tf-iamgen generate ./terraform

# Generate to file
$ tf-iamgen generate ./terraform --output policy.json

# Generate grouped by service
$ tf-iamgen generate ./terraform --group-by service --output policy.json

# Output (example)
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3Permissions",
      "Effect": "Allow",
      "Action": [
        "s3:CreateBucket",
        "s3:DeleteBucket",
        "s3:GetObject",
        "s3:ListBucket"
      ],
      "Resource": "*"
    },
    {
      "Sid": "Ec2Permissions",
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:RunInstances",
        "ec2:TerminateInstances"
      ],
      "Resource": "*"
    }
  ]
}
```

## 🏗️ Project Structure

```
tf-iamgen/
├── cmd/                    # CLI commands (analyze, generate, version)
├── internal/               # Core business logic
│   ├── parser/            # Terraform HCL parser
│   ├── mapping/           # Resource-to-IAM action mappings
│   ├── policy/            # Policy generation logic
│   └── cloudtrail/        # CloudTrail integration (Phase 2+)
├── mappings/              # YAML/JSON mapping database (34 resources)
├── examples/              # Example Terraform projects
├── tests/                 # Unit and integration tests (67+ tests)
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
└── main.go                # Application entry point
```

## 🚀 Roadmap

| Phase | Features | Status |
|-------|----------|--------|
| **1: MVP** | Static analysis, policy generation, CLI | ✅ Complete |
| **2: CloudTrail Learning** | Dynamic analysis, least-privilege refinement, dashboard | 📋 Planned |
| **3: CI/CD Integration** | GitHub Actions, GitLab CI, access analyzer | 📋 Planned |
| **4: Enterprise** | Multi-account, SaaS portal, ML optimization | 📋 Planned |

## 📊 What's Included

- **Terraform Parser**: Full HCL & JSON parsing support
- **34 Resource Types**: Pre-mapped to IAM actions (S3, EC2, IAM, RDS, Lambda, etc.)
- **67+ Unit Tests**: Comprehensive test coverage
- **Policy Generator**: Generates least-privilege policies
- **CLI Tools**: Easy-to-use command-line interface
- **Provider Specs**: Foundation for provider-aware permission analysis

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
