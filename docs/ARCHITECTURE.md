# tf-iamgen Architecture

## System Overview

```
┌────────────────────────────────────────────────────────────────┐
│                    Terraform Project Files                      │
│                  (*.tf, *.tf.json in directory)                 │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ▼
┌────────────────────────────────────────────────────────────────┐
│              Parser (HCL → AST Extraction)                      │
│  - Parse Terraform HCL using hashicorp/hcl/v2                  │
│  - Extract resource definitions                                │
│  - Identify: type, name, attributes                            │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ▼
┌────────────────────────────────────────────────────────────────┐
│         Mapper (Resource Type → IAM Actions)                    │
│  - Load YAML mapping database at startup                       │
│  - Cache results for performance                               │
│  - Lookup: aws_s3_bucket → [s3:CreateBucket, ...]             │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                          ▼
┌────────────────────────────────────────────────────────────────┐
│         Policy Generator (Actions → IAM Policy)                │
│  - Aggregate discovered actions                                │
│  - Group by AWS service                                        │
│  - Build valid AWS IAM policy JSON                             │
└─────────────────────────┬──────────────────────────────────────┘
                          │
                ┌─────────┴──────────┬──────────────┐
                │                    │              │
                ▼                    ▼              ▼
          JSON Policy         Summary Report    stdout/file
         (policy.json)     (actions grouped)   (human readable)
```

## Phase Architecture

### Phase 1: Static Analysis (MVP)

**Single-pass workflow:**
1. **Input**: Directory containing Terraform files
2. **Processing**: Static parsing → mapping lookup → policy generation
3. **Output**: IAM policy JSON or human-readable summary

**Components:**
- `cmd/analyze`: List discovered resources and actions
- `cmd/generate`: Output IAM policy JSON
- `internal/parser`: HCL parsing
- `internal/mapping`: YAML database + cache
- `internal/policy`: Policy generation and formatting

**Data Storage:**
- `mappings/aws_mapping.yaml`: Resource → IAM actions mapping
- Local files only, no external services

### Phase 2: CloudTrail Learning (Commercial)

**Two-path workflow:**

```
Static Analysis Path (existing)
        │
        ├─→ Predicted IAM Actions
        │
        ▼
┌─────────────────────────┐
│   Differential Engine   │
│  (Predicted vs Actual)  │
└──────────┬──────────────┘
           │
           ▼
    Refined Policy
   (Least-Privilege)

        │
CloudTrail Analysis Path (new)
        │
        ├─→ Query CloudTrail events
        ├─→ Extract actual IAM actions used
        ├─→ Aggregate and deduplicate
        │
        └─→ Observed IAM Actions
```

**New Components:**
- `cmd/learn`: CloudTrail integration command
- `internal/cloudtrail/collector`: Event retrieval
- `internal/cloudtrail/analyzer`: Event processing
- `ui/backend`: API for storing/analyzing data (optional SaaS)
- `ui/frontend`: Dashboard for visualization

**Data Storage:**
- AWS CloudTrail events (read-only)
- Optional backend: PostgreSQL for aggregation
- Optional SaaS: Anonymized event data

### Phase 3: CI/CD Integration

**Integration points:**
- GitHub Actions plugin: Pre-deployment validation
- GitLab CI: Policy compliance checks
- AWS IAM Access Analyzer: Validation against best practices
- Role simulation: Dry-run execution with predicted permissions

### Phase 4: Enterprise

**Multi-tenant SaaS:**
- Centralized dashboard
- Multi-account management
- Role recommendation engine (ML)
- Compliance reporting

## Core Components Deep Dive

### Parser (`internal/parser/`)

**Responsibility**: Convert Terraform HCL into structured data

**Implementation:**
```
HCL Code
   │
   ├─→ hashicorp/hcl/v2 parser
   │
   ├─→ AST (Abstract Syntax Tree)
   │
   └─→ Extracted Resources
       ├─ Type: aws_s3_bucket
       ├─ Name: my_bucket
       └─ Attributes: {...}
```

**Handles:**
- `.tf` files (HCL syntax)
- `.tf.json` files (JSON syntax)
- Variable references and interpolations (basic)
- Nested blocks (modules, resources, locals)

**Output**: `[]Resource` — slice of discovered resources

### Mapper (`internal/mapping/`)

