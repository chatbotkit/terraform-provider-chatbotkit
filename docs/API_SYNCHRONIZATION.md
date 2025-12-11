# API Synchronization

This document describes how the ChatBotKit Terraform Provider stays synchronized with the ChatBotKit API.

## Overview

The provider is designed to maintain synchronization with the ChatBotKit API v1 specification. This ensures that any changes to the API are reflected in the provider implementation.

## API Specification

- **Base URL**: `https://api.chatbotkit.com/v1`
- **Spec URL**: `https://api.chatbotkit.com/v1/spec`
- **Documentation**: https://chatbotkit.com/docs/api

## Resource Mapping

The provider implements Terraform resources that map directly to ChatBotKit API endpoints:

| Terraform Resource | API Endpoints |
|-------------------|---------------|
| `chatbotkit_bot` | `/bot/list`, `/bot/create`, `/bot/{id}/fetch`, `/bot/{id}/update`, `/bot/{id}/delete` |
| `chatbotkit_dataset` | `/dataset/list`, `/dataset/create`, `/dataset/{id}/fetch`, `/dataset/{id}/update`, `/dataset/{id}/delete` |
| `chatbotkit_skillset` | `/skillset/list`, `/skillset/create`, `/skillset/{id}/fetch`, `/skillset/{id}/update`, `/skillset/{id}/delete` |
| `chatbotkit_file` | `/file/list`, `/file/upload`, `/file/{id}/fetch`, `/file/{id}/update`, `/file/{id}/delete` |
| `chatbotkit_integration` | `/integration/list`, `/integration/create`, `/integration/{id}/fetch`, `/integration/{id}/update`, `/integration/{id}/delete` |
| `chatbotkit_secret` | `/secret/list`, `/secret/create`, `/secret/{id}/fetch`, `/secret/{id}/update`, `/secret/{id}/delete` |

## Data Source Mapping

Each resource has corresponding data sources:

- **Single resource data sources**: Fetch a single resource by ID (e.g., `chatbotkit_bot`)
- **List data sources**: List all resources (e.g., `chatbotkit_bots`)

## Synchronization Process

### 1. Regular Review

Periodically review the ChatBotKit API specification for:
- New resources
- Changed resource schemas
- Deprecated endpoints
- New required or optional fields

### 2. Schema Alignment

Ensure that Terraform resource schemas match the API:

```go
// Example: Bot resource schema must match API structure
type Bot struct {
    ID              string  `json:"id,omitempty"`
    Name            string  `json:"name"`
    Description     string  `json:"description,omitempty"`
    Model           string  `json:"model,omitempty"`
    // ... other fields matching API
}
```

### 3. Validation Tool

The validation tool automatically checks API synchronization by:
- Comparing resource field definitions against the API schema
- Validating field types (string, int64, float64, bool, etc.)
- Detecting missing required and optional fields
- Checking JSON tag correctness
- Attempting to fetch the latest schema from the API endpoint
- Using an embedded schema as a fallback when the API is unavailable

Run the validation tool:

```bash
# Run validation
go run tools/validate-api-sync.go

# Or use the Makefile
make validate-api

# Export the embedded schema for reference
go run tools/validate-api-sync.go -export-schema schema.json
```

The validation tool provides detailed output:
- ✓ Successfully validated resources
- ✗ Validation errors (missing required fields, type mismatches)
- ⚠ Warnings (missing optional fields, extra fields)

### 4. Build-Time Checks

The provider includes build-time validation:

```bash
make check  # Runs fmt, lint, validate-api, and test
```

## Excluded Resources

As per the design specification, the following resource types are explicitly excluded:

- Contacts
- Conversations
- Tasks
- Memory
- Spaces
- Ratings

These are intentionally not implemented in the provider.

## Adding New Resources

When the ChatBotKit API adds new resources:

1. **Check the API spec** to understand the resource structure
2. **Create the client model** in `internal/client/`
3. **Implement the resource** in `internal/resources/`
4. **Add data sources** in `internal/datasources/`
5. **Register in provider** (`internal/provider/provider.go`)
6. **Add tests** for the new resource
7. **Update documentation** and examples
8. **Update embedded schema** in `tools/validate-api-sync.go` (function `getEmbeddedAPISchema()`)
9. **Run validation** to ensure the new resource is correctly synchronized

## Handling API Changes

### Breaking Changes

If the API introduces breaking changes:

1. Update the provider version (major version bump)
2. Update resource schemas to match new API structure
3. Add migration guides if state changes are required
4. Document breaking changes in CHANGELOG

### Non-Breaking Changes

For backward-compatible changes:

1. Add new optional fields to resource schemas
2. Update client models
3. Bump minor version
4. Document new features

### Deprecations

If the API deprecates features:

1. Mark resources/attributes as deprecated in schema
2. Add deprecation warnings
3. Document migration path
4. Plan removal for next major version

## Testing Synchronization

### Manual Testing

1. Compare provider schemas with API spec
2. Test CRUD operations for each resource
3. Verify data sources return expected data
4. Check error handling matches API responses

### Automated Testing

```bash
# Run all tests
make test

# Run specific resource tests
go test ./internal/resources/... -v

# Run data source tests
go test ./internal/datasources/... -v
```

## Continuous Monitoring

To maintain synchronization:

1. **Subscribe to API updates**: Monitor ChatBotKit API changelog
2. **Regular audits**: Review API spec quarterly
3. **Community feedback**: Track issues reporting API mismatches
4. **Automated checks**: Consider CI/CD integration to validate against API spec

## Resources

- [ChatBotKit API Documentation](https://chatbotkit.com/docs/api)
- [ChatBotKit API Spec](https://api.chatbotkit.com/v1/spec)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
