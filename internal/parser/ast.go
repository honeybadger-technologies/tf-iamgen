// Package parser provides Terraform HCL file parsing functionality.
package parser

import (
	"fmt"
	"strings"
)

// Resource represents a single AWS Terraform resource.
type Resource struct {
	Type       string                 // Resource type (e.g., "aws_s3_bucket")
	Name       string                 // Resource name (e.g., "my_bucket")
	Attributes map[string]interface{} // Resource attributes and their values
	FilePath   string                 // Path to the file where this resource is defined
	LineNumber int                    // Line number where resource starts
}

// String returns a formatted string representation of a resource.
func (r *Resource) String() string {
	return fmt.Sprintf("resource \"%s\" \"%s\" (%s:%d)", r.Type, r.Name, r.FilePath, r.LineNumber)
}

// FullName returns the fully qualified resource name.
func (r *Resource) FullName() string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// Block represents a configuration block (resource, variable, data, etc.)
type Block struct {
	Type       string   // Block type (e.g., "resource", "variable", "data", "module")
	Labels     []string // Labels (e.g., ["aws_s3_bucket", "my_bucket"])
	Attributes map[string]interface{}
	Blocks     []Block // Nested blocks
	FilePath   string
	LineNumber int
}

// Attribute represents a configuration attribute with its value.
type Attribute struct {
	Name  string
	Value interface{}
	Type  string // "string", "number", "bool", "list", "map", "object", "reference"
}

// ParseResult contains all resources and metadata from parsing.
type ParseResult struct {
	Resources      []Resource             // Discovered AWS resources
	Variables      map[string]*Block      // Declared variables
	Modules        map[string]*Block      // Module declarations
	LocalValues    map[string]interface{} // Local values
	FilesProcessed int                    // Number of files parsed
	TotalResources int                    // Total resources found
	Errors         []ParseError           // Any errors encountered during parsing
}

// ParseError represents an error during parsing.
type ParseError struct {
	FilePath  string
	Line      int
	Column    int
	Message   string
	ErrorType string // "syntax", "file", "unknown_resource_type", etc.
}

// Error returns the string representation of a ParseError.
func (e ParseError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("%s:%d:%d: %s (%s)", e.FilePath, e.Line, e.Column, e.Message, e.ErrorType)
	}
	if e.Line > 0 {
		return fmt.Sprintf("%s:%d: %s (%s)", e.FilePath, e.Line, e.Message, e.ErrorType)
	}
	return fmt.Sprintf("%s: %s (%s)", e.FilePath, e.Message, e.ErrorType)
}

// HasErrors returns true if there are any critical errors.
func (pr *ParseResult) HasErrors() bool {
	for _, err := range pr.Errors {
		if err.ErrorType == "syntax" || err.ErrorType == "file" {
			return true
		}
	}
	return false
}

// GetResourcesByType returns all resources matching the given type.
func (pr *ParseResult) GetResourcesByType(resourceType string) []Resource {
	var result []Resource
	for _, res := range pr.Resources {
		if res.Type == resourceType {
			result = append(result, res)
		}
	}
	return result
}

// GetResourcesByService returns all resources for a given AWS service.
// e.g., "s3", "ec2", "rds", "iam"
func (pr *ParseResult) GetResourcesByService(service string) []Resource {
	var result []Resource
	prefix := "aws_" + service + "_"
	for _, res := range pr.Resources {
		if strings.HasPrefix(res.Type, prefix) {
			result = append(result, res)
		}
	}
	return result
}

// Summary returns a human-readable summary of the parse result.
func (pr *ParseResult) Summary() string {
	serviceMap := make(map[string]int)
	for _, res := range pr.Resources {
		parts := strings.Split(res.Type, "_")
		if len(parts) >= 2 {
			service := parts[1]
			serviceMap[service]++
		}
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Parsed %d files", pr.FilesProcessed))
	lines = append(lines, fmt.Sprintf("Found %d resources:", pr.TotalResources))

	for service, count := range serviceMap {
		lines = append(lines, fmt.Sprintf("  - %s: %d resources", service, count))
	}

	if len(pr.Errors) > 0 {
		lines = append(lines, fmt.Sprintf("Warnings: %d", len(pr.Errors)))
	}

	return strings.Join(lines, "\n")
}
