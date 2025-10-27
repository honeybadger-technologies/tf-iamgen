package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/honeybadger/tf-iamgen/internal/parser"
)

// TestNewTerraformParser tests that a new parser can be created
func TestNewTerraformParser(t *testing.T) {
	p := parser.NewTerraformParser()
	if p == nil {
		t.Fatal("NewTerraformParser returned nil")
	}
}

// TestParseSimpleVPC tests parsing the simple VPC example
func TestParseSimpleVPC(t *testing.T) {
	parser := parser.NewTerraformParser()

	result, err := parser.ParseDirectory("../examples/simple_vpc")
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}

	if result == nil {
		t.Fatal("ParseDirectory returned nil result")
	}

	// Check that we found resources
	if len(result.Resources) == 0 {
		t.Fatal("Expected to find resources in simple_vpc example, but found none")
	}

	// Check specific resources
	expectedResources := map[string]bool{
		"aws_vpc":              false,
		"aws_internet_gateway": false,
		"aws_subnet":           false,
		"aws_security_group":   false,
		"aws_route_table":      false,
	}

	for _, res := range result.Resources {
		if _, exists := expectedResources[res.Type]; exists {
			expectedResources[res.Type] = true
		}
	}

	// Verify we found each expected resource type
	for resType, found := range expectedResources {
		if !found {
			t.Errorf("Did not find resource type %s", resType)
		}
	}

	t.Logf("Found %d resources in simple_vpc example", len(result.Resources))
	for _, res := range result.Resources {
		t.Logf("  - %s.%s", res.Type, res.Name)
	}
}

// TestParseSimpleS3 tests parsing the simple S3 example
func TestParseSimpleS3(t *testing.T) {
	parser := parser.NewTerraformParser()

	result, err := parser.ParseDirectory("../examples/simple_s3")
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}

	if result == nil {
		t.Fatal("ParseDirectory returned nil result")
	}

	// Check that we found resources
	if len(result.Resources) == 0 {
		t.Fatal("Expected to find resources in simple_s3 example, but found none")
	}

	// Check for S3 specific resources
	hasS3Bucket := false
	for _, res := range result.Resources {
		if res.Type == "aws_s3_bucket" {
			hasS3Bucket = true
		}
	}

	if !hasS3Bucket {
		t.Error("Expected to find aws_s3_bucket resource")
	}

	t.Logf("Found %d resources in simple_s3 example", len(result.Resources))
	for _, res := range result.Resources {
		t.Logf("  - %s.%s", res.Type, res.Name)
	}
}

// TestResourceMetadata tests the GetResourceMetadata function
func TestResourceMetadata(t *testing.T) {
	tests := []struct {
		resourceType    string
		expectedType    string
		expectedService string
	}{
		{"aws_s3_bucket", "aws_s3_bucket", "s3"},
		{"aws_instance", "aws_instance", "ec2"},
		{"aws_db_instance", "aws_db_instance", "db"},
		{"aws_iam_role", "aws_iam_role", "iam"},
		{"aws_lambda_function", "aws_lambda_function", "lambda"},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			metadata := parser.GetResourceMetadata(tt.resourceType)
			if metadata.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, metadata.Type)
			}
			if metadata.Service != tt.expectedService {
				t.Errorf("Expected service %s, got %s", tt.expectedService, metadata.Service)
			}
		})
	}
}

// TestIsAWSResource tests the IsAWSResource function
func TestIsAWSResource(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{"aws_s3_bucket", true},
		{"aws_instance", true},
		{"aws_vpc", true},
		{"local_file", false},
		{"random_id", false},
		{"google_storage_bucket", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			result := parser.IsAWSResource(tt.resourceType)
			if result != tt.expected {
				t.Errorf("IsAWSResource(%s) = %v, expected %v", tt.resourceType, result, tt.expected)
			}
		})
	}
}

// TestIsKnownResource tests the IsKnownResource function
func TestIsKnownResource(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{"aws_s3_bucket", true},
		{"aws_instance", true},
		{"aws_custom_resource", false},
		{"local_file", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			result := parser.IsKnownResource(tt.resourceType)
			if result != tt.expected {
				t.Errorf("IsKnownResource(%s) = %v, expected %v", tt.resourceType, result, tt.expected)
			}
		})
	}
}

// TestGetServiceFromResourceType tests service extraction
func TestGetServiceFromResourceType(t *testing.T) {
	tests := []struct {
		resourceType    string
		expectedService string
	}{
		{"aws_s3_bucket", "s3"},
		{"aws_instance", "instance"},
		{"aws_db_instance", "db"},
		{"aws_lambda_function", "lambda"},
		{"aws_iam_role", "iam"},
		{"aws_vpc", "vpc"},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			service := parser.GetServiceFromResourceType(tt.resourceType)
			if service != tt.expectedService {
				t.Errorf("GetServiceFromResourceType(%s) = %s, expected %s", tt.resourceType, service, tt.expectedService)
			}
		})
	}
}

