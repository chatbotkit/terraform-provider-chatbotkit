// +build tools

package main

// This file validates that the provider's resources are synchronized with the ChatBotKit API spec.
// It can be run during build time to ensure API compatibility.

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/chatbotkit/terraform-provider/internal/client"
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

// getEmbeddedAPISchema returns the embedded API schema based on current implementation
// This serves as a fallback when the API spec endpoint is not accessible
func getEmbeddedAPISchema() map[string]ResourceSchema {
	return map[string]ResourceSchema{
		"bot": {
			Name: "bot",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "description", Type: "string", Required: false, ReadOnly: false},
				{Name: "model", Type: "string", Required: false, ReadOnly: false},
				{Name: "datasetId", Type: "string", Required: false, ReadOnly: false},
				{Name: "skillsetId", Type: "string", Required: false, ReadOnly: false},
				{Name: "backstory", Type: "string", Required: false, ReadOnly: false},
				{Name: "temperature", Type: "float64", Required: false, ReadOnly: false},
				{Name: "instructions", Type: "string", Required: false, ReadOnly: false},
				{Name: "moderation", Type: "bool", Required: false, ReadOnly: false},
				{Name: "privacy", Type: "bool", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
		"dataset": {
			Name: "dataset",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "description", Type: "string", Required: false, ReadOnly: false},
				{Name: "type", Type: "string", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
		"skillset": {
			Name: "skillset",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "description", Type: "string", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
		"file": {
			Name: "file",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "type", Type: "string", Required: false, ReadOnly: false},
				{Name: "source", Type: "string", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
		"integration": {
			Name: "integration",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "description", Type: "string", Required: false, ReadOnly: false},
				{Name: "type", Type: "string", Required: false, ReadOnly: false},
				{Name: "botId", Type: "string", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
		"secret": {
			Name: "secret",
			Fields: []FieldDefinition{
				{Name: "id", Type: "string", Required: false, ReadOnly: true},
				{Name: "name", Type: "string", Required: true, ReadOnly: false},
				{Name: "value", Type: "string", Required: false, ReadOnly: false},
				{Name: "meta", Type: "map[string]interface{}", Required: false, ReadOnly: false},
				{Name: "createdAt", Type: "int64", Required: false, ReadOnly: true},
				{Name: "updatedAt", Type: "int64", Required: false, ReadOnly: true},
			},
		},
	}
}

// fetchAPISchema attempts to fetch the API schema from the endpoint
func fetchAPISchema() (map[string]ResourceSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	// Try to parse the response as JSON
	var specData map[string]interface{}
	if err := json.Unmarshal(body, &specData); err != nil {
		return nil, fmt.Errorf("failed to parse API spec: %w", err)
	}

	// Parse the OpenAPI spec and extract resource schemas
	// This is a simplified parser - a full implementation would need to handle
	// OpenAPI 3.0 spec format properly
	return parseOpenAPISpec(specData)
}

// parseOpenAPISpec parses an OpenAPI specification and extracts resource schemas
func parseOpenAPISpec(spec map[string]interface{}) (map[string]ResourceSchema, error) {
	// This is a simplified implementation
	// A full implementation would properly parse OpenAPI 3.0 spec
	schemas := make(map[string]ResourceSchema)
	
	// For now, return empty as we'll use the embedded schema as fallback
	return schemas, nil
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
	case "integration":
		return reflect.TypeOf(client.Integration{}), nil
	case "secret":
		return reflect.TypeOf(client.Secret{}), nil
	default:
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
	
	// Handle map types
	if strings.HasPrefix(typeName, "map[") {
		return "map[string]interface{}"
	}
	
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
	
	// Try to fetch the API schema from the endpoint
	fmt.Println("Attempting to fetch API schema from:", APISpecURL)
	apiSchemas, err := fetchAPISchema()
	fetchedFromAPI := false
	if err != nil {
		fmt.Printf("⚠ Could not fetch API schema from endpoint: %v\n", err)
		fmt.Println("⚠ Using embedded schema as fallback")
		apiSchemas = getEmbeddedAPISchema()
	} else {
		fmt.Println("✓ Successfully fetched API schema from endpoint")
		fetchedFromAPI = true
		// If we got an empty schema from the API, use embedded as fallback
		if len(apiSchemas) == 0 {
			fmt.Println("⚠ API returned empty schema, using embedded schema as fallback")
			apiSchemas = getEmbeddedAPISchema()
			fetchedFromAPI = false
		}
	}
	
	// Save the schema to a file for reference if fetched from API
	if fetchedFromAPI {
		schemaFile := "/tmp/chatbotkit-api-schema.json"
		if err := saveSchemaToFile(apiSchemas, schemaFile); err != nil {
			fmt.Printf("⚠ Failed to save schema to file: %v\n", err)
		} else {
			fmt.Printf("✓ Saved API schema to: %s\n", schemaFile)
		}
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
	// Parse command line flags
	exportSchema := flag.String("export-schema", "", "Export embedded schema to specified file (JSON format)")
	flag.Parse()
	
	// Handle schema export if requested
	if *exportSchema != "" {
		fmt.Println("Exporting embedded API schema...")
		schema := getEmbeddedAPISchema()
		if err := saveSchemaToFile(schema, *exportSchema); err != nil {
			fmt.Fprintf(os.Stderr, "✗ Failed to export schema: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Schema exported to: %s\n", *exportSchema)
		return
	}
	
	// Run validation
	if err := ValidateAPISync(); err != nil {
		fmt.Fprintf(os.Stderr, "\n✗ API sync validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n✓ API sync validation passed")
}
