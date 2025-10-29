package policy

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/honeybadger/tf-iamgen/internal/mapping"
	"github.com/honeybadger/tf-iamgen/internal/parser"
)

// TestGeneratePolicyBasic tests basic policy generation
func TestGeneratePolicyBasic(t *testing.T) {
	// Create a mock mapping service
	mappingService := createMockMappingService()

	// Create generator with default options
	opts := PolicyGenerationOptions{
		GroupBy:              "flat",
		UseWildcardResources: true,
		IncludeSids:          true,
	}
	gen := NewGenerator(mappingService, opts)

	// Create a parse result with resources
	parseResult := &parser.ParseResult{
		Resources: []parser.Resource{
			{
				Type: "aws_s3_bucket",
				Name: "my_bucket",
			},
		},
	}

	policy, metadata, err := gen.GeneratePolicy(parseResult)

	if err != nil {
		t.Fatalf("GeneratePolicy failed: %v", err)
	}

	if policy == nil {
		t.Fatalf("Generated policy is nil")
	}

	if len(policy.Statement) == 0 {
		t.Fatalf("Policy has no statements")
	}

	if metadata.ResourceCount != 1 {
		t.Errorf("Expected 1 resource, got %d", metadata.ResourceCount)
	}

	if metadata.ActionCount == 0 {
		t.Errorf("Expected non-zero action count")
	}
}

// TestGeneratePolicyGroupedByService tests policy generation grouped by service
func TestGeneratePolicyGroupedByService(t *testing.T) {
	mappingService := createMockMappingService()

	opts := PolicyGenerationOptions{
		GroupBy:              "service",
		UseWildcardResources: true,
		IncludeSids:          true,
	}
	gen := NewGenerator(mappingService, opts)

	parseResult := &parser.ParseResult{
		Resources: []parser.Resource{
			{Type: "aws_s3_bucket", Name: "bucket1"},
			{Type: "aws_instance", Name: "instance1"},
		},
	}

	policy, _, err := gen.GeneratePolicy(parseResult)

	if err != nil {
		t.Fatalf("GeneratePolicy failed: %v", err)
	}

	// Should have statements grouped by service
	if len(policy.Statement) == 0 {
		t.Fatalf("Policy has no statements")
	}

	// Check that statements are properly grouped
	hasS3 := false

	for _, stmt := range policy.Statement {
		for _, action := range stmt.Action {
			if len(action) > 0 && action[0:1] == "s" && len(action) > 1 && action[1:2] == "3" {
				hasS3 = true
			}
		}
	}

	if !hasS3 {
		t.Error("Expected S3 actions in policy")
	}
}

