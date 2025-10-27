# tf-iamgen Folder Structure Reference

## 🗂️ Complete Directory Tree

```
tf-iamgen/
│
├─ 📄 main.go                     ← Application entry point
├─ 📄 go.mod                      ← Go module dependencies
├─ 📄 Makefile                    ← Build automation commands
├─ 📄 README.md                   ← Project overview & quick start
├─ 📄 LICENSE                     ← MIT License
│
├─ 📁 cmd/                        ← CLI Command implementations (Phase 1-3)
│  ├─ root.go                    ← Root command setup (Cobra)
│  ├─ analyze.go                 ← "tf-iamgen analyze" command
│  ├─ generate.go                ← "tf-iamgen generate" command
│  └─ learn.go                   ← "tf-iamgen learn" command (Phase 2)
│
├─ 📁 internal/                  ← Private business logic (Go convention)
│  ├─ parser/                    ← Terraform HCL parsing
│  │  ├─ terraform_parser.go    ← Main parser implementation
│  │  ├─ ast.go                 ← AST data structures
│  │  └─ resources.go           ← Resource definitions
│  │
│  ├─ mapping/                   ← Resource-to-IAM action mapping
│  │  ├─ aws_mapping.go         ← Mapping manager
│  │  ├─ loader.go              ← Load YAML mappings
│  │  └─ cache.go               ← LRU cache for lookups
│  │
│  ├─ policy/                    ← IAM policy generation
│  │  ├─ generator.go           ← Policy generation logic
│  │  ├─ formatter.go           ← JSON/text output formatting
│  │  └─ models.go              ← Policy data structures
│  │
│  └─ cloudtrail/               ← CloudTrail integration (Phase 2+)
│     ├─ collector.go           ← Event collection
│     ├─ analyzer.go            ← Event analysis
│     └─ models.go              ← CloudTrail data models
│
├─ 📁 mappings/                 ← AWS IAM Mapping Database
│  ├─ aws_mapping.yaml          ← Main resource → IAM action mappings
│  ├─ README.md                 ← Mapping format documentation
│  └─ services/                 ← (Optional) Service-specific mappings
│     ├─ s3_mapping.yaml       ← S3-specific mappings
│     ├─ ec2_mapping.yaml      ← EC2-specific mappings
│     ├─ rds_mapping.yaml      ← RDS-specific mappings
│     └─ iam_mapping.yaml      ← IAM-specific mappings
│
├─ 📁 examples/                 ← Example Terraform projects
│  ├─ simple_vpc/               ← Basic VPC example
│  │  ├─ main.tf               ← VPC resource definitions
│  │  ├─ variables.tf          ← Variable declarations
│  │  └─ outputs.tf            ← Output declarations
│  │
│  ├─ simple_s3/                ← S3 bucket example
│  │  ├─ main.tf               ← S3 resource definitions
│  │  ├─ variables.tf          ← Variable declarations
│  │  └─ outputs.tf            ← Output declarations
│  │
│  └─ complex_deployment/       ← (Future) Multi-service architecture
│     └─ ...                    ← Various resource types
│
├─ 📁 tests/                    ← Test suite
│  ├─ unit/                     ← Unit tests (fast, isolated)
│  │  ├─ parser_test.go        ← Parser unit tests
│  │  ├─ mapping_test.go       ← Mapping unit tests
│  │  ├─ generator_test.go     ← Policy generator tests
│  │  └─ cache_test.go         ← Cache tests
│  │
│  └─ integration/              ← Integration tests (slow, realistic)
│     ├─ e2e_test.go           ← End-to-end tests
│     └─ fixtures/             ← Test data and fixtures
│
├─ 📁 docs/                     ← Project documentation
│  ├─ ARCHITECTURE.md          ← System design & data flows
│  ├─ CONTRIBUTING.md          ← Developer guidelines
│  ├─ STRUCTURE_GUIDE.md       ← This file
│  ├─ API.md                   ← CLI commands & flags (future)
│  ├─ MAPPING_FORMAT.md        ← How to write resource mappings (future)
│  ├─ EXAMPLES.md              ← Usage examples & workflows (future)
│  ├─ PHASES.md                ← Detailed roadmap (future)
│  └─ SECURITY.md              ← Security & privacy (future)
│
├─ 📁 scripts/                  ← Build & utility scripts
│  ├─ build.sh                 ← Cross-platform build
│  ├─ install-deps.sh          ← Dependency installation
│  ├─ test.sh                  ← Test runner
│  ├─ lint.sh                  ← Code quality checks
│  └─ generate-mappings.sh     ← Mapping generation helper
│
├─ 📁 ui/                       ← Web Dashboard (Phase 2+)
│  ├─ frontend/                ← React dashboard
│  │  ├─ src/                 ← React components & logic
│  │  ├─ public/              ← Static assets
│  │  ├─ package.json         ← NPM dependencies
│  │  └─ ...
│  │
│  └─ backend/                 ← FastAPI backend (Phase 2+)
│     ├─ app/                 ← API endpoints
│     ├─ requirements.txt      ← Python dependencies
│     ├─ main.py              ← FastAPI app entry
│     └─ ...
│
├─ 📁 .github/                  ← GitHub configuration
│  └─ workflows/               ← CI/CD pipelines
│     ├─ test.yml             ← Run tests on push/PR
│     ├─ build.yml            ← Build releases on tag
│     └─ lint.yml             ← Code quality on push
│
└─ 📁 build/                    ← Build output (git-ignored)
   └─ bin/                      ← Compiled binaries
      ├─ tf-iamgen            ← Main executable (macOS/Linux)
      └─ tf-iamgen.exe        ← Main executable (Windows)
```

