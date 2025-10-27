# 🚀 START HERE

Welcome to **tf-iamgen** — *Automatically discover and generate AWS IAM permissions for Terraform!*

This guide will help you get oriented quickly. **Please read in this order** to understand the project.

---

## 📖 Reading Guide (5 minutes)

### 1. **Overview** (You are here!)
This file — Quick orientation and navigation

### 2. **README.md** (3 min)
- What tf-iamgen does
- Quick start commands
- Key features
- 4-phase roadmap

### 3. **SETUP_COMPLETE.md** (5 min)
- What was created in this setup
- Verification checklist
- Next steps

### 4. **docs/ARCHITECTURE.md** (10 min)
- System design overview
- Component breakdown
- Data flows
- Phase-based architecture

### 5. **docs/CONTRIBUTING.md** (15 min)
- Development setup
- Code style guide
- Testing guidelines
- Workflow for contributions

### 6. **docs/STRUCTURE_GUIDE.md** (5 min)
- Quick reference for file locations
- "I want to..." navigation
- Phase-based focus areas

---

## 🎯 Quick Navigation

| Need | File | Time |
|------|------|------|
| **Project overview** | README.md | 3 min |
| **What was created** | SETUP_COMPLETE.md | 5 min |
| **System design** | docs/ARCHITECTURE.md | 10 min |
| **Development setup** | docs/CONTRIBUTING.md | 15 min |
| **Quick reference** | docs/STRUCTURE_GUIDE.md | 5 min |
| **Detailed structure** | PROJECT_STRUCTURE.md | 10 min |
| **Visual tree** | FOLDER_STRUCTURE.txt | 2 min |

---

## ✅ Project Status

**Phase:** 1 (MVP - Open Source)  
**Status:** ✅ Project structure ready for development  
**Created:** October 27, 2025  
**License:** MIT

### What's Done ✅
- Project structure (18 directories)
- Comprehensive documentation (6+ files)
- Build automation (Makefile)
- Go module configuration
- Contributing guidelines
- MIT License

