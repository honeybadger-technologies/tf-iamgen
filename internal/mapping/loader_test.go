package mapping

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadMappingsFromDirectory tests loading YAML files from a directory
func TestLoadMappingsFromDirectory(t *testing.T) {
	db := NewMappingDatabase()

	// Try to load from the actual mappings directory
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	if !db.IsLoaded() {
		t.Error("Expected database to be marked as loaded")
	}

	// Check that mappings were loaded
	allMappings := db.GetAllMappings()
	if len(allMappings) == 0 {
		t.Error("Expected mappings to be loaded, got none")
	}
}

// TestLoadMappingsDirectoryNotFound tests error when directory doesn't exist
func TestLoadMappingsDirectoryNotFound(t *testing.T) {
	db := NewMappingDatabase()

	err := db.LoadMappings("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error when loading from nonexistent directory")
	}
}

// TestGetMapping tests retrieving a specific mapping
func TestGetMapping(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	mapping, exists := db.GetMapping("aws_s3_bucket")
	if !exists {
		t.Error("Expected aws_s3_bucket mapping to exist")
	}

	if mapping == nil {
		t.Error("Expected non-nil mapping for aws_s3_bucket")
	}

	if mapping.Service != "s3" {
		t.Errorf("Expected service s3, got %s", mapping.Service)
	}
}

// TestHasMapping tests checking if mapping exists
func TestHasMapping(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	if !db.HasMapping("aws_s3_bucket") {
		t.Error("Expected HasMapping to return true for aws_s3_bucket")
	}

	if db.HasMapping("nonexistent_resource") {
		t.Error("Expected HasMapping to return false for nonexistent resource")
	}
}

// TestGetAllMappings tests retrieving all mappings
func TestGetAllMappings(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	allMappings := db.GetAllMappings()

	if len(allMappings) == 0 {
		t.Error("Expected GetAllMappings to return non-empty map")
	}

	// Verify we can access specific mappings
	if _, ok := allMappings["aws_s3_bucket"]; !ok {
		t.Error("Expected aws_s3_bucket in returned mappings")
	}
}

// TestClearDatabase tests clearing the database
func TestClearDatabase(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	if !db.IsLoaded() {
		t.Error("Expected database to be loaded before clear")
	}

	db.Clear()

	if db.IsLoaded() {
		t.Error("Expected database to not be loaded after clear")
	}

	allMappings := db.GetAllMappings()
	if len(allMappings) > 0 {
		t.Errorf("Expected empty mappings after clear, got %d", len(allMappings))
	}
}

// TestS3MappingLoaded tests that S3 mappings are correctly loaded
func TestS3MappingLoaded(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Check S3 bucket mapping
	mapping, exists := db.GetMapping("aws_s3_bucket")
	if !exists {
		t.Skip("aws_s3_bucket mapping not found (may be testing without mappings)")
	}

	if mapping.Service != "s3" {
		t.Errorf("Expected s3 service, got %s", mapping.Service)
	}

	if len(mapping.Actions) == 0 {
		t.Error("Expected aws_s3_bucket to have actions")
	}

	// Check that actions were loaded
	for actionName, actionSet := range mapping.Actions {
		if actionSet.IsEmpty() {
			t.Errorf("Action %s has no IAM actions", actionName)
		}
	}
}

// TestEC2MappingLoaded tests that EC2 mappings are correctly loaded
func TestEC2MappingLoaded(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Check multiple EC2 resources
	ec2Resources := []string{
		"aws_instance",
		"aws_security_group",
		"aws_vpc",
		"aws_subnet",
	}

	for _, resourceType := range ec2Resources {
		mapping, exists := db.GetMapping(resourceType)
		if !exists {
			t.Logf("Skipping %s (mapping not found)", resourceType)
			continue
		}

		if mapping.Service != "ec2" {
			t.Errorf("Expected ec2 service for %s, got %s", resourceType, mapping.Service)
		}

		if len(mapping.Actions) == 0 {
			t.Errorf("Expected %s to have actions", resourceType)
		}
	}
}

// TestIAMMappingLoaded tests that IAM mappings are correctly loaded
func TestIAMMappingLoaded(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Check IAM resources
	iamResources := []string{
		"aws_iam_role",
		"aws_iam_policy",
		"aws_iam_user",
	}

	for _, resourceType := range iamResources {
		mapping, exists := db.GetMapping(resourceType)
		if !exists {
			t.Logf("Skipping %s (mapping not found)", resourceType)
			continue
		}

		if mapping.Service != "iam" {
			t.Errorf("Expected iam service for %s, got %s", resourceType, mapping.Service)
		}
	}
}

// TestAttributeActionsLoaded tests that attribute-specific actions are loaded
func TestAttributeActionsLoaded(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	mapping, exists := db.GetMapping("aws_s3_bucket")
	if !exists {
		t.Skip("aws_s3_bucket mapping not found")
	}

	if len(mapping.AttributeActions) == 0 {
		t.Error("Expected aws_s3_bucket to have attribute actions")
	}

	// Check versioning attribute
	if _, ok := mapping.AttributeActions["versioning"]; !ok {
		t.Error("Expected versioning attribute in aws_s3_bucket")
	}
}

// TestMappingDescriptions tests that descriptions are loaded
func TestMappingDescriptions(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	mapping, exists := db.GetMapping("aws_s3_bucket")
	if !exists {
		t.Skip("aws_s3_bucket mapping not found")
	}

	if mapping.Description == "" {
		t.Error("Expected aws_s3_bucket to have a description")
	}
}

// TestMappingThreadSafety tests that concurrent reads are safe
func TestMappingThreadSafety(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Perform concurrent reads
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			db.GetMapping("aws_s3_bucket")
			db.HasMapping("aws_s3_bucket")
			db.GetAllMappings()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestCreateTempMappingFile creates a temporary mapping file for testing
func createTempMappingFile(t *testing.T) string {
	tmpDir := t.TempDir()

	content := `
test_resource:
  service: test
  description: "Test resource"
  actions:
    create:
      - test:CreateResource
    read:
      - test:GetResource
    delete:
      - test:DeleteResource
`

	tmpFile := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp mapping file: %v", err)
	}

	return tmpDir
}

// TestLoadCustomMappingFile tests loading a custom mapping file
func TestLoadCustomMappingFile(t *testing.T) {
	tmpDir := createTempMappingFile(t)

	db := NewMappingDatabase()
	err := db.LoadMappings(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load custom mappings: %v", err)
	}

	mapping, exists := db.GetMapping("test_resource")
	if !exists {
		t.Error("Expected test_resource mapping to exist")
	}

	if mapping.Service != "test" {
		t.Errorf("Expected service test, got %s", mapping.Service)
	}

	if len(mapping.Actions) != 3 {
		t.Errorf("Expected 3 actions, got %d", len(mapping.Actions))
	}
}
