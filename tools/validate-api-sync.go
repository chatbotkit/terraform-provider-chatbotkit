// +build tools

package main

// This file validates that the provider's resources are synchronized with the ChatBotKit API spec.
// It can be run during build time to ensure API compatibility.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/chatbotkit/terraform-provider/internal/client"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	APISpecURL = "https://api.chatbotkit.com/v1/spec"
)

// FieldDefinition represents a field in the API schema
type FieldDefinition struct {
	Name     string
	Type     string
	Required bool
	ReadOnly bool
}

// ResourceSchema represents the expected schema for a resource from the API
type ResourceSchema struct {
	Name   string
	Fields []FieldDefinition
}



// fetchAPISchema attempts to fetch the API schema from the endpoint
func fetchAPISchema() (map[string]ResourceSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch the spec as JSON first to work around missing $ref issues
	req, err := http.NewRequestWithContext(ctx, "GET", APISpecURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API spec: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API spec endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse with a loader that's more lenient
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	
	doc, err := loader.LoadFromData(body)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	fmt.Printf("Successfully loaded OpenAPI spec\n")

	// Parse the OpenAPI spec and extract resource schemas
	return parseOpenAPISpec(doc)
}

// isIntegrationType checks if a schema name represents an integration type
func isIntegrationType(name string) bool {
	lowerName := strings.ToLower(name)
	// Check for integration suffix or exact match
	return strings.HasSuffix(lowerName, "integration") || lowerName == "integration"
}

// parseOpenAPISpec parses an OpenAPI specification and extracts resource schemas
func parseOpenAPISpec(doc *openapi3.T) (map[string]ResourceSchema, error) {
	schemas := make(map[string]ResourceSchema)
	
	if doc.Components == nil || doc.Components.Schemas == nil {
		return schemas, fmt.Errorf("no schemas found in OpenAPI spec components")
	}
	
	fmt.Printf("Found %d schemas in OpenAPI spec\n", len(doc.Components.Schemas))
	
	// Parse each schema for our resources
	// Note: The API may have multiple integration types as separate schemas
	resourceNames := []string{"Bot", "Dataset", "Skillset", "File", "Secret"}
	
	// Add all integration types found in the API spec
	for schemaName := range doc.Components.Schemas {
		// Check if this is an integration type (e.g., SlackIntegration, DiscordIntegration, etc.)
		if isIntegrationType(schemaName) {
			resourceNames = append(resourceNames, schemaName)
		}
	}
	
	fmt.Printf("Looking for %d resource schemas: %v\n", len(resourceNames), resourceNames)
	
	for _, resourceName := range resourceNames {
		schemaRef, ok := doc.Components.Schemas[resourceName]
		if !ok {
			fmt.Printf("  Schema '%s' not found in spec\n", resourceName)
			continue
		}
		
		schema := schemaRef.Value
		if schema == nil {
			fmt.Printf("  Schema '%s' has nil value\n", resourceName)
			continue
		}
		
		// Extract properties
		if schema.Properties == nil {
			fmt.Printf("  Schema '%s' has no properties\n", resourceName)
			continue
		}
		
		// Get required fields set
		requiredMap := make(map[string]bool)
		for _, r := range schema.Required {
			requiredMap[r] = true
		}
		
		// Build field definitions
		var fields []FieldDefinition
		for propName, propRef := range schema.Properties {
			if propRef.Value == nil {
				continue
			}
			
			prop := propRef.Value
			fieldType := mapOpenAPITypeToGo(prop)
			
			fields = append(fields, FieldDefinition{
				Name:     propName,
				Type:     fieldType,
				Required: requiredMap[propName],
				ReadOnly: prop.ReadOnly,
			})
		}
		
		schemas[strings.ToLower(resourceName)] = ResourceSchema{
			Name:   strings.ToLower(resourceName),
			Fields: fields,
		}
		
		fmt.Printf("  ✓ Parsed schema '%s' with %d fields\n", resourceName, len(fields))
	}
	
	return schemas, nil
}

// mapOpenAPITypeToGo maps OpenAPI types to Go types
func mapOpenAPITypeToGo(schema *openapi3.Schema) string {
	// Type is *openapi3.Types which is a slice of strings
	if schema.Type == nil || len(*schema.Type) == 0 {
		return "interface{}"
	}
	
	// Get the first type (OpenAPI 3.1 allows multiple types)
	typeStr := (*schema.Type)[0]
	
	switch typeStr {
	case "number":
		if schema.Format == "float" || schema.Format == "double" {
			return "float64"
		}
		return "float64"
	case "integer":
		if schema.Format == "int32" {
			return "int32"
		}
		return "int64"
	case "boolean":
		return "bool"
	case "object":
		return "map[string]interface{}"
	case "array":
		return "[]interface{}"
	case "string":
		return "string"
	default:
		return "interface{}"
	}
}

// getClientModel returns the reflect.Type for a client model by resource name
func getClientModel(resourceName string) (reflect.Type, error) {
	switch resourceName {
	case "bot":
		return reflect.TypeOf(client.Bot{}), nil
	case "dataset":
		return reflect.TypeOf(client.Dataset{}), nil
	case "skillset":
		return reflect.TypeOf(client.Skillset{}), nil
	case "file":
		return reflect.TypeOf(client.File{}), nil
	case "secret":
		return reflect.TypeOf(client.Secret{}), nil
	default:
		// Handle integration types - all integration types map to the same Integration struct
		if isIntegrationType(resourceName) {
			return reflect.TypeOf(client.Integration{}), nil
		}
		return nil, fmt.Errorf("unknown resource: %s", resourceName)
	}
}

// StructFieldInfo contains information about a struct field
type StructFieldInfo struct {
	Name     string
	Type     string
	JSONName string
	Omitempty bool
}

// extractFieldsFromStruct extracts field information from a Go struct using reflection
func extractFieldsFromStruct(t reflect.Type) map[string]StructFieldInfo {
	fields := make(map[string]StructFieldInfo)
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		
		// Parse the JSON tag to get the field name and options
		parts := strings.Split(jsonTag, ",")
		jsonName := parts[0]
		if jsonName == "" {
			jsonName = field.Name
		}
		
		omitempty := false
		for _, part := range parts[1:] {
			if part == "omitempty" {
				omitempty = true
				break
			}
		}
		
		// Get the Go type
		fieldType := field.Type.String()
		
		fields[jsonName] = StructFieldInfo{
			Name:      field.Name,
			Type:      fieldType,
			JSONName:  jsonName,
			Omitempty: omitempty,
		}
	}
	
	return fields
}

