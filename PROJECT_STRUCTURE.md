# tf-iamgen Project Structure

This document describes the organization and purpose of each directory in the tf-iamgen project.

## Directory Layout

```
tf-iamgen/
├── cmd/                          # CLI Command implementations
│   ├── root.go                   # Root command setup
│   ├── analyze.go                # 'analyze' command: scan Terraform for resources
│   ├── generate.go               # 'generate' command: output IAM policy
│   └── learn.go                  # 'learn' command: CloudTrail integration (Phase 2)
│
├── internal/                     # Private application code
│   ├── parser/
│   │   ├── terraform_parser.go   # HCL parser using hashicorp/hcl
│   │   ├── ast.go               # AST data structures for parsed Terraform
│   │   └── resources.go         # Resource type definitions
│   │
│   ├── mapping/
│   │   ├── aws_mapping.go       # In-memory mapping manager
│   │   ├── loader.go            # Load mappings from YAML
│   │   └── cache.go             # Caching layer for performance
│   │
│   ├── policy/
│   │   ├── generator.go         # Core policy generation logic
│   │   ├── formatter.go         # JSON formatting and output
│   │   └── models.go            # Policy data structures
│   │
│   └── cloudtrail/              # CloudTrail integration (Phase 2+)
│       ├── collector.go         # CloudTrail event collection
│       ├── analyzer.go          # Event analysis and action extraction
│       └── models.go            # CloudTrail data models
│
├── mappings/                     # AWS IAM Mapping Database
│   ├── aws_mapping.yaml         # Main resource → IAM action mappings
│   ├── services/                # Service-specific mappings (optional)
│   │   ├── s3_mapping.yaml
│   │   ├── ec2_mapping.yaml
│   │   ├── rds_mapping.yaml
│   │   └── iam_mapping.yaml
│   └── README.md                # Mapping format documentation
│
├── examples/                     # Example Terraform Projects
│   ├── simple_vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   │
│   ├── simple_s3/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   │
│   └── complex_deployment/      # (Future) multi-service example
│       └── ...
│
├── tests/                        # Test Suite
│   ├── unit/
│   │   ├── parser_test.go       # Parser unit tests
│   │   ├── mapping_test.go      # Mapping unit tests
│   │   ├── generator_test.go    # Policy generator tests
│   │   └── cache_test.go        # Cache tests
│   │
│   └── integration/              # (Future) integration tests
│       ├── e2e_test.go
│       └── fixtures/
│
├── ui/                           # Web Dashboard (Phase 2+)
│   ├── frontend/                # React app
│   │   ├── src/
│   │   ├── public/
│   │   └── package.json
│   │
│   └── backend/                 # FastAPI backend (Phase 2+)
│       ├── app/
│       ├── requirements.txt
│       └── main.py
│
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md          # System design and architecture
│   ├── CONTRIBUTING.md          # Contribution guidelines
│   ├── API.md                   # CLI API documentation
│   ├── MAPPING_FORMAT.md        # How to write resource mappings
│   ├── EXAMPLES.md              # Usage examples
│   ├── PHASES.md                # Roadmap details
│   └── SECURITY.md              # Security considerations
│
├── scripts/                      # Build and Utility Scripts
│   ├── build.sh                 # Build script
│   ├── install-deps.sh          # Dependency installation
│   ├── test.sh                  # Test runner
│   ├── lint.sh                  # Code quality checks
│   └── generate-mappings.sh     # Mapping generation helper
│
├── .github/                      # GitHub Configuration
│   └── workflows/               # CI/CD Workflows
│       ├── test.yml             # Run tests on push
│       ├── build.yml            # Build releases
│       └── lint.yml             # Code quality checks
│
├── build/                        # Build Artifacts (git-ignored)
│   └── bin/                      # Compiled binaries
│
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
├── go.sum                        # Go dependency checksums
├── Makefile                      # Build automation
├── LICENSE                       # MIT License
└── README.md                     # Project overview
```

## Key Directories Explained

### `/cmd` - CLI Commands
Implements Cobra CLI framework commands:
- **root.go**: Main command router and app setup
- **analyze.go**: Scans Terraform directory, identifies resources
- **generate.go**: Outputs IAM policy JSON
- **learn.go**: (Phase 2) Connects to CloudTrail for dynamic analysis