### What's Next 🔨
- Bootstrap core Go files (main.go, cmd/*, internal/*)
- Create IAM mapping database (YAML)
- Write example Terraform projects
- Build test suite
- Implement CLI commands

---

## 🏗️ Project Structure

```
tf-iamgen/
├── cmd/                    # CLI commands (analyze, generate, learn)
├── internal/               # Core business logic
│   ├── parser/            # Terraform HCL parsing
│   ├── mapping/           # Resource → IAM mapping
│   ├── policy/            # Policy generation
│   └── cloudtrail/        # CloudTrail integration (Phase 2+)
├── mappings/              # IAM mapping database (YAML)
├── examples/              # Example Terraform projects
├── tests/                 # Test suite
├── docs/                  # Documentation
├── scripts/               # Build scripts
└── ui/                    # Web UI (Phase 2+)
```

See `FOLDER_STRUCTURE.txt` for detailed tree view.

---

## 🚀 Getting Started (5 minutes)

### 1. Install Dependencies
```bash
make install-deps
```

### 2. Build the Project
```bash
make build
```

### 3. Run Tests
```bash
make test
```

### 4. See All Commands
```bash
make help
```

---

## 📊 The Big Picture

### What tf-iamgen Does

1. **Analyze** → Scan your Terraform project
2. **Map** → Identify required AWS IAM actions
3. **Generate** → Create least-privilege IAM policies
4. **Output** → Valid JSON ready to deploy

### Why It Matters

- ✅ **Least-Privilege**: Only grant permissions you actually need
- ✅ **Automated**: No manual IAM policy writing
- ✅ **Secure**: Prevent over-privileged deployments
- ✅ **DevOps-Friendly**: Integrate into CI/CD

---

## 💡 Key Concepts

| Term | Meaning |
|------|---------|
| **Resource Type** | Terraform AWS resource (e.g., `aws_s3_bucket`) |
| **IAM Action** | AWS permission (e.g., `s3:CreateBucket`) |
| **Mapping** | YAML file: Resource Type → IAM Actions |
| **Least-Privilege** | Minimal permissions needed |
| **Policy** | AWS IAM policy JSON document |

---

## 🎯 4-Phase Roadmap

### Phase 1: MVP (Open Source) ← **You are here** ✨
- CLI tool in Go
- Terraform static analysis
- IAM policy generation
- Free and open source

### Phase 2: CloudTrail Learning (Commercial)
- Analyze actual AWS API calls
- Dynamic least-privilege refinement
- SaaS dashboard

### Phase 3: CI/CD Integration (Commercial)
- GitHub Actions plugin
- GitLab CI integration
- Pre-deployment validation

### Phase 4: Enterprise (Commercial Advanced)
- Multi-account management
- ML-driven optimization
- Compliance reporting

See `docs/PHASES.md` for details (coming soon).

---

## 🛠️ Development Commands

```bash
# Quick reference
make help              # Show all commands
make build             # Compile
make test              # Run tests
make lint              # Check code quality
make fmt               # Format code
make clean             # Clean artifacts
make all               # Format + lint + test + build
```

---

## 📚 Documentation Index

### Getting Started
- `README.md` - Project overview
- `START_HERE.md` - This file
- `SETUP_COMPLETE.md` - Setup summary

### Architecture & Design
- `docs/ARCHITECTURE.md` - System design
- `PROJECT_STRUCTURE.md` - Directory details
- `FOLDER_STRUCTURE.txt` - Visual tree

### Development
- `docs/CONTRIBUTING.md` - Developer guide
- `docs/STRUCTURE_GUIDE.md` - Quick reference
- `Makefile` - Build tasks

### Future Documentation
- `docs/MAPPING_FORMAT.md` - How to write mappings
- `docs/API.md` - CLI reference
- `docs/EXAMPLES.md` - Usage examples
- `docs/SECURITY.md` - Security & privacy

---

## ❓ Frequently Asked Questions

### Q: Where do I start coding?
**A:** Begin with `main.go` (entry point) and `cmd/root.go` (CLI setup). See `docs/CONTRIBUTING.md` for guidance.

### Q: How do I add a new AWS resource type?
**A:** Edit `mappings/aws_mapping.yaml` to add resource → IAM actions mappings. See future `docs/MAPPING_FORMAT.md`.

### Q: How does the system work end-to-end?
**A:** Read `docs/ARCHITECTURE.md` for a complete walkthrough of data flows and components.

### Q: Where is the code I should write?
**A:** Look at these files to create:
- `main.go` - Entry point
- `cmd/analyze.go` - Analyze command
- `cmd/generate.go` - Generate command
- `internal/parser/terraform_parser.go` - HCL parser
- `internal/mapping/aws_mapping.go` - Mapping manager
- `internal/policy/generator.go` - Policy generator

### Q: How do I contribute?
**A:** See `docs/CONTRIBUTING.md` for the full workflow: fork, branch, code, test, commit, PR.

---

## 🔍 Quick Help

### I want to understand...
| Topic | File |
|-------|------|
| Project overview | README.md |
| System architecture | docs/ARCHITECTURE.md |
| Folder structure | docs/STRUCTURE_GUIDE.md |
| Development setup | docs/CONTRIBUTING.md |
| Detailed directories | PROJECT_STRUCTURE.md |

### I want to...
| Task | Location |
|------|----------|
| Add CLI command | `cmd/my_command.go` |
| Add resource type | `mappings/aws_mapping.yaml` |
| Fix parser bug | `internal/parser/*.go` |
| Improve policy generation | `internal/policy/generator.go` |
| Add tests | `tests/unit/*_test.go` |

---

## 📞 Summary

**You have:**
- ✨ Clean, organized project structure
- 📚 Comprehensive documentation
- 🛠️ Build automation ready
- 🧪 Test framework in place
- 📋 Contributing guidelines

**Next steps:**
1. Read `README.md` (quick overview)
2. Read `docs/ARCHITECTURE.md` (system design)
3. Create bootstrap Go files
4. Build and test
5. Start coding!

---

## 🎓 Learning Resources

- [Effective Go](https://golang.org/doc/effective_go) - Go style guide
- [Cobra CLI](https://github.com/spf13/cobra) - CLI framework
- [HCL2 Parser](https://github.com/hashicorp/hcl) - Terraform parser
- [AWS IAM](https://docs.aws.amazon.com/iam/) - IAM reference

---

## ✨ You're Ready!

Everything is set up. The structure is clean. The docs are comprehensive.

**Start with:** `README.md` → `docs/ARCHITECTURE.md` → `docs/CONTRIBUTING.md`

Then build something awesome! 🚀

---

**Status:** ✅ Ready for Phase 1 Development  
**Created:** October 27, 2025  
**License:** MIT  

**Let's make infrastructure security better, one policy at a time.** 🛡️