// normalizeTypeName normalizes type names for comparison
func normalizeTypeName(typeName string) string {
	// Handle common type variations
	switch typeName {
	case "float64", "float32", "float":
		return "float64"
	case "int", "int64", "int32":
		return "int64"
	case "bool", "boolean":
		return "bool"
	case "string":
		return "string"
	}
	
	// Handle map types - only normalize generic interface maps
	if typeName == "map[string]interface{}" || strings.HasPrefix(typeName, "map[string]interface") {
		return "map[string]interface{}"
	}
	
	// Return other map types unchanged to catch specific map type mismatches
	return typeName
}

// validateResource validates a single resource against the API schema
func validateResource(resourceName string, schema ResourceSchema) (bool, []string, []string) {
	errors := []string{}
	warnings := []string{}
	
	// Get the client model for this resource
	modelType, err := getClientModel(resourceName)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to get client model: %v", err))
		return false, errors, warnings
	}
	
	// Extract fields from the struct
	structFields := extractFieldsFromStruct(modelType)
	
	// Create a map of expected fields from schema
	expectedFields := make(map[string]FieldDefinition)
	for _, field := range schema.Fields {
		expectedFields[field.Name] = field
	}
	
	// Check if all schema fields are present in the struct
	for fieldName, fieldDef := range expectedFields {
		structField, exists := structFields[fieldName]
		if !exists {
			if fieldDef.Required {
				errors = append(errors, fmt.Sprintf("  ✗ Missing required field '%s' (expected type: %s)", 
					fieldName, fieldDef.Type))
			} else {
				warnings = append(warnings, fmt.Sprintf("  ⚠ Missing optional field '%s' (expected type: %s)", 
					fieldName, fieldDef.Type))
			}
		} else {
			// Field exists, check if type matches
			expectedType := normalizeTypeName(fieldDef.Type)
			actualType := normalizeTypeName(structField.Type)
			
			if expectedType != actualType {
				errors = append(errors, fmt.Sprintf("  ✗ Field '%s' type mismatch: expected %s, got %s", 
					fieldName, expectedType, actualType))
			}
			
			// Check if read-only fields have omitempty tag
			if fieldDef.ReadOnly && !structField.Omitempty {
				warnings = append(warnings, fmt.Sprintf("  ⚠ Read-only field '%s' should have 'omitempty' tag", 
					fieldName))
			}
		}
	}
	
	// Check for extra fields in struct that aren't in schema
	for fieldName, structField := range structFields {
		if _, exists := expectedFields[fieldName]; !exists {
			warnings = append(warnings, fmt.Sprintf("  ⚠ Extra field '%s' (type: %s) not in API schema - may be provider-specific", 
				fieldName, structField.Type))
		}
	}
	
	return len(errors) == 0, errors, warnings
}

