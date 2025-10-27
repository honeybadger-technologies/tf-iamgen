package mapping

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadMappings loads all mapping YAML files from the mappings directory
func (db *MappingDatabase) LoadMappings(mappingsDir string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if directory exists
	if _, err := os.Stat(mappingsDir); os.IsNotExist(err) {
		return fmt.Errorf("mappings directory not found: %s", mappingsDir)
	}

	// Find all YAML files in the mappings directory
	files, err := filepath.Glob(filepath.Join(mappingsDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to list mapping files: %w", err)
	}

	// Also check for .yml files
	ymlFiles, err := filepath.Glob(filepath.Join(mappingsDir, "*.yml"))
	if err != nil {
		return fmt.Errorf("failed to list mapping files: %w", err)
	}
	files = append(files, ymlFiles...)

	if len(files) == 0 {
		return fmt.Errorf("no mapping files found in %s", mappingsDir)
	}

	// Load each file
	for _, file := range files {
		if err := db.loadMappingFile(file); err != nil {
			return fmt.Errorf("failed to load mapping file %s: %w", file, err)
		}
	}

	db.loaded = true
	return nil
}

// loadMappingFile loads a single YAML mapping file
func (db *MappingDatabase) loadMappingFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the file structure
	var fileData map[string]interface{}
	if err := yaml.Unmarshal(content, &fileData); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Each file contains resource mappings
	// Structure: resource_type -> mapping definition
	for resourceType, mappingData := range fileData {
		if mappingData == nil {
			continue
		}

		// Convert to map for easier handling
		mappingMap, ok := mappingData.(map[string]interface{})
		if !ok {
			continue
		}

		// Parse the mapping
		mapping, err := parseMappingData(mappingMap)
		if err != nil {
			return fmt.Errorf("failed to parse mapping for %s: %w", resourceType, err)
		}

		db.mappings[resourceType] = mapping
	}

	return nil
}

// parseMappingData converts raw YAML data to ResourceActionMap
func parseMappingData(data map[string]interface{}) (*ResourceActionMap, error) {
	mapping := &ResourceActionMap{
		Actions:          make(map[string]ActionSet),
		AttributeActions: make(map[string]map[string]ActionSet),
	}

	// Parse service field
	if service, ok := data["service"].(string); ok {
		mapping.Service = service
	}

	// Parse description field
	if desc, ok := data["description"].(string); ok {
		mapping.Description = desc
	}

	// Parse actions field
	if actionsData, ok := data["actions"]; ok {
		if actionsMap, ok := actionsData.(map[string]interface{}); ok {
			for actionName, actionValue := range actionsMap {
				actions, err := parseActionList(actionValue)
				if err != nil {
					return nil, fmt.Errorf("failed to parse action %s: %w", actionName, err)
				}
				mapping.Actions[actionName] = actions
			}
		}
	}

	// Parse attribute_actions field
	if attrActionsData, ok := data["attribute_actions"]; ok {
		if attrActionsMap, ok := attrActionsData.(map[string]interface{}); ok {
			for attrName, attrValue := range attrActionsMap {
				// attrValue should be a list of actions directly, not a nested map
				actions, err := parseActionList(attrValue)
				if err != nil {
					return nil, fmt.Errorf("failed to parse attribute action %s: %w", attrName, err)
				}

				// Store as a map with a dummy key since our structure expects map[string]map[string]ActionSet
				if mapping.AttributeActions[attrName] == nil {
					mapping.AttributeActions[attrName] = make(map[string]ActionSet)
				}
				mapping.AttributeActions[attrName]["_default"] = actions
			}
		}
	}

	return mapping, nil
}

// parseActionList converts various formats to ActionSet
func parseActionList(data interface{}) (ActionSet, error) {
	actionSet := make(ActionSet)

	switch v := data.(type) {
	case []interface{}:
		// List of strings
		for _, item := range v {
			if action, ok := item.(string); ok {
				actionSet[action] = true
			}
		}
	case string:
		// Single string
		actionSet[v] = true
	default:
		return nil, fmt.Errorf("unexpected action format: %T", v)
	}

	return actionSet, nil
}

// GetMapping retrieves a resource mapping from the database
func (db *MappingDatabase) GetMapping(resourceType string) (*ResourceActionMap, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	mapping, exists := db.mappings[resourceType]
	return mapping, exists
}

// HasMapping checks if a resource type has a mapping
func (db *MappingDatabase) HasMapping(resourceType string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, exists := db.mappings[resourceType]
	return exists
}

// GetAllMappings returns all mappings
func (db *MappingDatabase) GetAllMappings() map[string]*ResourceActionMap {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Return a copy to prevent external modifications
	result := make(map[string]*ResourceActionMap)
	for k, v := range db.mappings {
		result[k] = v
	}
	return result
}

// Clear clears all mappings from the database
func (db *MappingDatabase) Clear() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.mappings = make(map[string]*ResourceActionMap)
	db.loaded = false
}

// IsLoaded returns whether mappings have been loaded
func (db *MappingDatabase) IsLoaded() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.loaded
}
