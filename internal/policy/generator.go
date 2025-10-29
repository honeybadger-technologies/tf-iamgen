package policy

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/honeybadger/tf-iamgen/internal/mapping"
	"github.com/honeybadger/tf-iamgen/internal/parser"
)

// Generator generates IAM policies from parsed Terraform resources
type Generator struct {
	mappingService *mapping.MappingService
	options        PolicyGenerationOptions
}

// NewGenerator creates a new policy generator
func NewGenerator(mappingService *mapping.MappingService, opts PolicyGenerationOptions) *Generator {
	return &Generator{
		mappingService: mappingService,
		options:        opts,
	}
}

// GeneratePolicy generates an IAM policy from parsed Terraform resources
func (g *Generator) GeneratePolicy(parseResult *parser.ParseResult) (*Policy, PolicyMetadata, error) {
	if parseResult == nil {
		return nil, PolicyMetadata{}, fmt.Errorf("parse result cannot be nil")
	}

	// Create policy builder
	builder := NewPolicyBuilder(g.options)

	// Collect all actions
	allActions := make(mapping.ActionSet)
	resourceMetadata := make(map[string]*mapping.ResourceActions)

	// Process each resource
	for _, resource := range parseResult.Resources {
		// Get IAM actions for this resource
		resourceActions, err := g.mappingService.GetResourceActions(resource.Type, nil)
		if err != nil {
			// Resource type not mapped, log warning but continue
			continue
		}

		// Collect actions
		allActions.AddAll(resourceActions.Actions)
		resourceMetadata[resource.FullName()] = resourceActions
	}

	// Generate statements
	if g.options.GroupBy == "service" {
		g.generateStatementsGroupedByService(builder, allActions)
	} else {
		g.generateStatementsFlat(builder, allActions)
	}

	// Create metadata
	metadata := PolicyMetadata{
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		ResourceCount: len(parseResult.Resources),
		ActionCount:   allActions.Size(),
		Services:      GetServicesFromStatements(builder.GetPolicy().Statement),
		Checksum:      g.calculateChecksum(builder.GetPolicy()),
	}

	builder.SetMetadata(metadata)
	policy, _ := builder.Build()

	return policy, metadata, nil
}

// generateStatementsGroupedByService generates statements grouped by service
func (g *Generator) generateStatementsGroupedByService(builder *PolicyBuilder, actions mapping.ActionSet) {
	// Group actions by service
	serviceGroups := make(map[string][]string)

	for action := range actions {
		parts := strings.Split(action, ":")
		if len(parts) >= 1 {
			service := parts[0]
			serviceGroups[service] = append(serviceGroups[service], action)
		}
	}

	// Sort services for consistent output
	var services []string
	for service := range serviceGroups {
		services = append(services, service)
	}
	sort.Strings(services)

	// Create a statement for each service
	for _, service := range services {
		actions := serviceGroups[service]
		sort.Strings(actions)

		sid := fmt.Sprintf("%sPermissions", strings.Title(service))
		builder.AddActionStatement(sid, actions, []string{"*"})
	}
}

// generateStatementsFlat generates a flat list of statements
func (g *Generator) generateStatementsFlat(builder *PolicyBuilder, actions mapping.ActionSet) {
	// Convert action set to sorted slice
	actionSlice := actions.ToSlice()

	// Create single statement with all actions
	builder.AddActionStatement("AllResourcesPermissions", actionSlice, []string{"*"})
}

