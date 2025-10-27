package mapping

import (
	"sync"
)

// ActionSet represents a set of IAM actions
type ActionSet map[string]bool

// ResourceActionMap maps AWS resource types to their required IAM actions
type ResourceActionMap struct {
	// ResourceType -> Actions mapping (e.g., "aws_s3_bucket" -> [s3:GetObject, s3:PutObject])
	Actions map[string]ActionSet `yaml:"actions"`

	// Resource attributes that determine which actions are needed
	// e.g., aws_s3_bucket with "versioning" attribute needs s3:ListBucketVersions
	AttributeActions map[string]map[string]ActionSet `yaml:"attribute_actions"`

	// Service this resource belongs to (e.g., "s3", "ec2", "iam")
	Service string `yaml:"service"`

	// Description of what permissions this resource typically needs
	Description string `yaml:"description"`
}

// MappingDatabase holds all resource-to-action mappings with caching
type MappingDatabase struct {
	mu       sync.RWMutex
	mappings map[string]*ResourceActionMap
	loaded   bool
}

// MappingService provides lookup functionality for IAM mappings
type MappingService struct {
	db    *MappingDatabase
	cache map[string]ActionSet // Simple cache for computed action sets
	mu    sync.RWMutex
}

// ResourceActions represents all IAM actions needed for a specific resource
type ResourceActions struct {
	ResourceType string    // e.g., "aws_s3_bucket"
	ResourceName string    // e.g., "my_bucket"
	Service      string    // e.g., "s3"
	Actions      ActionSet // Set of required IAM actions
	Reason       string    // Why these actions are needed
}

// PolicyStatement represents a single IAM policy statement
type PolicyStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

// IAMPolicy represents a complete IAM policy
type IAMPolicy struct {
	Version   string            `json:"Version"`
	Statement []PolicyStatement `json:"Statement"`
}

// NewActionSet creates a new ActionSet from a slice of action strings
func NewActionSet(actions ...string) ActionSet {
	set := make(ActionSet)
	for _, action := range actions {
		set[action] = true
	}
	return set
}

// Add adds an action to the set
func (as ActionSet) Add(action string) {
	as[action] = true
}

// AddAll adds all actions from another ActionSet
func (as ActionSet) AddAll(other ActionSet) {
	for action := range other {
		as[action] = true
	}
}

// Contains checks if an action exists in the set
func (as ActionSet) Contains(action string) bool {
	return as[action]
}

// ToSlice converts ActionSet to a sorted slice
func (as ActionSet) ToSlice() []string {
	result := make([]string, 0, len(as))
	for action := range as {
		result = append(result, action)
	}
	return result
}

// IsEmpty checks if the ActionSet is empty
func (as ActionSet) IsEmpty() bool {
	return len(as) == 0
}

// Size returns the number of actions in the set
func (as ActionSet) Size() int {
	return len(as)
}

// NewMappingDatabase creates a new mapping database
func NewMappingDatabase() *MappingDatabase {
	return &MappingDatabase{
		mappings: make(map[string]*ResourceActionMap),
		loaded:   false,
	}
}

// NewMappingService creates a new mapping service
func NewMappingService(db *MappingDatabase) *MappingService {
	return &MappingService{
		db:    db,
		cache: make(map[string]ActionSet),
	}
}