// TestGeneratePolicyJSON tests JSON output formatting
func TestGeneratePolicyJSON(t *testing.T) {
	policy := NewPolicy()
	policy.AddStatement(Statement{
		Sid:      "TestStatement",
		Effect:   EffectAllow,
		Action:   []string{"s3:GetObject"},
		Resource: []string{"arn:aws:s3:::bucket/*"},
	})

	jsonStr, err := policy.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if jsonStr == "" {
		t.Fatalf("JSON output is empty")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Check structure
	if version, ok := result["Version"].(string); !ok || version != "2012-10-17" {
		t.Errorf("Expected valid policy version")
	}

	if statements, ok := result["Statement"].([]interface{}); !ok || len(statements) != 1 {
		t.Errorf("Expected 1 statement in policy")
	}
}

// TestGeneratePolicyCompactJSON tests compact JSON output
func TestGeneratePolicyCompactJSON(t *testing.T) {
	policy := NewPolicy()
	policy.AddStatement(Statement{
		Sid:      "TestStatement",
		Effect:   EffectAllow,
		Action:   []string{"s3:GetObject"},
		Resource: []string{"arn:aws:s3:::bucket/*"},
	})

	compactJSON, err := policy.ToCompactJSON()
	if err != nil {
		t.Fatalf("ToCompactJSON failed: %v", err)
	}

	formattedJSON, err := policy.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Compact should be smaller
	if len(compactJSON) >= len(formattedJSON) {
		t.Errorf("Compact JSON should be smaller than formatted JSON")
	}
}

// TestValidatePolicy tests policy validation
func TestValidatePolicy(t *testing.T) {
	mappingService := createMockMappingService()
	gen := NewGenerator(mappingService, PolicyGenerationOptions{})

	// Valid policy
	validPolicy := NewPolicy()
	validPolicy.AddStatement(Statement{
		Effect:   EffectAllow,
		Action:   []string{"s3:GetObject"},
		Resource: []string{"*"},
	})

	warnings, err := gen.ValidatePolicy(validPolicy)
	if err != nil {
		t.Fatalf("ValidatePolicy failed: %v", err)
	}

	// Should have no critical errors
	if len(warnings) > 0 {
		t.Logf("Validation warnings: %v", warnings)
	}

	// Empty policy
	emptyPolicy := NewPolicy()
	warnings, err = gen.ValidatePolicy(emptyPolicy)
	if err != nil {
		t.Fatalf("ValidatePolicy failed: %v", err)
	}

	if len(warnings) == 0 {
		t.Errorf("Expected warnings for empty policy")
	}
}

// TestMergeActions tests action merging
func TestMergeActions(t *testing.T) {
	actions1 := []string{"s3:GetObject", "s3:ListBucket"}
	actions2 := []string{"s3:GetObject", "s3:DeleteObject"} // Duplicate GetObject
	actions3 := []string{"ec2:RunInstances"}

	merged := MergeActions(actions1, actions2, actions3)

	if len(merged) != 4 {
		t.Errorf("Expected 4 unique actions, got %d", len(merged))
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, action := range merged {
		if seen[action] {
			t.Errorf("Found duplicate action: %s", action)
		}
		seen[action] = true
	}
}

// TestMergeResources tests resource merging
func TestMergeResources(t *testing.T) {
	resources1 := []string{"arn:aws:s3:::bucket1", "arn:aws:s3:::bucket2"}
	resources2 := []string{"arn:aws:s3:::bucket1", "arn:aws:s3:::bucket3"}

	merged := MergeResources(resources1, resources2)

	if len(merged) != 3 {
		t.Errorf("Expected 3 unique resources, got %d", len(merged))
	}
}

// TestPolicyBuilder tests the policy builder
func TestPolicyBuilder(t *testing.T) {
	opts := PolicyGenerationOptions{GroupBy: "service"}
	builder := NewPolicyBuilder(opts)

	builder.AddActionStatement("S3Access", []string{"s3:GetObject"}, []string{"*"})
	builder.AddActionStatement("EC2Access", []string{"ec2:RunInstances"}, []string{"*"})

	policy, _ := builder.Build()

	if len(policy.Statement) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(policy.Statement))
	}

	// Verify statements
	if policy.Statement[0].Sid != "EC2Access" && policy.Statement[0].Sid != "S3Access" {
		t.Errorf("Unexpected statement Sid: %s", policy.Statement[0].Sid)
	}
}

// TestPolicyMetadata tests metadata generation
func TestPolicyMetadata(t *testing.T) {
	mappingService := createMockMappingService()
	gen := NewGenerator(mappingService, PolicyGenerationOptions{})

	parseResult := &parser.ParseResult{
		Resources: []parser.Resource{
			{Type: "aws_s3_bucket", Name: "bucket1"},
		},
	}

	_, metadata, err := gen.GeneratePolicy(parseResult)
	if err != nil {
		t.Fatalf("GeneratePolicy failed: %v", err)
	}

	// Check metadata fields
	if metadata.ResourceCount != 1 {
		t.Errorf("Expected 1 resource in metadata")
	}

	if metadata.ActionCount == 0 {
		t.Errorf("Expected non-zero action count")
	}

	if metadata.GeneratedAt == "" {
		t.Errorf("GeneratedAt is empty")
	}

	// Verify timestamp format
	if _, err := time.Parse(time.RFC3339, metadata.GeneratedAt); err != nil {
		t.Errorf("GeneratedAt is not valid RFC3339: %v", err)
	}

	if len(metadata.Services) == 0 {
		t.Errorf("Expected services in metadata")
	}
}

// TestGroupStatementsByService tests service grouping
func TestGroupStatementsByService(t *testing.T) {
	statements := []Statement{
		{
			Effect:   EffectAllow,
			Action:   []string{"s3:GetObject", "s3:ListBucket"},
			Resource: []string{"*"},
		},
		{
			Effect:   EffectAllow,
			Action:   []string{"ec2:RunInstances", "ec2:TerminateInstances"},
			Resource: []string{"*"},
		},
	}

	grouped := GroupStatementsByService(statements)

	if len(grouped) != 2 {
		t.Errorf("Expected 2 service groups, got %d", len(grouped))
	}

	if _, hasS3 := grouped["s3"]; !hasS3 {
		t.Error("Expected S3 group")
	}

	if _, hasEC2 := grouped["ec2"]; !hasEC2 {
		t.Error("Expected EC2 group")
	}
}