// saveSchemaToFile saves the schema to a JSON file for reference
func saveSchemaToFile(schemas map[string]ResourceSchema, filename string) error {
	data, err := json.MarshalIndent(schemas, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	
	return nil
}

// ValidateAPISync ensures that all expected resources are implemented and match the API schema
func ValidateAPISync() error {
	fmt.Println("API Sync Validation")
	fmt.Println("===================")
	fmt.Println()
	
	// Fetch the API schema from the endpoint
	fmt.Println("Fetching API schema from:", APISpecURL)
	apiSchemas, err := fetchAPISchema()
	if err != nil {
		return fmt.Errorf("failed to fetch API schema: %w", err)
	}
	
	if len(apiSchemas) == 0 {
		return fmt.Errorf("API returned empty schema")
	}
	
	fmt.Println("✓ Successfully fetched API schema from endpoint")
	
	// Save the schema to a file for reference
	schemaFile := fmt.Sprintf("%s/chatbotkit-api-schema-%d.json", os.TempDir(), time.Now().Unix())
	if err := saveSchemaToFile(apiSchemas, schemaFile); err != nil {
		fmt.Printf("⚠ Failed to save schema to file: %v\n", err)
	} else {
		fmt.Printf("✓ Saved API schema to: %s\n", schemaFile)
	}
	
	fmt.Println()
	
	// Validate each resource
	allValid := true
	validationResults := make(map[string]bool)
	totalWarnings := 0
	
	for resourceName, schema := range apiSchemas {
		fmt.Printf("Validating resource: %s\n", resourceName)
		isValid, errors, warnings := validateResource(resourceName, schema)
		validationResults[resourceName] = isValid
		totalWarnings += len(warnings)
		
		if isValid && len(warnings) == 0 {
			fmt.Printf("  ✓ Resource '%s' is correctly synchronized\n", resourceName)
		} else if isValid && len(warnings) > 0 {
			fmt.Printf("  ✓ Resource '%s' is synchronized (with warnings)\n", resourceName)
			for _, warnMsg := range warnings {
				fmt.Println(warnMsg)
			}
		} else {
			fmt.Printf("  ✗ Resource '%s' has validation errors:\n", resourceName)
			for _, errMsg := range errors {
				fmt.Println(errMsg)
			}
			if len(warnings) > 0 {
				fmt.Println("  Warnings:")
				for _, warnMsg := range warnings {
					fmt.Println(warnMsg)
				}
			}
			allValid = false
		}
		fmt.Println()
	}
	
	// Summary
	fmt.Println("Validation Summary")
	fmt.Println("==================")
	validCount := 0
	for _, isValid := range validationResults {
		if isValid {
			validCount++
		}
	}
	fmt.Printf("Resources validated: %d/%d\n", validCount, len(validationResults))
	if totalWarnings > 0 {
		fmt.Printf("Total warnings: %d\n", totalWarnings)
	}
	
	if !allValid {
		return fmt.Errorf("validation failed for one or more resources")
	}
	
	return nil
}

func main() {
	if err := ValidateAPISync(); err != nil {
		fmt.Fprintf(os.Stderr, "\n✗ API sync validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n✓ API sync validation passed")
}
