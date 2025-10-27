package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

// TerraformParser is the main parser for Terraform HCL files
type TerraformParser struct {
	hclParser *hclparse.Parser
	files     []string
	result    *ParseResult
}

// NewTerraformParser creates a new Terraform HCL parser
func NewTerraformParser() *TerraformParser {
	return &TerraformParser{
		hclParser: hclparse.NewParser(),
		result: &ParseResult{
			Resources:   []Resource{},
			Variables:   make(map[string]*Block),
			Modules:     make(map[string]*Block),
			LocalValues: make(map[string]interface{}),
			Errors:      []ParseError{},
		},
	}
}

// ParseDirectory parses all Terraform files in a directory recursively
func (tp *TerraformParser) ParseDirectory(dirPath string) (*ParseResult, error) {
	// Normalize path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if directory exists
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("directory not found: %w", err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Find all .tf and .tf.json files
	files, err := tp.findTerraformFiles(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find terraform files: %w", err)
	}

	if len(files) == 0 {
		return tp.result, nil // Empty directory is valid
	}

	tp.files = files

	// Parse each file
	for _, filePath := range files {
		if err := tp.parseFile(filePath); err != nil {
			// Add warning but continue parsing other files
			tp.result.Errors = append(tp.result.Errors, ParseError{
				FilePath:  filePath,
				Message:   err.Error(),
				ErrorType: "parse_error",
			})
		}
	}

	tp.result.FilesProcessed = len(files)
	tp.result.TotalResources = len(tp.result.Resources)

	return tp.result, nil
}

// ParseFile parses a single Terraform file
func (tp *TerraformParser) parseFile(filePath string) error {
	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Determine file type
	isTFJSON := strings.HasSuffix(filePath, ".tf.json") || strings.HasSuffix(filePath, ".json")

	var file *hcl.File
	var diags hcl.Diagnostics

	if isTFJSON {
		// Parse JSON format
		file, diags = tp.hclParser.ParseJSON(content, filePath)
	} else {
		// Parse HCL format
		file, diags = tp.hclParser.ParseHCL(content, filePath)
	}

	// Check for parse errors
	if diags.HasErrors() {
		var errMessages []string
		for _, diag := range diags {
			errMessages = append(errMessages, diag.Error())
		}
		return fmt.Errorf("parse errors: %s", strings.Join(errMessages, "; "))
	}

	// Extract resources from the file
	tp.extractResources(file.Body, filePath)

	return nil
}

// extractResources extracts all resources from an HCL body
func (tp *TerraformParser) extractResources(body hcl.Body, filePath string) {
	if body == nil {
		return
	}

	// Get the schema for the body
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
			{
				Type:       "local",
				LabelNames: []string{},
			},
			{
				Type:       "terraform",
				LabelNames: []string{},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}

	// Parse the body content
	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		for _, diag := range diags {
			tp.result.Errors = append(tp.result.Errors, ParseError{
				FilePath:  filePath,
				Message:   diag.Error(),
				ErrorType: "parse_error",
			})
		}
		return
	}

	// Process resource blocks
	for _, resourceBlock := range content.Blocks {
		if resourceBlock.Type != "resource" {
			continue
		}

		// Resource blocks have two labels: type and name
		if len(resourceBlock.Labels) < 2 {
			continue
		}

		resourceType := resourceBlock.Labels[0]
		resourceName := resourceBlock.Labels[1]

		// Only extract AWS resources
		if !IsAWSResource(resourceType) {
			continue
		}

		// Parse the resource attributes
		attrs, _ := resourceBlock.Body.JustAttributes()
		attributes := tp.extractAttributes(attrs)

		// Create Resource struct
		resource := Resource{
			Type:       resourceType,
			Name:       resourceName,
			Attributes: attributes,
			FilePath:   filePath,
			LineNumber: resourceBlock.DefRange.Start.Line,
		}

		tp.result.Resources = append(tp.result.Resources, resource)

		// Add warning if resource type is unknown
		if !IsKnownResource(resourceType) {
			tp.result.Errors = append(tp.result.Errors, ParseError{
				FilePath:  filePath,
				Line:      resourceBlock.DefRange.Start.Line,
				Message:   fmt.Sprintf("unknown resource type: %s", resourceType),
				ErrorType: "unknown_resource_type",
			})
		}
	}

	// Process variable blocks
	for _, varBlock := range content.Blocks {
		if varBlock.Type != "variable" && len(varBlock.Labels) == 0 {
			continue
		}
		if len(varBlock.Labels) > 0 {
			varName := varBlock.Labels[0]
			tp.result.Variables[varName] = &Block{
				Type:   varBlock.Type,
				Labels: varBlock.Labels,
			}
		}
	}

	// Process module blocks
	for _, modBlock := range content.Blocks {
		if modBlock.Type != "module" {
			continue
		}
		if len(modBlock.Labels) > 0 {
			modName := modBlock.Labels[0]
			tp.result.Modules[modName] = &Block{
				Type:   modBlock.Type,
				Labels: modBlock.Labels,
			}
		}
	}
}

