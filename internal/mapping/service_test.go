package mapping

import (
	"testing"
)

// TestGetResourceActions tests basic resource action lookup
func TestGetResourceActions(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	actions, err := service.GetResourceActions("aws_s3_bucket", nil)
	if err != nil {
		t.Fatalf("Failed to get resource actions: %v", err)
	}

	if actions == nil {
		t.Fatal("Expected non-nil ResourceActions")
	}

	if actions.ResourceType != "aws_s3_bucket" {
		t.Errorf("Expected resource type aws_s3_bucket, got %s", actions.ResourceType)
	}

	if actions.Service != "s3" {
		t.Errorf("Expected service s3, got %s", actions.Service)
	}

	if actions.Actions.IsEmpty() {
		t.Error("Expected non-empty action set")
	}
}

// TestGetResourceActionsWithAttributes tests action lookup with attributes
func TestGetResourceActionsWithAttributes(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	// Get base actions
	baseActions, _ := service.GetResourceActions("aws_s3_bucket", nil)
	baseSize := baseActions.Actions.Size()

	// Get actions with attributes
	attrs := map[string]interface{}{
		"versioning": true,
		"logging":    true,
	}
	attrActions, _ := service.GetResourceActions("aws_s3_bucket", attrs)
	attrSize := attrActions.Actions.Size()

	// Attributes should add more actions
	if attrSize <= baseSize {
		t.Errorf("Expected more actions with attributes: base=%d, with_attrs=%d", baseSize, attrSize)
	}
}

// TestGetResourceActionsUnmappedResource tests error for unmapped resource
func TestGetResourceActionsUnmappedResource(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	_, err = service.GetResourceActions("nonexistent_resource", nil)
	if err == nil {
		t.Error("Expected error for unmapped resource")
	}
}

