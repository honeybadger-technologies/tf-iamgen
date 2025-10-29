package policy

import (
	"encoding/json"
	"sort"
	"strings"
)

// PolicyVersion is the IAM policy language version
const PolicyVersion = "2012-10-17"

// Effect is the effect of a policy statement
type Effect string

const (
	EffectAllow Effect = "Allow"
	EffectDeny  Effect = "Deny"
)

// Principal represents who the policy applies to
type Principal struct {
	Service   []string            `json:"Service,omitempty"`
	AWS       []string            `json:"AWS,omitempty"`
	Principal map[string][]string `json:"Principal,omitempty"` // For federated principals
}

// Statement represents a single IAM policy statement
type Statement struct {
	Sid       string      `json:"Sid,omitempty"`
	Effect    Effect      `json:"Effect"`
	Principal *Principal  `json:"Principal,omitempty"`
	Action    []string    `json:"Action"`
	Resource  []string    `json:"Resource"`
	Condition interface{} `json:"Condition,omitempty"`
}

// Policy represents a complete IAM policy document
type Policy struct {
	Version   string      `json:"Version"`
	Sid       string      `json:"Sid,omitempty"`
	Statement []Statement `json:"Statement"`
}

// PolicyMetadata contains metadata about generated policies
type PolicyMetadata struct {
	GeneratedAt      string   // ISO 8601 timestamp
	TerraformVersion string   // Version of Terraform analyzed
	ResourceCount    int      // Number of resources analyzed
	ActionCount      int      // Number of IAM actions
	Services         []string // AWS services used
	Checksum         string   // Hash of policy for validation
}

// PolicyGenerationOptions controls policy generation behavior
type PolicyGenerationOptions struct {
	// Include wildcard resources or specific ARNs
	UseWildcardResources bool

	// Group statements by service or resource
	GroupBy string // "service" or "resource"

	// Include resource conditions
	IncludeConditions bool

	// Add statement IDs for clarity
	IncludeSids bool

	// Minimize policy size
	Minimize bool

	// Custom resource ARN mappings
	ResourceMappings map[string]string
}

// NewPolicy creates a new empty policy
func NewPolicy() *Policy {
	return &Policy{
		Version:   PolicyVersion,
		Statement: []Statement{},
	}
}

// AddStatement adds a new statement to the policy
func (p *Policy) AddStatement(statement Statement) {
	// Ensure actions and resources are sorted for consistency
	sort.Strings(statement.Action)
	sort.Strings(statement.Resource)

	p.Statement = append(p.Statement, statement)
}

// ToJSON converts the policy to JSON with proper formatting
func (p *Policy) ToJSON() (string, error) {
	// Sort statements for consistency
	sortStatements(p.Statement)

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToCompactJSON converts the policy to minified JSON
func (p *Policy) ToCompactJSON() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MergeActions merges multiple action slices removing duplicates
func MergeActions(actionSlices ...[]string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, actions := range actionSlices {
		for _, action := range actions {
			if !seen[action] {
				result = append(result, action)
				seen[action] = true
			}
		}
	}

	sort.Strings(result)
	return result
}

// MergeResources merges multiple resource slices removing duplicates
func MergeResources(resourceSlices ...[]string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, resources := range resourceSlices {
		for _, resource := range resources {
			if !seen[resource] {
				result = append(result, resource)
				seen[resource] = true
			}
		}
	}

	sort.Strings(result)
	return result
}

// sortStatements sorts statements for consistent output
func sortStatements(statements []Statement) {
	sort.Slice(statements, func(i, j int) bool {
		// Sort by Sid first, then by first action
		if statements[i].Sid != statements[j].Sid {
			return statements[i].Sid < statements[j].Sid
		}
		if len(statements[i].Action) > 0 && len(statements[j].Action) > 0 {
			return statements[i].Action[0] < statements[j].Action[0]
		}
		return i < j
	})
}

// PolicyBuilder helps construct policies step by step
type PolicyBuilder struct {
	policy   *Policy
	options  PolicyGenerationOptions
	metadata PolicyMetadata
}

// NewPolicyBuilder creates a new policy builder with options
func NewPolicyBuilder(opts PolicyGenerationOptions) *PolicyBuilder {
	return &PolicyBuilder{
		policy:  NewPolicy(),
		options: opts,
	}
}

// AddActionStatement adds a statement for specific actions and resources
func (pb *PolicyBuilder) AddActionStatement(sid string, actions []string, resources []string) {
	statement := Statement{
		Sid:      sid,
		Effect:   EffectAllow,
		Action:   MergeActions([][]string{actions}...),
		Resource: MergeResources([][]string{resources}...),
	}

	pb.policy.AddStatement(statement)
}

// Build returns the constructed policy and metadata
func (pb *PolicyBuilder) Build() (*Policy, PolicyMetadata) {
	return pb.policy, pb.metadata
}

// GetPolicy returns the current policy
func (pb *PolicyBuilder) GetPolicy() *Policy {
	return pb.policy
}

// SetMetadata sets the policy metadata
func (pb *PolicyBuilder) SetMetadata(metadata PolicyMetadata) {
	pb.metadata = metadata
}

// GroupStatementsByService groups statements by AWS service
func GroupStatementsByService(statements []Statement) map[string][]Statement {
	groups := make(map[string][]Statement)

	for _, stmt := range statements {
		for _, action := range stmt.Action {
			parts := strings.Split(action, ":")
			if len(parts) >= 1 {
				service := parts[0]
				groups[service] = append(groups[service], stmt)
				break // Only group by first service
			}
		}
	}

	return groups
}

// GetServicesFromStatements extracts unique services from statements
func GetServicesFromStatements(statements []Statement) []string {
	services := make(map[string]bool)

	for _, stmt := range statements {
		for _, action := range stmt.Action {
			parts := strings.Split(action, ":")
			if len(parts) >= 1 {
				services[parts[0]] = true
			}
		}
	}

	var result []string
	for service := range services {
		result = append(result, service)
	}
	sort.Strings(result)
	return result
}

// ActionToResource converts an action to a resource ARN pattern
func ActionToResource(action string) string {
	// Default to wildcard for most cases
	// Can be overridden with ResourceMappings
	return "*"
}

// NormalizeAction ensures action is in correct format (service:Action)
func NormalizeAction(action string) string {
	// Ensure action contains service prefix
	if !strings.Contains(action, ":") {
		return "*:" + action
	}
	return action
}