**Responsibility**: Map Terraform resources to IAM actions

**Architecture:**
```
YAML Mapping File
   │
   ├─→ Load at startup
   │
   ├─→ In-memory cache (map[string][]string)
   │   Key: resource type (e.g., "aws_s3_bucket")
   │   Value: IAM actions (e.g., ["s3:CreateBucket", ...])
   │
   └─→ Lookup service
       Input: Resource type
       Output: []Action
```

**Features:**
- Lazy-loading: Mappings loaded once at startup
- LRU cache: Recent lookups cached in memory
- Hierarchical structure: Support for grouping actions by operation type

**Example Mapping:**
```yaml
aws_s3_bucket:
  all:
    - s3:CreateBucket
    - s3:PutBucketAcl
    - s3:GetBucket
    - s3:ListBucket
```

### Policy Generator (`internal/policy/`)

**Responsibility**: Generate AWS IAM policy JSON

**Workflow:**
```
Discovered Resources
   │
   ├─→ Get actions for each resource type
   │
   ├─→ Aggregate by AWS service (s3:, ec2:, etc.)
   │
   ├─→ Deduplicate actions
   │
   ├─→ Build IAM Statement
   │   {
   │     "Effect": "Allow",
   │     "Action": ["s3:*"],
   │     "Resource": "*"
   │   }
   │
   └─→ Output valid IAM Policy JSON
```

**Features:**
- Statement grouping by service
- Action deduplication
- Customizable Resource ARNs (phase 2+)
- Conditions support (future)

### CloudTrail Collector (`internal/cloudtrail/`)

**Responsibility**: Query and analyze CloudTrail events (Phase 2+)

**Workflow:**
```
AWS CloudTrail
   │
   ├─→ Query recent events
   │
   ├─→ Filter by principal/role
   │
   ├─→ Extract EventName → IAM Action mapping
   │
   ├─→ Aggregate actions used
   │
   └─→ Output observed IAM actions
```

**Note**: Not implemented in Phase 1 MVP

## Data Structures

### Resource (Parser Output)

```go
type Resource struct {
    Type       string                 // e.g., "aws_s3_bucket"
    Name       string                 // e.g., "my_bucket"
    Attributes map[string]interface{} // Parsed attributes
}
```

### IAM Action

```go
type Action string // e.g., "s3:CreateBucket"
```

### IAM Policy

```go
type Policy struct {
    Version   string      // "2012-10-17"
    Statement []Statement
}

type Statement struct {
    Effect   string   // "Allow" | "Deny"
    Action   []string // IAM actions
    Resource []string // Resource ARNs
}
```

## Performance Considerations

### Caching Strategy

1. **Mapping Cache**: In-memory LRU cache for resource → actions lookups
   - Reduces YAML parsing overhead
   - Configurable size (default: 1000 entries)

2. **Parser Cache**: (Future) Cache Terraform AST between runs
   - File modification time check
   - Only re-parse changed files

### Benchmarks (Target for Phase 1)

- Parse 100 Terraform files: < 500ms
- Generate policy: < 50ms
- Total execution: < 1s

## Error Handling

### Parser Errors
- Invalid HCL syntax → Error message + line number
- Unsupported resource types → Warn, continue parsing

### Mapping Errors
- Unknown resource type → Warn, no actions added
- Malformed YAML → Fatal error

### Policy Generation Errors
- Empty action set → Valid policy with no permissions
- JSON encoding error → Fatal error

## Security Considerations

### Phase 1 (Open Source)
- No Terraform code transmitted anywhere
- All processing local
- No credentials required
- No network access

### Phase 2+ (Commercial)
- Optional CloudTrail analysis (AWS credentials required)
- Optional SaaS backend (data encryption in transit/at rest)
- GDPR-compliant data handling
- Anonymized event data storage

## Extension Points

### Adding New Resource Type
1. Create YAML entry in `mappings/aws_mapping.yaml`
2. Test with example Terraform
3. Add unit test

### Adding New IAM Action Discovery Method
1. Create new package in `internal/`
2. Implement common interface
3. Integrate with policy generator

### Adding New Output Format
1. Implement formatter in `internal/policy/formatter.go`
2. Add CLI flag in `cmd/generate.go`
3. Add tests

---

**Next Steps:** See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup.