// TestGetActionsForMultipleResources tests bulk lookup
func TestGetActionsForMultipleResources(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	resources := map[string]map[string]interface{}{
		"aws_s3_bucket": nil,
		"aws_iam_role":  nil,
		"aws_instance":  nil,
	}

	results, err := service.GetActionsForMultipleResources(resources)
	if err != nil {
		t.Fatalf("Failed to get actions for multiple resources: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	if _, ok := results["aws_s3_bucket"]; !ok {
		t.Error("Expected aws_s3_bucket in results")
	}
}

// TestCombineActions tests merging multiple action sets
func TestCombineActions(t *testing.T) {
	service := &MappingService{
		db:    nil,
		cache: make(map[string]ActionSet),
	}

	ra1 := &ResourceActions{
		ResourceType: "aws_s3_bucket",
		Actions:      NewActionSet("s3:GetObject", "s3:PutObject"),
	}

	ra2 := &ResourceActions{
		ResourceType: "aws_iam_role",
		Actions:      NewActionSet("iam:CreateRole", "iam:DeleteRole"),
	}

	combined := service.CombineActions(ra1, ra2)

	if combined.Size() != 4 {
		t.Errorf("Expected 4 combined actions, got %d", combined.Size())
	}

	if !combined.Contains("s3:GetObject") {
		t.Error("Expected s3:GetObject in combined actions")
	}

	if !combined.Contains("iam:CreateRole") {
		t.Error("Expected iam:CreateRole in combined actions")
	}
}

// TestGetActionsByService tests grouping actions by service
func TestGetActionsByService(t *testing.T) {
	service := &MappingService{
		db:    nil,
		cache: make(map[string]ActionSet),
	}

	actions := NewActionSet(
		"s3:GetObject", "s3:PutObject",
		"ec2:RunInstances", "ec2:TerminateInstances",
		"iam:CreateRole",
	)

	grouped := service.GetActionsByService(actions)

	if len(grouped) != 3 {
		t.Errorf("Expected 3 services, got %d", len(grouped))
	}

	if _, ok := grouped["s3"]; !ok {
		t.Error("Expected s3 in grouped actions")
	}

	if grouped["s3"].Size() != 2 {
		t.Errorf("Expected 2 s3 actions, got %d", grouped["s3"].Size())
	}

	if grouped["ec2"].Size() != 2 {
		t.Errorf("Expected 2 ec2 actions, got %d", grouped["ec2"].Size())
	}

	if grouped["iam"].Size() != 1 {
		t.Errorf("Expected 1 iam action, got %d", grouped["iam"].Size())
	}
}

// TestGetResourcesWithoutMapping tests identifying unmapped resources
func TestGetResourcesWithoutMapping(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	resourceTypes := []string{
		"aws_s3_bucket",
		"nonexistent_resource1",
		"aws_iam_role",
		"nonexistent_resource2",
	}

	unmapped := service.GetResourcesWithoutMapping(resourceTypes)

	if len(unmapped) != 2 {
		t.Errorf("Expected 2 unmapped resources, got %d", len(unmapped))
	}

	foundNonexistent1 := false
	foundNonexistent2 := false

	for _, resource := range unmapped {
		if resource == "nonexistent_resource1" {
			foundNonexistent1 = true
		}
		if resource == "nonexistent_resource2" {
			foundNonexistent2 = true
		}
	}

	if !foundNonexistent1 {
		t.Error("Expected nonexistent_resource1 in unmapped")
	}

	if !foundNonexistent2 {
		t.Error("Expected nonexistent_resource2 in unmapped")
	}
}

// TestGetCoverageStats tests coverage statistics
func TestGetCoverageStats(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	stats := service.GetCoverageStats()

	if stats == nil {
		t.Fatal("Expected non-nil coverage stats")
	}

	totalMappings, ok := stats["total_mappings"]
	if !ok {
		t.Error("Expected total_mappings in stats")
	}

	if totalMappings.(int) == 0 {
		t.Error("Expected non-zero total mappings")
	}

	totalActions, ok := stats["total_actions"]
	if !ok {
		t.Error("Expected total_actions in stats")
	}

	if totalActions.(int) == 0 {
		t.Error("Expected non-zero total actions")
	}

	services, ok := stats["services"]
	if !ok {
		t.Error("Expected services in stats")
	}

	if len(services.(map[string]int)) == 0 {
		t.Error("Expected non-empty services map")
	}
}

// TestGetMappingInfo tests detailed mapping information
func TestGetMappingInfo(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	info, err := service.GetMappingInfo("aws_s3_bucket")
	if err != nil {
		t.Fatalf("Failed to get mapping info: %v", err)
	}

	if info == nil {
		t.Fatal("Expected non-nil mapping info")
	}

	if info["resource_type"] != "aws_s3_bucket" {
		t.Errorf("Expected resource_type aws_s3_bucket, got %v", info["resource_type"])
	}

	if info["service"] != "s3" {
		t.Errorf("Expected service s3, got %v", info["service"])
	}

	if baseActionCount, ok := info["base_action_count"]; ok {
		if baseActionCount.(int) == 0 {
			t.Error("Expected non-zero base action count")
		}
	}
}

// TestGetMappingInfoUnmapped tests error for unmapped resource info
func TestGetMappingInfoUnmapped(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	_, err = service.GetMappingInfo("nonexistent_resource")
	if err == nil {
		t.Error("Expected error for unmapped resource info")
	}
}

// TestClearCache tests cache clearing
func TestClearCache(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	// Populate cache
	service.GetResourceActions("aws_s3_bucket", nil)

	if len(service.cache) == 0 {
		t.Skip("Cache not populated (may be expected)")
	}

	service.ClearCache()

	if len(service.cache) > 0 {
		t.Errorf("Expected empty cache after clear, got %d entries", len(service.cache))
	}
}

// TestCaching tests that caching works correctly
func TestCaching(t *testing.T) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	// First call
	actions1, _ := service.GetResourceActions("aws_s3_bucket", nil)

	// Cache should have entry
	initialCacheSize := len(service.cache)

	// Second call with same parameters
	actions2, _ := service.GetResourceActions("aws_s3_bucket", nil)

	// Cache size shouldn't increase
	if len(service.cache) != initialCacheSize {
		t.Errorf("Expected cache size to stay %d, got %d", initialCacheSize, len(service.cache))
	}

	// Results should be same
	if actions1.Actions.Size() != actions2.Actions.Size() {
		t.Error("Expected cached result to be identical")
	}
}

// TestBenchmarkGetResourceActions benchmarks lookup performance
func BenchmarkGetResourceActions(b *testing.B) {
	db := NewMappingDatabase()
	err := db.LoadMappings("../../mappings")
	if err != nil {
		b.Fatalf("Failed to load mappings: %v", err)
	}

	service := NewMappingService(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetResourceActions("aws_s3_bucket", nil)
	}
}

// TestBenchmarkCombineActions benchmarks combining actions
func BenchmarkCombineActions(b *testing.B) {
	service := &MappingService{
		db:    nil,
		cache: make(map[string]ActionSet),
	}

	ra1 := &ResourceActions{
		ResourceType: "aws_s3_bucket",
		Actions:      NewActionSet("s3:GetObject", "s3:PutObject"),
	}

	ra2 := &ResourceActions{
		ResourceType: "aws_iam_role",
		Actions:      NewActionSet("iam:CreateRole", "iam:DeleteRole"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CombineActions(ra1, ra2)
	}
}

// TestBenchmarkGroupByService benchmarks grouping actions
func BenchmarkGroupByService(b *testing.B) {
	service := &MappingService{
		db:    nil,
		cache: make(map[string]ActionSet),
	}

	actions := NewActionSet(
		"s3:GetObject", "s3:PutObject",
		"ec2:RunInstances", "ec2:TerminateInstances",
		"iam:CreateRole", "iam:DeleteRole",
		"lambda:InvokeFunction",
		"rds:CreateDBInstance",
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetActionsByService(actions)
	}
}