// TestGetResourceCategoryFromResourceType tests category extraction
func TestGetResourceCategoryFromResourceType(t *testing.T) {
	tests := []struct {
		resourceType     string
		expectedCategory string
	}{
		{"aws_s3_bucket", "bucket"},
		{"aws_instance", "instance"},
		{"aws_db_instance", "db_instance"},
		{"aws_lambda_function", "lambda_function"},
		{"aws_iam_role", "iam_role"},
		{"aws_security_group_rule", "security_group_rule"},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			category := parser.GetResourceCategoryFromResourceType(tt.resourceType)
			if category != tt.expectedCategory {
				t.Errorf("GetResourceCategoryFromResourceType(%s) = %s, expected %s", tt.resourceType, category, tt.expectedCategory)
			}
		})
	}
}

// TestParseResultMethods tests ParseResult helper methods
func TestParseResultMethods(t *testing.T) {
	result := &parser.ParseResult{
		Resources: []parser.Resource{
			{Type: "aws_s3_bucket", Name: "bucket1"},
			{Type: "aws_s3_bucket", Name: "bucket2"},
			{Type: "aws_instance", Name: "app1"},
			{Type: "aws_vpc", Name: "main"},
		},
		FilesProcessed: 3,
		TotalResources: 4,
	}

	// Test GetResourcesByType
	s3Buckets := result.GetResourcesByType("aws_s3_bucket")
	if len(s3Buckets) != 2 {
		t.Errorf("Expected 2 S3 buckets, got %d", len(s3Buckets))
	}

	// Test GetResourcesByService
	ec2Resources := result.GetResourcesByService("ec2")
	if len(ec2Resources) != 1 {
		t.Errorf("Expected 1 EC2 resource, got %d", len(ec2Resources))
	}

	// Test Summary
	summary := result.Summary()
	if len(summary) == 0 {
		t.Error("Summary returned empty string")
	}
	t.Logf("ParseResult Summary:\n%s", summary)
}

// TestParseInvalidDirectory tests parsing an invalid directory
func TestParseInvalidDirectory(t *testing.T) {
	parser := parser.NewTerraformParser()
	_, err := parser.ParseDirectory("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error when parsing nonexistent directory")
	}
}

// TestResourceString tests the Resource String method
func TestResourceString(t *testing.T) {
	resource := parser.Resource{
		Type:       "aws_s3_bucket",
		Name:       "my_bucket",
		FilePath:   "main.tf",
		LineNumber: 10,
	}

	str := resource.String()
	if str != `resource "aws_s3_bucket" "my_bucket" (main.tf:10)` {
		t.Errorf("Unexpected string representation: %s", str)
	}
}

// TestResourceFullName tests the Resource FullName method
func TestResourceFullName(t *testing.T) {
	resource := parser.Resource{
		Type: "aws_s3_bucket",
		Name: "my_bucket",
	}

	fullName := resource.FullName()
	if fullName != "aws_s3_bucket.my_bucket" {
		t.Errorf("Expected 'aws_s3_bucket.my_bucket', got '%s'", fullName)
	}
}

// TestParseErrorString tests the ParseError Error method
func TestParseErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      parser.ParseError
		expected string
	}{
		{
			name: "With line and column",
			err: parser.ParseError{
				FilePath:  "main.tf",
				Line:      10,
				Column:    5,
				Message:   "Syntax error",
				ErrorType: "syntax",
			},
			expected: "main.tf:10:5: Syntax error (syntax)",
		},
		{
			name: "With line only",
			err: parser.ParseError{
				FilePath:  "main.tf",
				Line:      10,
				Message:   "Parse error",
				ErrorType: "parse_error",
			},
			expected: "main.tf:10: Parse error (parse_error)",
		},
		{
			name: "File only",
			err: parser.ParseError{
				FilePath:  "main.tf",
				Message:   "File error",
				ErrorType: "file",
			},
			expected: "main.tf: File error (file)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// BenchmarkParseDirectory benchmarks parsing a directory
func BenchmarkParseDirectory(b *testing.B) {
	// Get the absolute path to examples/simple_vpc
	exPath, err := filepath.Abs("../examples/simple_vpc")
	if err != nil {
		b.Fatalf("Failed to get absolute path: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(exPath); os.IsNotExist(err) {
		b.Skipf("Example directory not found at %s", exPath)
	}

	for i := 0; i < b.N; i++ {
		p := parser.NewTerraformParser()
		_, err := p.ParseDirectory(exPath)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// BenchmarkResourceTypeExtraction benchmarks service/category extraction
func BenchmarkResourceTypeExtraction(b *testing.B) {
	resourceTypes := []string{
		"aws_s3_bucket",
		"aws_instance",
		"aws_db_instance",
		"aws_lambda_function",
		"aws_iam_role",
		"aws_security_group_rule",
	}

	for i := 0; i < b.N; i++ {
		for _, rt := range resourceTypes {
			parser.GetServiceFromResourceType(rt)
			parser.GetResourceCategoryFromResourceType(rt)
		}
	}
}