// extractAttributes extracts attributes from an HCL attribute map
func (tp *TerraformParser) extractAttributes(attrs hcl.Attributes) map[string]interface{} {
	result := make(map[string]interface{})

	for name, attr := range attrs {
		// Decode the value
		val, diags := attr.Expr.Value(nil) // nil context for simple evaluation

		if diags.HasErrors() {
			// If we can't evaluate, store the raw expression
			result[name] = attr.Expr.Range().String()
			continue
		}

		result[name] = ctyValueToInterface(val)
	}

	return result
}

// ctyValueToInterface converts a cty.Value to a native Go interface
func ctyValueToInterface(val cty.Value) interface{} {
	if val.IsNull() {
		return nil
	}

	if !val.IsKnown() {
		return "<unknown>"
	}

	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Number:
		// Try to return as int if possible, otherwise float64
		f := val.AsBigFloat()
		i, _ := f.Int64()
		return i
	case cty.Bool:
		return val.True()
	case cty.List(cty.String), cty.List(cty.Number), cty.List(cty.Bool):
		// Handle list types
		return convertCtyListToSlice(val)
	case cty.Map(cty.String), cty.Map(cty.Number), cty.Map(cty.Bool):
		// Handle map types
		return convertCtyMapToMap(val)
	case cty.Object(nil):
		// Handle object types
		return convertCtyObjectToMap(val)
	default:
		return val.GoString()
	}
}

// convertCtyListToSlice converts a cty list to a Go slice
func convertCtyListToSlice(val cty.Value) []interface{} {
	var result []interface{}
	for it := val.ElementIterator(); it.Next(); {
		_, v := it.Element()
		result = append(result, ctyValueToInterface(v))
	}
	return result
}

// convertCtyMapToMap converts a cty map to a Go map
func convertCtyMapToMap(val cty.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for it := val.ElementIterator(); it.Next(); {
		k, v := it.Element()
		result[k.AsString()] = ctyValueToInterface(v)
	}
	return result
}

// convertCtyObjectToMap converts a cty object to a Go map
func convertCtyObjectToMap(val cty.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for it := val.ElementIterator(); it.Next(); {
		k, v := it.Element()
		result[k.AsString()] = ctyValueToInterface(v)
	}
	return result
}

// findTerraformFiles finds all .tf and .tf.json files in a directory
func (tp *TerraformParser) findTerraformFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories like .terraform, .git, etc.
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file has .tf or .tf.json extension
		if strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".tf.json") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetResources returns all discovered resources
func (tp *TerraformParser) GetResources() []Resource {
	if tp.result == nil {
		return []Resource{}
	}
	return tp.result.Resources
}

// GetResult returns the complete parse result
func (tp *TerraformParser) GetResult() *ParseResult {
	return tp.result
}
