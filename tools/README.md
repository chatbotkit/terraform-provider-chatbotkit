# Terraform Provider Tools

This directory contains tools for maintaining and validating the ChatBotKit Terraform Provider.

## validate-api-sync.go

A validation tool that ensures the provider's resource definitions are synchronized with the ChatBotKit API specification.

### Features

- **Automatic Schema Validation**: Compares provider resource fields against the API schema
- **Field Type Checking**: Validates that field types match between the provider and API
- **Missing Field Detection**: Identifies required and optional fields that may be missing
- **API Endpoint Support**: Can fetch the latest schema from the ChatBotKit API endpoint
- **Embedded Schema Fallback**: Uses an embedded schema when the API endpoint is unavailable
- **Schema Export**: Export the embedded schema to a JSON file for reference

### Usage

#### Run Validation

Run the validation tool to check all resources:

```bash
go run tools/validate-api-sync.go
```

Or use the Makefile target:

```bash
make validate-api
```

#### Export Schema

Export the embedded API schema to a JSON file:

```bash
go run tools/validate-api-sync.go -export-schema schema.json
```

### Output

The tool provides detailed output including:

- ✓ Successfully validated resources
- ✗ Validation errors (missing required fields, type mismatches)
- ⚠ Warnings (missing optional fields, extra fields, JSON tag issues)

Example output:

```
API Sync Validation
===================

Attempting to fetch API schema from: https://api.chatbotkit.com/v1/spec
✓ Successfully fetched API schema from endpoint

Validating resource: bot
  ✓ Resource 'bot' is correctly synchronized

Validating resource: dataset
  ✓ Resource 'dataset' is correctly synchronized

...

Validation Summary
==================
Resources validated: 6/6

✓ API sync validation passed
```

### Validation Rules

The tool checks for:

1. **Required Fields**: All required fields from the API schema must be present in the client model
2. **Field Types**: Field types must match between the schema and client model (e.g., string, float64, int64, bool)
3. **Optional Fields**: Warns if optional fields are missing but doesn't fail validation
4. **Read-Only Fields**: Checks that read-only fields have the `omitempty` JSON tag
5. **Extra Fields**: Warns about fields in the client model that aren't in the API schema

### Schema Structure

The embedded schema defines each resource with its fields:

```go
type ResourceSchema struct {
    Name   string
    Fields []FieldDefinition
}

type FieldDefinition struct {
    Name     string  // Field name (JSON name)
    Type     string  // Expected Go type
    Required bool    // Whether the field is required
    ReadOnly bool    // Whether the field is read-only
}
```

### Integration with CI/CD

Include this validation in your build pipeline:

```bash
make check  # Runs fmt, lint, validate-api, and test
```

Or run it individually:

```bash
make validate-api
```

### Updating the Embedded Schema

When the ChatBotKit API changes:

1. Update the embedded schema in `getEmbeddedAPISchema()` function
2. Run validation to ensure all resources are synchronized
3. Update client models in `internal/client/` if needed
4. Update resource definitions in `internal/resources/` if needed

### Resources Validated

The tool validates the following resources:

- **bot**: Bot resource with AI model configuration
- **dataset**: Dataset resource for knowledge bases
- **skillset**: Skillset resource for custom abilities
- **file**: File resource for document management
- **integration**: Integration resource for connecting bots to platforms
- **secret**: Secret resource for sensitive data management

### Exit Codes

- `0`: Validation passed successfully
- `1`: Validation failed (one or more resources have errors)

### Dependencies

The tool uses:

- Standard Go libraries (reflect, encoding/json, net/http)
- ChatBotKit Terraform Provider client models (`github.com/chatbotkit/terraform-provider/internal/client`)

### Troubleshooting

**Issue**: "Could not fetch API schema from endpoint"

**Solution**: This is expected if the API endpoint is unavailable. The tool will automatically fall back to the embedded schema.

**Issue**: "Missing required field"

**Solution**: Add the field to the corresponding client model in `internal/client/` and the resource definition in `internal/resources/`.

**Issue**: "Field type mismatch"

**Solution**: Ensure the Go type in the client model matches the expected type in the schema. Common types: `string`, `int64`, `float64`, `bool`, `map[string]interface{}`.

### Future Enhancements

Possible improvements for the validation tool:

- Parse OpenAPI 3.0 specification format directly
- Validate data source schemas
- Check for deprecated fields
- Validate field constraints (min/max values, patterns)
- Integration tests with actual API calls
- Generate resource definitions from schema