### `/internal` - Core Business Logic
Private packages following Go conventions:

**parser/**: Parses Terraform HCL files into abstract syntax tree
- Uses `hashicorp/hcl/v2` library
- Extracts resource types and configurations
- Handles both `.tf` and `.tf.json` formats

**mapping/**: Manages Terraform → IAM action mappings
- Loads YAML mapping database at startup
- Caches results for performance
- Provides lookup: `aws_s3_bucket` → `[s3:CreateBucket, s3:PutBucketPolicy, ...]`

**policy/**: Generates AWS IAM policy documents
- Aggregates actions discovered by parser
- Groups by AWS service
- Outputs valid IAM policy JSON

**cloudtrail/**: (Phase 2+) Learns from actual AWS API calls
- Connects to AWS CloudTrail
- Analyzes historical events
- Refines predicted permissions

### `/mappings` - IAM Mapping Database
YAML/JSON files defining resource-to-action relationships:
```yaml
aws_s3_bucket:
  create:
    - s3:CreateBucket
  manage:
    - s3:PutBucketPolicy
    - s3:PutBucketAcl
    - s3:PutBucketVersioning
```

**Supports:**
- Service-specific files (`s3_mapping.yaml`, `ec2_mapping.yaml`)
- Hierarchical action grouping (create, manage, delete, etc.)
- Action pre-conditions and resource ARN patterns

### `/examples` - Example Projects
Real Terraform projects for testing and documentation:
- **simple_vpc/**: Basic VPC creation (good for CLI demo)
- **simple_s3/**: S3 bucket setup (test S3 mapping)
- **complex_deployment/** (future): Full multi-service architecture

Each example should work end-to-end: `tf-iamgen analyze → generate → validate`

### `/tests` - Test Suite

**unit/**: Fast unit tests for individual components
- Parser correctness with sample HCL
- Mapping database accuracy
- Policy generation JSON validity
- Cache performance

**integration/**: (Future) End-to-end tests
- Analyze example projects
- Validate output policies
- Test against real AWS (in sandbox)

### `/docs` - Documentation

**ARCHITECTURE.md**: System design, data flows, component interactions

**MAPPING_FORMAT.md**: How contributors add new resource types

**API.md**: CLI usage, flags, output formats

**EXAMPLES.md**: Real-world usage scenarios and workflows

**SECURITY.md**: How data is handled, privacy guarantees, compliance

### `/scripts` - Build Automation

**build.sh**: Cross-platform compilation (Linux, macOS, Windows)

**test.sh**: Runs unit and integration tests

**lint.sh**: Code quality (go fmt, go vet, golangci-lint)

**generate-mappings.sh**: Helper to extract IAM actions from AWS docs

### `/.github/workflows` - CI/CD Pipeline

**test.yml**: Run tests on every push/PR

**build.yml**: Create releases for each tag

**lint.yml**: Enforce code quality standards

## Development Workflow

1. **Parse** Terraform files → Extract AWS resources (cmd/analyze.go → internal/parser)
2. **Map** Resource types → IAM actions (internal/mapping)
3. **Generate** Valid IAM policy JSON (internal/policy)
4. **Output** Results to file or stdout (cmd/generate.go)

## Phase Considerations

### Phase 1 (MVP)
- Focus: cmd/, internal/{parser, mapping, policy}, mappings/, examples/, tests/unit/
- Ignore: ui/, cloudtrail/, Phase 2+ documentation

### Phase 2 (CloudTrail)
- Add: internal/cloudtrail/, ui/backend
- Expand: cmd/learn.go, tests/integration

### Phase 3+ (CI/CD & Enterprise)
- Add: GitHub Actions plugins, API server
- Expand: ui/frontend, API endpoints

## Naming Conventions

- **Files**: lowercase_with_underscores.go
- **Packages**: short, descriptive, no underscores
- **Functions**: PascalCase (exported), camelCase (private)
- **Constants**: ALL_CAPS_WITH_UNDERSCORES

## Adding New Features

1. Create new file in appropriate `/internal` package
2. Add unit tests in `/tests/unit`
3. Update mappings if adding resource support
4. Document in `/docs`
5. Example in `/examples` if user-facing

---

**Happy coding!** 🚀