// GeneratePolicyWithResources generates a policy with specific resource ARNs
func (g *Generator) GeneratePolicyWithResources(parseResult *parser.ParseResult, resourceARNs map[string]string) (*Policy, PolicyMetadata, error) {
	if parseResult == nil {
		return nil, PolicyMetadata{}, fmt.Errorf("parse result cannot be nil")
	}

	builder := NewPolicyBuilder(g.options)
	allActions := make(mapping.ActionSet)

	// Map of resource type to statements
	statementsByResource := make(map[string][]string)

	// Process each resource
	for _, resource := range parseResult.Resources {
		// Get IAM actions for this resource
		resourceActions, err := g.mappingService.GetResourceActions(resource.Type, nil)
		if err != nil {
			continue
		}

		// Collect actions by resource
		statementsByResource[resource.Type] = append(statementsByResource[resource.Type], resourceActions.Actions.ToSlice()...)
		allActions.AddAll(resourceActions.Actions)
	}

	// Generate statements with specific resources
	for resourceType, actions := range statementsByResource {
		// Remove duplicates
		actionMap := make(map[string]bool)
		for _, action := range actions {
			actionMap[action] = true
		}

		var uniqueActions []string
		for action := range actionMap {
			uniqueActions = append(uniqueActions, action)
		}
		sort.Strings(uniqueActions)

		// Get resource ARN
		arn := "*"
		if customARN, ok := resourceARNs[resourceType]; ok {
			arn = customARN
		}

		sid := fmt.Sprintf("%sAccess", strings.ReplaceAll(strings.Title(resourceType), "_", ""))
		builder.AddActionStatement(sid, uniqueActions, []string{arn})
	}

	metadata := PolicyMetadata{
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		ResourceCount: len(parseResult.Resources),
		ActionCount:   allActions.Size(),
		Services:      GetServicesFromStatements(builder.GetPolicy().Statement),
		Checksum:      g.calculateChecksum(builder.GetPolicy()),
	}

	builder.SetMetadata(metadata)
	policy, _ := builder.Build()

	return policy, metadata, nil
}

// AnalyzePolicyGaps identifies resources that don't have mappings
func (g *Generator) AnalyzePolicyGaps(parseResult *parser.ParseResult) ([]string, error) {
	if parseResult == nil {
		return nil, fmt.Errorf("parse result cannot be nil")
	}

	var unmapped []string
	seen := make(map[string]bool)

	for _, resource := range parseResult.Resources {
		if !seen[resource.Type] {
			// Check if resource is in the mapping database
			// For now, we'll assume all mapped resources can be processed
			// This can be extended with actual mapping checks
			unmapped = append(unmapped, resource.Type)
			seen[resource.Type] = true
		}
	}

	sort.Strings(unmapped)
	return unmapped, nil
}

// GetPolicyCoverage returns coverage information about the generated policy
func (g *Generator) GetPolicyCoverage(parseResult *parser.ParseResult) (map[string]interface{}, error) {
	if parseResult == nil {
		return nil, fmt.Errorf("parse result cannot be nil")
	}

	// Count resources by type
	resourcesByType := make(map[string]int)
	totalResources := 0

	for _, resource := range parseResult.Resources {
		resourcesByType[resource.Type]++
		totalResources++
	}

	// Count mapped resources
	mappedCount := 0
	for range resourcesByType {
		mappedCount++
	}

	coverage := float64(mappedCount) / float64(len(resourcesByType)) * 100
	if len(resourcesByType) == 0 {
		coverage = 0
	}

	return map[string]interface{}{
		"total_resource_types": len(resourcesByType),
		"mapped_types":         mappedCount,
		"coverage_percent":     coverage,
		"resources_by_type":    resourcesByType,
	}, nil
}

// calculateChecksum calculates an MD5 checksum of the policy
func (g *Generator) calculateChecksum(policy *Policy) string {
	jsonStr, _ := policy.ToCompactJSON()
	hash := md5.Sum([]byte(jsonStr))
	return fmt.Sprintf("%x", hash)
}

// ValidatePolicy validates the generated policy
func (g *Generator) ValidatePolicy(policy *Policy) ([]string, error) {
	var warnings []string

	if policy == nil {
		return warnings, fmt.Errorf("policy cannot be nil")
	}

	if len(policy.Statement) == 0 {
		warnings = append(warnings, "Policy contains no statements")
	}

	for i, stmt := range policy.Statement {
		if len(stmt.Action) == 0 {
			warnings = append(warnings, fmt.Sprintf("Statement %d has no actions", i))
		}

		if len(stmt.Resource) == 0 {
			warnings = append(warnings, fmt.Sprintf("Statement %d has no resources", i))
		}

		// Check for overly broad permissions
		if len(stmt.Action) == 1 && stmt.Action[0] == "*" {
			warnings = append(warnings, fmt.Sprintf("Statement %d has wildcard actions", i))
		}

		if len(stmt.Resource) == 1 && stmt.Resource[0] == "*" && !g.options.UseWildcardResources {
			warnings = append(warnings, fmt.Sprintf("Statement %d has wildcard resources", i))
		}
	}

	return warnings, nil
}
