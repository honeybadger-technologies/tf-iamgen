package mapping

import (
	"testing"
)

// TestActionSetAdd tests adding single actions to a set
func TestActionSetAdd(t *testing.T) {
	as := make(ActionSet)

	as.Add("s3:GetObject")
	if !as.Contains("s3:GetObject") {
		t.Error("Expected s3:GetObject to be in ActionSet after Add")
	}

	if as.Size() != 1 {
		t.Errorf("Expected size 1, got %d", as.Size())
	}
}

// TestActionSetAddAll tests merging two action sets
func TestActionSetAddAll(t *testing.T) {
	as1 := NewActionSet("s3:GetObject", "s3:PutObject")
	as2 := NewActionSet("s3:DeleteObject", "s3:ListBucket")

	as1.AddAll(as2)

	if as1.Size() != 4 {
		t.Errorf("Expected size 4, got %d", as1.Size())
	}

	if !as1.Contains("s3:DeleteObject") {
		t.Error("Expected s3:DeleteObject after AddAll")
	}
}

// TestActionSetContains tests membership checking
func TestActionSetContains(t *testing.T) {
	as := NewActionSet("ec2:RunInstances", "ec2:TerminateInstances")

	if !as.Contains("ec2:RunInstances") {
		t.Error("Expected ec2:RunInstances to be in set")
	}

	if as.Contains("ec2:StartInstances") {
		t.Error("Expected ec2:StartInstances not to be in set")
	}
}

// TestActionSetToSlice tests conversion to slice
func TestActionSetToSlice(t *testing.T) {
	as := NewActionSet("iam:CreateRole", "iam:DeleteRole", "iam:GetRole")

	slice := as.ToSlice()

	if len(slice) != 3 {
		t.Errorf("Expected 3 actions in slice, got %d", len(slice))
	}

	// Check all actions are present
	expected := map[string]bool{
		"iam:CreateRole": false,
		"iam:DeleteRole": false,
		"iam:GetRole":    false,
	}

	for _, action := range slice {
		if _, ok := expected[action]; ok {
			expected[action] = true
		}
	}

	for action, found := range expected {
		if !found {
			t.Errorf("Expected action %s not found in slice", action)
		}
	}
}

// TestActionSetIsEmpty tests empty check
func TestActionSetIsEmpty(t *testing.T) {
	emptySet := make(ActionSet)
	if !emptySet.IsEmpty() {
		t.Error("Expected empty set to return true for IsEmpty()")
	}

	filledSet := NewActionSet("lambda:InvokeFunction")
	if filledSet.IsEmpty() {
		t.Error("Expected non-empty set to return false for IsEmpty()")
	}
}

// TestActionSetSize tests size calculation
func TestActionSetSize(t *testing.T) {
	as := NewActionSet()

	if as.Size() != 0 {
		t.Errorf("Expected size 0 for empty set, got %d", as.Size())
	}

	as.Add("rds:CreateDBInstance")
	as.Add("rds:DeleteDBInstance")

	if as.Size() != 2 {
		t.Errorf("Expected size 2, got %d", as.Size())
	}
}

// TestActionSetNoDuplicates tests that duplicates aren't added
func TestActionSetNoDuplicates(t *testing.T) {
	as := NewActionSet("dynamodb:GetItem")
	as.Add("dynamodb:GetItem")
	as.Add("dynamodb:GetItem")

	if as.Size() != 1 {
		t.Errorf("Expected size 1 (no duplicates), got %d", as.Size())
	}
}

// TestNewActionSet tests creating ActionSet with initial values
func TestNewActionSet(t *testing.T) {
	actions := []string{
		"s3:GetBucketLocation",
		"s3:ListBucket",
		"s3:GetObject",
	}

	as := NewActionSet(actions...)

	if as.Size() != 3 {
		t.Errorf("Expected size 3, got %d", as.Size())
	}

	for _, action := range actions {
		if !as.Contains(action) {
			t.Errorf("Expected %s in set", action)
		}
	}
}

// TestMappingDatabase tests database creation
func TestMappingDatabase(t *testing.T) {
	db := NewMappingDatabase()

	if db == nil {
		t.Fatal("Expected NewMappingDatabase to return non-nil database")
	}

	if db.IsLoaded() {
		t.Error("Expected newly created database to not be loaded")
	}
}

// TestMappingService tests service creation
func TestMappingService(t *testing.T) {
	db := NewMappingDatabase()
	service := NewMappingService(db)

	if service == nil {
		t.Fatal("Expected NewMappingService to return non-nil service")
	}
}

// TestResourceActions tests ResourceActions structure
func TestResourceActions(t *testing.T) {
	actions := NewActionSet("s3:GetObject", "s3:PutObject")

	ra := &ResourceActions{
		ResourceType: "aws_s3_bucket",
		ResourceName: "my-bucket",
		Service:      "s3",
		Actions:      actions,
		Reason:       "Basic bucket operations",
	}

	if ra.ResourceType != "aws_s3_bucket" {
		t.Errorf("Expected resource type aws_s3_bucket, got %s", ra.ResourceType)
	}

	if ra.Service != "s3" {
		t.Errorf("Expected service s3, got %s", ra.Service)
	}

	if !ra.Actions.Contains("s3:GetObject") {
		t.Error("Expected s3:GetObject in actions")
	}
}

// TestResourceActionMapStructure tests ResourceActionMap
func TestResourceActionMapStructure(t *testing.T) {
	mapping := &ResourceActionMap{
		Actions:          make(map[string]ActionSet),
		AttributeActions: make(map[string]map[string]ActionSet),
		Service:          "s3",
		Description:      "S3 bucket operations",
	}

	mapping.Actions["create"] = NewActionSet("s3:CreateBucket")
	mapping.Actions["read"] = NewActionSet("s3:GetBucketLocation")

	if mapping.Service != "s3" {
		t.Errorf("Expected service s3, got %s", mapping.Service)
	}

	if len(mapping.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(mapping.Actions))
	}
}
