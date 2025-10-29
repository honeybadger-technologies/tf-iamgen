package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProviderSpec represents metadata about a Terraform provider
type ProviderSpec struct {
	Provider    string                   `json:"provider"`
	Version     string                   `json:"version"`
	Resources   map[string]*ResourceSpec `json:"resources"`
	DataSources map[string]*ResourceSpec `json:"data_sources"`
}

// ResourceSpec represents a Terraform resource's schema and required operations
type ResourceSpec struct {
	Name             string              `json:"name"`
	CreateOperations []string            `json:"create_operations"` // API calls for Create
	ReadOperations   []string            `json:"read_operations"`   // API calls for Read
	UpdateOperations []string            `json:"update_operations"` // API calls for Update
	DeleteOperations []string            `json:"delete_operations"` // API calls for Delete
	ListOperations   []string            `json:"list_operations"`   // API calls for List
	Arguments        map[string]Argument `json:"arguments"`
}

// Argument represents a resource argument with its required API operations
type Argument struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Required           bool     `json:"required"`
	RequiredOperations []string `json:"required_operations"` // Operations needed for this argument
}

// ProviderSpecLoader loads provider specifications
type ProviderSpecLoader struct {
	specDir string
	cache   map[string]*ProviderSpec
}

// NewProviderSpecLoader creates a new provider spec loader
func NewProviderSpecLoader(specDir string) *ProviderSpecLoader {
	return &ProviderSpecLoader{
		specDir: specDir,
		cache:   make(map[string]*ProviderSpec),
	}
}

// LoadProviderSpec loads a provider specification
func (psl *ProviderSpecLoader) LoadProviderSpec(provider string, version string) (*ProviderSpec, error) {
	cacheKey := fmt.Sprintf("%s:%s", provider, version)

	if spec, cached := psl.cache[cacheKey]; cached {
		return spec, nil
	}

	specPath := filepath.Join(psl.specDir, provider, fmt.Sprintf("%s-spec.json", version))

	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load provider spec: %w", err)
	}

	var spec ProviderSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse provider spec: %w", err)
	}

	psl.cache[cacheKey] = &spec
	return &spec, nil
}