## 📍 Quick Navigation

### I want to...

**Add a new CLI command**
- Edit: `cmd/root.go` (register command)
- Create: `cmd/my_command.go`
- Reference: [CONTRIBUTING.md](./CONTRIBUTING.md)

**Add support for a new resource type**
- Edit: `mappings/aws_mapping.yaml`
- Add tests in: `tests/unit/mapping_test.go`
- Create example in: `examples/` (if needed)
- Reference: [docs/MAPPING_FORMAT.md](./MAPPING_FORMAT.md)

**Fix a bug in the parser**
- Edit: `internal/parser/terraform_parser.go`
- Add test in: `tests/unit/parser_test.go`
- Run: `make test`

**Add a new feature to policy generation**
- Edit: `internal/policy/generator.go`
- Format output in: `internal/policy/formatter.go`
- Add tests in: `tests/unit/generator_test.go`

**Improve performance**
- Cache implementation: `internal/mapping/cache.go`
- Performance tests in: `tests/unit/cache_test.go`
- Benchmarks: `make test`

**Add documentation**
- Project docs: `docs/` directory
- Code comments: In source files (Go convention)
- Examples: `examples/` directory

**Set up CI/CD**
- GitHub workflows: `.github/workflows/`
- Build scripts: `scripts/` directory

**Deploy the application**
- Build: `make build`
- Binary output: `build/bin/tf-iamgen`

---

## 📊 Phase-Based Focus Areas

### Phase 1: MVP (Current)
🎯 **Primary directories:**
- `cmd/` - CLI commands
- `internal/parser/` - HCL parsing
- `internal/mapping/` - IAM mappings
- `internal/policy/` - Policy generation
- `mappings/` - YAML database
- `examples/` - Terraform examples
- `tests/unit/` - Unit tests

### Phase 2: CloudTrail Learning
🎯 **Add/expand:**
- `cmd/learn.go` - New command
- `internal/cloudtrail/` - New package
- `ui/backend/` - API server
- Tests in `tests/integration/`

### Phase 3: CI/CD Integration
🎯 **Add/expand:**
- `.github/workflows/` - Plugin actions
- `cmd/` - New commands (simulate, etc.)
- Documentation in `docs/`

### Phase 4: Enterprise
🎯 **Add/expand:**
- `ui/frontend/` - React dashboard
- `ui/backend/` - Full SaaS backend
- Multi-account support

---

## 🔑 Key Files Explained

| File | Purpose |
|------|---------|
| `main.go` | Entry point, initializes CLI |
| `go.mod` | Go dependencies (Cobra, HCL, AWS SDK, etc.) |
| `Makefile` | Common tasks: build, test, lint, etc. |
| `cmd/root.go` | Cobra CLI framework setup |
| `internal/parser/terraform_parser.go` | Terraform HCL parsing logic |
| `internal/mapping/aws_mapping.go` | Resource → action lookup |
| `internal/policy/generator.go` | IAM policy JSON generation |
| `mappings/aws_mapping.yaml` | Terraform resource mappings database |
| `docs/ARCHITECTURE.md` | System design documentation |
| `.github/workflows/*.yml` | CI/CD pipeline definitions |

---

## 💡 Development Workflow

```
1. Pick a task (bug fix, feature, mapping)
   └─ Check docs/CONTRIBUTING.md

2. Create feature branch
   └─ git checkout -b feature/description

3. Make changes
   └─ Edit files in appropriate directories
   └─ Add tests in tests/unit/ or tests/integration/

4. Quality checks
   └─ make fmt        (format code)
   └─ make lint       (check quality)
   └─ make test       (run tests)

5. Commit with clear message
   └─ git commit -m "feat: add XYZ mapping"

6. Push and create PR
   └─ git push origin feature/description
   └─ Create PR on GitHub
```

---

## ✨ Pro Tips

- Use `make help` to see all available tasks
- Run `make all` to format, lint, test, and build
- Use `make dev` for live reloading during development (requires `entr`)
- Check `PROJECT_STRUCTURE.md` for detailed descriptions
- Reference `ARCHITECTURE.md` for system design understanding

---

**Happy coding!** 🚀