// TestGetServicesFromStatements tests service extraction
func TestGetServicesFromStatements(t *testing.T) {
	statements := []Statement{
		{
			Effect:   EffectAllow,
			Action:   []string{"s3:GetObject"},
			Resource: []string{"*"},
		},
		{
			Effect:   EffectAllow,
			Action:   []string{"ec2:RunInstances"},
			Resource: []string{"*"},
		},
		{
			Effect:   EffectAllow,
			Action:   []string{"s3:ListBucket"},
			Resource: []string{"*"},
		},
	}

	services := GetServicesFromStatements(statements)

	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}

	serviceMap := make(map[string]bool)
	for _, s := range services {
		serviceMap[s] = true
	}

	if !serviceMap["s3"] || !serviceMap["ec2"] {
		t.Error("Expected S3 and EC2 services")
	}
}

// TestGeneratePolicyNilInput tests nil input handling
func TestGeneratePolicyNilInput(t *testing.T) {
	mappingService := createMockMappingService()
	gen := NewGenerator(mappingService, PolicyGenerationOptions{})

	_, _, err := gen.GeneratePolicy(nil)

	if err == nil {
		t.Fatalf("Expected error for nil parse result")
	}
}

// TestPermissionNarrower tests the permission narrower
func TestPermissionNarrower(t *testing.T) {
	narrower := NewPermissionNarrower()

	actions := []string{"s3:GetObject", "s3:ListBucket", "s3:DeleteObject"}

	narrowed := narrower.NarrowPermissions("aws_s3_bucket", actions)

	if len(narrowed) == 0 {
		t.Errorf("Expected narrowed actions")
	}
}

// TestActionUsageRecording tests recording action usage
func TestActionUsageRecording(t *testing.T) {
	narrower := NewPermissionNarrower()

	narrower.RecordActionUsage("aws_s3_bucket", "s3:GetObject")
	narrower.RecordActionUsage("aws_s3_bucket", "s3:ListBucket")

	used := narrower.GetUsedActions("aws_s3_bucket")

	if len(used) != 2 {
		t.Errorf("Expected 2 used actions, got %d", len(used))
	}
}

// Helper function to create a mock mapping service
func createMockMappingService() *mapping.MappingService {
	db := mapping.NewMappingDatabase()

	// Create mock mappings for testing
	s3Actions := mapping.NewActionSet(
		"s3:GetObject",
		"s3:ListBucket",
		"s3:CreateBucket",
		"s3:DeleteBucket",
	)

	ec2Actions := mapping.NewActionSet(
		"ec2:RunInstances",
		"ec2:TerminateInstances",
		"ec2:DescribeInstances",
	)

	iamActions := mapping.NewActionSet(
		"iam:CreateRole",
		"iam:DeleteRole",
		"iam:AttachRolePolicy",
	)

	s3Mapping := &mapping.ResourceActionMap{
		Actions: map[string]mapping.ActionSet{
			"aws_s3_bucket": s3Actions,
		},
		Service:     "s3",
		Description: "S3 bucket",
	}

	ec2Mapping := &mapping.ResourceActionMap{
		Actions: map[string]mapping.ActionSet{
			"aws_instance": ec2Actions,
		},
		Service:     "ec2",
		Description: "EC2 instance",
	}

	iamMapping := &mapping.ResourceActionMap{
		Actions: map[string]mapping.ActionSet{
			"aws_iam_role": iamActions,
		},
		Service:     "iam",
		Description: "IAM role",
	}

	db.AddMappingForTesting("aws_s3_bucket", s3Mapping)
	db.AddMappingForTesting("aws_instance", ec2Mapping)
	db.AddMappingForTesting("aws_iam_role", iamMapping)

	service := mapping.NewMappingService(db)
	return service
}

// BenchmarkGeneratePolicy benchmarks policy generation
func BenchmarkGeneratePolicy(b *testing.B) {
	mappingService := createMockMappingService()
	gen := NewGenerator(mappingService, PolicyGenerationOptions{GroupBy: "service"})

	parseResult := &parser.ParseResult{
		Resources: []parser.Resource{
			{Type: "aws_s3_bucket", Name: "bucket1"},
			{Type: "aws_instance", Name: "instance1"},
			{Type: "aws_iam_role", Name: "role1"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.GeneratePolicy(parseResult)
	}
}

// BenchmarkPolicyJSONMarshaling benchmarks JSON marshaling
func BenchmarkPolicyJSONMarshaling(b *testing.B) {
	policy := NewPolicy()
	for i := 0; i < 10; i++ {
		policy.AddStatement(Statement{
			Effect:   EffectAllow,
			Action:   []string{"s3:GetObject", "s3:ListBucket"},
			Resource: []string{"*"},
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		policy.ToJSON()
	}
}
