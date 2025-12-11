# Terraform Provider Tools

This directory contains tools for maintaining and validating the ChatBotKit Terraform Provider.

## validate-api-sync.go

A validation tool that ensures the provider's resource definitions are synchronized with the ChatBotKit API specification.

### Features

- **Automatic Schema Validation**: Compares provider resource fields against the API schema
- **Field Type Checking**: Validates that field types match between the provider and API
- **Missing Field Detection**: Identifies required and optional fields that may be missing
- **Dynamic Schema Fetching**: Fetches the latest schema directly from the ChatBotKit API endpoint
- **Multiple Integration Support**: Automatically detects and validates all integration types from the API spec

### Usage

Run the validation tool to check all resources:

```bash
go run tools/validate-api-sync.go
```

Or use the Makefile target:

```bash
make validate-api
```

The tool will fetch the current API schema and save it to a temporary file for reference.

### Output

The tool provides detailed output including:

- ✓ Successfully validated resources
- ✗ Validation errors (missing required fields, type mismatches)
- ⚠ Warnings (missing optional fields, extra fields, JSON tag issues)

Example output:

```
API Sync Validation
===================

Fetching API schema from: https://api.chatbotkit.com/v1/spec
✓ Successfully fetched API schema from endpoint
✓ Saved API schema to: /tmp/chatbotkit-api-schema-1234567890.json

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

### Resources Validated

The tool automatically validates all resources defined in the ChatBotKit API specification:

- **bot**: Bot resource with AI model configuration
- **dataset**: Dataset resource for knowledge bases
- **skillset**: Skillset resource for custom abilities
- **file**: File resource for document management
- **integration types**: All integration resources (Slack, Discord, WhatsApp, etc.)
- **secret**: Secret resource for sensitive data management

When the ChatBotKit API changes, simply run the validation tool to ensure all resources remain synchronized.

### Exit Codes

- `0`: Validation passed successfully
- `1`: Validation failed (one or more resources have errors)

### Dependencies

The tool uses:

- Standard Go libraries (reflect, encoding/json, net/http)
- ChatBotKit Terraform Provider client models (`github.com/chatbotkit/terraform-provider/internal/client`)

### Troubleshooting

**Issue**: "Failed to fetch API schema"

**Solution**: Ensure you have internet connectivity and the ChatBotKit API endpoint is accessible. The validation requires a live connection to the API.

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
