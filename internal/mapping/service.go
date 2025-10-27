package mapping

import (
	"fmt"
	"sort"
	"strings"
)

// GetResourceActions retrieves all IAM actions needed for a resource
func (ms *MappingService) GetResourceActions(resourceType string, attributes map[string]interface{}) (*ResourceActions, error) {
	if !ms.db.HasMapping(resourceType) {
		return nil, fmt.Errorf("no mapping found for resource type: %s", resourceType)
	}

	// Check cache first
	cacheKey := ms.buildCacheKey(resourceType, attributes)
	ms.mu.RLock()
	if cached, exists := ms.cache[cacheKey]; exists {
		ms.mu.RUnlock()
		mapping, _ := ms.db.GetMapping(resourceType)
		return &ResourceActions{
			ResourceType: resourceType,
			Service:      mapping.Service,
			Actions:      cached,
			Reason:       "Retrieved from cache",
		}, nil
	}
	ms.mu.RUnlock()

	// Get the mapping
	mapping, _ := ms.db.GetMapping(resourceType)

	// Combine base actions with attribute-specific actions
	allActions := make(ActionSet)

	// Add all base actions
	for _, actionSet := range mapping.Actions {
		allActions.AddAll(actionSet)
	}

	// Add attribute-specific actions
	if attributes != nil && len(mapping.AttributeActions) > 0 {
		for attrName := range attributes {
			if attrActions, exists := mapping.AttributeActions[attrName]; exists {
				// Get the actions for this attribute (stored with "_default" key)
				for _, actionSet := range attrActions {
					allActions.AddAll(actionSet)
				}
			}
		}
	}

	// Cache the result
	ms.mu.Lock()
	ms.cache[cacheKey] = allActions
	ms.mu.Unlock()

	return &ResourceActions{
		ResourceType: resourceType,
		Service:      mapping.Service,
		Actions:      allActions,
		Reason:       fmt.Sprintf("Mapped from resource type %s with %d attributes", resourceType, len(attributes)),
	}, nil
}

// GetActionsForMultipleResources retrieves actions for multiple resources
func (ms *MappingService) GetActionsForMultipleResources(resources map[string]map[string]interface{}) (map[string]*ResourceActions, error) {
	result := make(map[string]*ResourceActions)

	for resourceType, attributes := range resources {
		actions, err := ms.GetResourceActions(resourceType, attributes)
		if err != nil {
			// Continue with other resources even if one fails
			continue
		}
		result[resourceType] = actions
	}

	return result, nil
}

// CombineActions combines multiple ResourceActions into a single ActionSet
func (ms *MappingService) CombineActions(resourceActions ...*ResourceActions) ActionSet {
	combined := make(ActionSet)

	for _, ra := range resourceActions {
		combined.AddAll(ra.Actions)
	}

	return combined
}

// GetActionsByService groups actions by their service
func (ms *MappingService) GetActionsByService(allActions ActionSet) map[string]ActionSet {
	result := make(map[string]ActionSet)

	for action := range allActions {
		// AWS actions are formatted as "service:Action"
		parts := strings.Split(action, ":")
		if len(parts) >= 1 {
			service := parts[0]
			if _, exists := result[service]; !exists {
				result[service] = make(ActionSet)
			}
			result[service][action] = true
		}
	}

	return result
}

// GetResourcesWithoutMapping returns resources that don't have mappings
func (ms *MappingService) GetResourcesWithoutMapping(resourceTypes []string) []string {
	var unmapped []string

	for _, resourceType := range resourceTypes {
		if !ms.db.HasMapping(resourceType) {
			unmapped = append(unmapped, resourceType)
		}
	}

	return unmapped
}

// GetCoverageStats returns coverage statistics for loaded mappings
func (ms *MappingService) GetCoverageStats() map[string]interface{} {
	allMappings := ms.db.GetAllMappings()

	totalActions := 0
	serviceCount := make(map[string]int)

	for _, mapping := range allMappings {
		for _, actionSet := range mapping.Actions {
			totalActions += len(actionSet)
		}
		if mapping.Service != "" {
			serviceCount[mapping.Service]++
		}
	}

	// Sort services for consistent output
	services := make([]string, 0, len(serviceCount))
	for service := range serviceCount {
		services = append(services, service)
	}
	sort.Strings(services)

	return map[string]interface{}{
		"total_mappings": len(allMappings),
		"total_actions":  totalActions,
		"services":       serviceCount,
		"service_list":   services,
	}
}

// ClearCache clears the lookup cache
func (ms *MappingService) ClearCache() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.cache = make(map[string]ActionSet)
}

// buildCacheKey builds a cache key from resource type and attributes
func (ms *MappingService) buildCacheKey(resourceType string, attributes map[string]interface{}) string {
	// Simple cache key: just use resource type if no attributes
	if len(attributes) == 0 {
		return resourceType
	}

	// Include attribute names in cache key
	var attrNames []string
	for attrName := range attributes {
		attrNames = append(attrNames, attrName)
	}
	sort.Strings(attrNames)

	return resourceType + ":" + strings.Join(attrNames, ",")
}

// GetMappingInfo returns information about a resource's mapping
func (ms *MappingService) GetMappingInfo(resourceType string) (map[string]interface{}, error) {
	mapping, exists := ms.db.GetMapping(resourceType)
	if !exists {
		return nil, fmt.Errorf("no mapping found for resource type: %s", resourceType)
	}

	// Count actions
	baseActionCount := 0
	for _, actionSet := range mapping.Actions {
		baseActionCount += len(actionSet)
	}

	attrActionCount := 0
	attrList := make([]string, 0)
	for attrName, attrActionsMap := range mapping.AttributeActions {
		attrList = append(attrList, attrName)
		for _, actionSet := range attrActionsMap {
			attrActionCount += len(actionSet)
		}
	}

	return map[string]interface{}{
		"resource_type":          resourceType,
		"service":                mapping.Service,
		"description":            mapping.Description,
		"base_action_count":      baseActionCount,
		"attribute_count":        len(mapping.AttributeActions),
		"attribute_names":        attrList,
		"attribute_action_count": attrActionCount,
		"total_possible_actions": baseActionCount + attrActionCount,
	}, nil
}