// GetResourceOperations returns the operations needed for a resource type
func (spec *ProviderSpec) GetResourceOperations(resourceType string, operation string) ([]string, error) {
	resource, exists := spec.Resources[resourceType]
	if !exists {
		return nil, fmt.Errorf("resource %s not found in provider %s", resourceType, spec.Provider)
	}

	switch strings.ToLower(operation) {
	case "create":
		return resource.CreateOperations, nil
	case "read", "get":
		return resource.ReadOperations, nil
	case "update":
		return resource.UpdateOperations, nil
	case "delete":
		return resource.DeleteOperations, nil
	case "list":
		return resource.ListOperations, nil
	case "all":
		return mergeOperations(
			resource.CreateOperations,
			resource.ReadOperations,
			resource.UpdateOperations,
			resource.DeleteOperations,
			resource.ListOperations,
		), nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

// GetArgumentOperations returns operations needed for a specific argument
func (spec *ProviderSpec) GetArgumentOperations(resourceType string, argumentName string) ([]string, error) {
	resource, exists := spec.Resources[resourceType]
	if !exists {
		return nil, fmt.Errorf("resource %s not found", resourceType)
	}

	arg, exists := resource.Arguments[argumentName]
	if !exists {
		return nil, fmt.Errorf("argument %s not found in resource %s", argumentName, resourceType)
	}

	return arg.RequiredOperations, nil
}

// GetResourceByOperation filters resources that use a specific operation
func (spec *ProviderSpec) GetResourcesByOperation(operation string) []string {
	var resources []string

	for resourceName, resource := range spec.Resources {
		var ops []string

		switch strings.ToLower(operation) {
		case "create":
			ops = resource.CreateOperations
		case "read":
			ops = resource.ReadOperations
		case "update":
			ops = resource.UpdateOperations
		case "delete":
			ops = resource.DeleteOperations
		case "list":
			ops = resource.ListOperations
		}

		if len(ops) > 0 {
			resources = append(resources, resourceName)
		}
	}

	return resources
}

// mergeOperations merges multiple operation slices
func mergeOperations(opSlices ...[]string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, ops := range opSlices {
		for _, op := range ops {
			if !seen[op] {
				result = append(result, op)
				seen[op] = true
			}
		}
	}

	return result
}

// PermissionNarrower narrows permissions based on provider specs and runtime analysis
type PermissionNarrower struct {
	providerSpecs map[string]*ProviderSpec
	actionHistory map[string][]string // Track which actions are actually used
}

// NewPermissionNarrower creates a new permission narrower
func NewPermissionNarrower() *PermissionNarrower {
	return &PermissionNarrower{
		providerSpecs: make(map[string]*ProviderSpec),
		actionHistory: make(map[string][]string),
	}
}

// RegisterProviderSpec registers a provider specification
func (pn *PermissionNarrower) RegisterProviderSpec(spec *ProviderSpec) {
	pn.providerSpecs[spec.Provider] = spec
}

// NarrowPermissions narrows permissions based on provider operations
func (pn *PermissionNarrower) NarrowPermissions(resourceType string, allActions []string) []string {
	// Extract service from resource type (e.g., aws_s3_bucket -> s3)
	parts := strings.Split(resourceType, "_")
	if len(parts) < 2 {
		return allActions
	}

	service := parts[1]
	if service == "db" && len(parts) > 2 {
		service = parts[1] + parts[2] // e.g., aws_db_instance -> rds
	}

	// For now, return all actions
	// In the future, this would:
	// 1. Look up provider spec for resource type
	// 2. Extract operations used
	// 3. Match those to specific IAM actions
	// 4. Return only necessary actions

	return allActions
}

// RecordActionUsage records that an action was used
func (pn *PermissionNarrower) RecordActionUsage(resourceType string, action string) {
	pn.actionHistory[resourceType] = append(pn.actionHistory[resourceType], action)
}

// GetUsedActions returns actions that were actually used
func (pn *PermissionNarrower) GetUsedActions(resourceType string) []string {
	return pn.actionHistory[resourceType]
}

// GetUnusedActions returns actions that were not used
func (pn *PermissionNarrower) GetUnusedActions(resourceType string, allActions []string) []string {
	usedMap := make(map[string]bool)
	for _, action := range pn.actionHistory[resourceType] {
		usedMap[action] = true
	}

	var unused []string
	for _, action := range allActions {
		if !usedMap[action] {
			unused = append(unused, action)
		}
	}

	return unused
}

// LeastPrivilegePolicy represents a policy optimized for least privilege
type LeastPrivilegePolicy struct {
	BasePolicy         *Policy
	NarrowedPolicy     *Policy
	RemovedPermissions map[string][]string
	CoveragePercent    float64
}

// CreateLeastPrivilegePolicy creates a least-privilege variant of a policy
func (pn *PermissionNarrower) CreateLeastPrivilegePolicy(basePolicy *Policy) *LeastPrivilegePolicy {
	lpp := &LeastPrivilegePolicy{
		BasePolicy:         basePolicy,
		NarrowedPolicy:     NewPolicy(),
		RemovedPermissions: make(map[string][]string),
	}

	removedCount := 0
	totalCount := 0

	// Copy statements but narrow actions based on used actions
	for _, stmt := range basePolicy.Statement {
		narrowedActions := []string{}

		for _, action := range stmt.Action {
			// This is where we'd check against actionHistory
			// For now, keep all actions
			narrowedActions = append(narrowedActions, action)
		}

		totalCount += len(stmt.Action)
		removedCount += len(stmt.Action) - len(narrowedActions)

		if len(narrowedActions) > 0 {
			narrowedStmt := stmt
			narrowedStmt.Action = narrowedActions
			lpp.NarrowedPolicy.AddStatement(narrowedStmt)
		}
	}

	if totalCount > 0 {
		lpp.CoveragePercent = float64(totalCount-removedCount) / float64(totalCount) * 100
	}

	return lpp
}
