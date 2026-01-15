# ChatBotKit Terraform Provider - Completion Plan

This document outlines the work required to complete the ChatBotKit Terraform Provider.

## Current State

The provider has a working code generator (`sites/main/scripts/gen-terraform-stubs.js`) that produces:

- ✅ 20 resource files with CRUD operations
- ✅ GraphQL client with HTTP transport and authentication
- ✅ Provider configuration with API key support
- ✅ Input/output types for all operations
- ✅ Go module configuration with `terraform-plugin-framework`

### Generated Resources

| Resource               | File                                |
| ---------------------- | ----------------------------------- |
| Blueprint              | `resource_blueprint.go`             |
| Bot                    | `resource_bot.go`                   |
| Dataset                | `resource_dataset.go`               |
| Discord Integration    | `resource_discord_integration.go`   |
| Email Integration      | `resource_email_integration.go`     |
| Extract Integration    | `resource_extract_integration.go`   |
| File                   | `resource_file.go`                  |
| MCP Server Integration | `resource_mcpserver_integration.go` |
| Messenger Integration  | `resource_messenger_integration.go` |
| Notion Integration     | `resource_notion_integration.go`    |
| Portal                 | `resource_portal.go`                |
| Secret                 | `resource_secret.go`                |
| Sitemap Integration    | `resource_sitemap_integration.go`   |
| Skillset               | `resource_skillset.go`              |
| Skillset Ability       | `resource_skillset_ability.go`      |
| Slack Integration      | `resource_slack_integration.go`     |
| Telegram Integration   | `resource_telegram_integration.go`  |
| Trigger Integration    | `resource_trigger_integration.go`   |
| Twilio Integration     | `resource_twilio_integration.go`    |
| WhatsApp Integration   | `resource_whats_app_integration.go` |

---

## ✅ Phase 1: Make It Runnable (COMPLETE)

**Goal**: Get a working Terraform provider that can be tested locally.  
**Effort**: ~30 minutes  
**Priority**: P0 - Required  
**Status**: ✅ Complete

### Task 1.1: Create `main.go`

Create the entry point for the Terraform provider binary.

**File**: `sdks/terraform/main.go`

```go
package main

import (
    "context"
    "flag"
    "log"

    "github.com/chatbotkit/terraform-sdk/internal/provider"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
    var debug bool
    flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
    flag.Parse()

    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/chatbotkit/chatbotkit",
        Debug:   debug,
    }

    err := providerserver.Serve(context.Background(), provider.New(version), opts)
    if err != nil {
        log.Fatal(err.Error())
    }
}
```

### Task 1.2: Verify Build

```bash
cd sdks/terraform
go build -o terraform-provider-chatbotkit
```

### Task 1.3: Create Local Dev Override

Create a Terraform CLI configuration for local development testing.

**File**: `~/.terraformrc` (or document in README)

```hcl
provider_installation {
  dev_overrides {
    "chatbotkit/chatbotkit" = "/path/to/sdks/terraform"
  }
  direct {}
}
```

### Task 1.4: Create Example Configuration

**File**: `sdks/terraform/examples/basic/main.tf`

```hcl
terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # api_key = "..." # Or set CHATBOTKIT_API_KEY env var
}

resource "chatbotkit_bot" "example" {
  name        = "My Test Bot"
  description = "Created via Terraform"
  backstory   = "You are a helpful assistant."
}
```

---

## ✅ Phase 2: Generator Improvements (COMPLETE)

**Goal**: Fix known issues in generated code.  
**Effort**: ~2 hours  
**Priority**: P1 - Important for production use  
**Status**: ✅ Complete

### Task 2.1: Fix Meta/Map Type Handling

**Problem**: All resources have `// Meta: TODO: convert map type` comments.

**Solution**: Update generator to handle `map[string]interface{}` conversion.

**File to modify**: `sites/main/scripts/gen-terraform-stubs.js`

Changes needed:

1. In `generateStringPtrHelper()`:

```javascript
} else if (field.goType === 'map[string]interface{}') {
  return `${field.goName}: convertMapToInterface(data.${field.goName})`
}
```

2. Add helper function to generated `client.go`:

```go
// convertMapToInterface converts types.Map to map[string]interface{}
func convertMapToInterface(ctx context.Context, m types.Map) map[string]interface{} {
    if m.IsNull() || m.IsUnknown() {
        return nil
    }
    result := make(map[string]interface{})
    m.ElementsAs(ctx, &result, false)
    return result
}
```

3. In `generateSetFromResponse()`:

```javascript
} else if (field.goType === 'map[string]interface{}') {
  return `if ${responseName}.${field.goName} != nil {
    mapValue, diags := types.MapValueFrom(ctx, types.StringType, ${responseName}.${field.goName})
    resp.Diagnostics.Append(diags...)
    data.${field.goName} = mapValue
  }`
}
```

### Task 2.2: Mark Sensitive Fields

**Problem**: API keys, tokens, and secrets are not marked as sensitive.

**Solution**: Update generator to detect and mark sensitive fields.

**File to modify**: `sites/main/scripts/gen-terraform-stubs.js`

Add to `generateResourceStub()` schema generation:

```javascript
const sensitivePatterns = [
  'token',
  'secret',
  'key',
  'password',
  'credential',
  'access_token',
]
const isSensitive = sensitivePatterns.some((p) =>
  field.name.toLowerCase().includes(p)
)

// In schema attribute generation:
if (isSensitive) {
  lines.push(`\t\t\t\tSensitive:           true,`)
}
```

### Task 2.3: Improve Read Operations

**Problem**: Current Read uses list+filter which is inefficient.

**Solution**: Use the GraphQL `node` query for direct ID lookup.

**Current approach** (inefficient):

```go
query Get${resource.name}($cursor: ID) {
    ${camelName}s(first: 1, after: $cursor) {
        edges { node { id ... } }
    }
}
// Then iterate to find matching ID
```

**Better approach**:

```go
query Get${resource.name}($id: ID!) {
    node(id: $id) {
        ... on ${resource.name} {
            id
            name
            description
            ...
        }
    }
}
```

Check if the GraphQL schema supports `node` query, otherwise use the specific fetch query for each resource type if available.

### Task 2.4: Handle 404/Not Found Gracefully

**Problem**: When a resource is deleted outside Terraform, Read fails with error.

**Solution**: Detect "not found" and remove from state.

Update Read methods:

```go
result, err := r.client.Get${resource.name}(ctx, data.ID.ValueString())
if err != nil {
    // Check if resource was deleted outside of Terraform
    if strings.Contains(err.Error(), "not found") {
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ${resourceNameLower}: %s", err))
    return
}
```

### Task 2.5: Add Computed Fields

**Problem**: Fields like `createdAt`, `updatedAt` are returned by API but not tracked.

**Solution**: Add computed-only fields to schema.

Update generator to add common computed fields:

```go
"created_at": schema.StringAttribute{
    MarkdownDescription: "Timestamp when the resource was created",
    Computed:            true,
},
"updated_at": schema.StringAttribute{
    MarkdownDescription: "Timestamp when the resource was last updated",
    Computed:            true,
},
```

---

## ✅ Phase 3: Data Sources (COMPLETE)

**Goal**: Allow reading existing resources without managing them.  
**Effort**: ~3 hours  
**Priority**: P2 - Nice to have  
**Status**: ✅ Complete

### Task 3.1: Add Data Source Generator

Create data sources for common resources that users might want to reference:

- `chatbotkit_bot` - Reference existing bots
- `chatbotkit_dataset` - Reference existing datasets
- `chatbotkit_skillset` - Reference existing skillsets
- `chatbotkit_blueprint` - Reference existing blueprints

**Example data source**:

```go
// datasource_bot.go
type BotDataSource struct {
    client *Client
}

type BotDataSourceModel struct {
    ID          types.String `tfsdk:"id"`
    Name        types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    // ... other fields
}

func (d *BotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BotDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    result, err := d.client.GetBot(ctx, data.ID.ValueString())
    // ... populate data from result
}
```

### Task 3.2: Register Data Sources in Provider

Update `provider.go`:

```go
func (p *ChatBotKitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewBotDataSource,
        NewDatasetDataSource,
        NewSkillsetDataSource,
        NewBlueprintDataSource,
    }
}
```

---

## Phase 4: Testing

**Goal**: Ensure provider works correctly and prevent regressions.  
**Effort**: ~4-6 hours  
**Priority**: P2 - Required for production release

### Task 4.1: Unit Tests for Client

**File**: `internal/provider/client_test.go`

Test the GraphQL client methods with mocked HTTP responses.

```go
func TestCreateBot(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request
        // Return mock response
    }))
    defer server.Close()

    client := NewClient("test-key", server.URL)
    result, err := client.CreateBot(context.Background(), CreateBotInput{
        Name: ptr("Test Bot"),
    })

    assert.NoError(t, err)
    assert.NotNil(t, result.ID)
}
```

### Task 4.2: Acceptance Tests

**File**: `internal/provider/resource_bot_test.go`

Terraform acceptance tests that create real resources.

```go
func TestAccBotResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read
            {
                Config: `
                    resource "chatbotkit_bot" "test" {
                        name = "test-bot"
                        description = "Test bot"
                    }
                `,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("chatbotkit_bot.test", "name", "test-bot"),
                    resource.TestCheckResourceAttrSet("chatbotkit_bot.test", "id"),
                ),
            },
            // Update
            {
                Config: `
                    resource "chatbotkit_bot" "test" {
                        name = "updated-bot"
                        description = "Updated description"
                    }
                `,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("chatbotkit_bot.test", "name", "updated-bot"),
                ),
            },
            // Import
            {
                ResourceName:      "chatbotkit_bot.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}
```

### Task 4.3: Test Configuration

**File**: `internal/provider/provider_test.go`

```go
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
    "chatbotkit": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
    if os.Getenv("CHATBOTKIT_API_KEY") == "" {
        t.Fatal("CHATBOTKIT_API_KEY must be set for acceptance tests")
    }
}
```

---

## Phase 5: Documentation

**Goal**: Prepare for Terraform Registry publication.  
**Effort**: ~2-3 hours  
**Priority**: P3 - Required for public release

### Task 5.1: Provider Documentation

**File**: `docs/index.md`

```markdown
# ChatBotKit Provider

The ChatBotKit provider allows you to manage ChatBotKit resources using Terraform.

## Example Usage

\`\`\`hcl
terraform {
required_providers {
chatbotkit = {
source = "chatbotkit/chatbotkit"
version = "~> 1.0"
}
}
}

provider "chatbotkit" {
api_key = var.chatbotkit_api_key
}
\`\`\`

## Authentication

The provider can be configured with an API key in two ways:

1. Set the `api_key` argument in the provider block
2. Set the `CHATBOTKIT_API_KEY` environment variable

## Schema

### Optional

- `api_key` (String, Sensitive) - The API key for ChatBotKit API
- `base_url` (String) - Custom API endpoint (defaults to https://api.chatbotkit.com/graphql)
```

### Task 5.2: Resource Documentation

Generate docs for each resource in `docs/resources/`.

**File**: `docs/resources/bot.md`

```markdown
# chatbotkit_bot

Manages a ChatBotKit Bot.

## Example Usage

\`\`\`hcl
resource "chatbotkit_bot" "assistant" {
name = "Customer Support Bot"
description = "Handles customer inquiries"
backstory = "You are a helpful customer support agent..."
model = "gpt-4"

dataset_id = chatbotkit_dataset.knowledge.id
skillset_id = chatbotkit_skillset.tools.id
}
\`\`\`

## Argument Reference

- `name` - (Optional) The name of the bot
- `description` - (Optional) Description of the bot
- `backstory` - (Optional) The system prompt/backstory for the bot
- `model` - (Optional) The AI model to use
- `dataset_id` - (Optional) ID of the dataset to attach
- `skillset_id` - (Optional) ID of the skillset to attach

## Attribute Reference

- `id` - The unique identifier of the bot
- `created_at` - Timestamp when the bot was created
- `updated_at` - Timestamp when the bot was last updated

## Import

Bots can be imported using their ID:

\`\`\`bash
terraform import chatbotkit_bot.assistant bot_abc123
\`\`\`
```

### Task 5.3: Generate Docs Automatically

Use `tfplugindocs` to generate documentation from schema:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
tfplugindocs generate
```

---

## Phase 6: Release Infrastructure

**Goal**: Enable automated releases to Terraform Registry.  
**Effort**: ~2 hours  
**Priority**: P3 - Required for public release

### Task 6.1: Goreleaser Configuration

**File**: `.goreleaser.yml`

```yaml
version: 2

builds:
  - env:
      - CGO_ENABLED=0
    mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - '-s -w -X main.version={{.Version}}'
    goos:
      - freebsd
      - windows
      - linux
      - darwin
    goarch:
      - amd64
      - '386'
      - arm
      - arm64
    ignore:
      - goos: darwin
        goarch: '386'
    binary: '{{ .ProjectName }}_v{{ .Version }}'

archives:
  - format: zip
    name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}'

checksum:
  name_template: '{{ .ProjectName }}_{{ .Version }}_SHA256SUMS'
  algorithm: sha256

signs:
  - artifacts: checksum
    args:
      - '--batch'
      - '--local-user'
      - '{{ .Env.GPG_FINGERPRINT }}'
      - '--output'
      - '${signature}'
      - '--detach-sign'
      - '${artifact}'

release:
  draft: true
```

### Task 6.2: GitHub Actions Workflow

**File**: `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'

      - name: Import GPG key
        uses: crazy-max/ghaction-import-gpg@v6
        id: import_gpg
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.GPG_PASSPHRASE }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
```

### Task 6.3: Terraform Registry Manifest

**File**: `terraform-registry-manifest.json`

```json
{
  "version": 1,
  "metadata": {
    "protocol_versions": ["6.0"]
  }
}
```

---

## Summary

| Phase                           | Effort    | Priority | Status      | Outcome                            |
| ------------------------------- | --------- | -------- | ----------- | ---------------------------------- |
| Phase 1: Make It Runnable       | 30 min    | P0       | ✅ Complete | Working provider for local testing |
| Phase 2: Generator Improvements | 2 hours   | P1       | ✅ Complete | Production-quality generated code  |
| Phase 3: Data Sources           | 3 hours   | P2       | ✅ Complete | Read-only resource access          |
| Phase 4: Testing                | 4-6 hours | P2       | ⬚ Pending   | Confidence in correctness          |
| Phase 5: Documentation          | 2-3 hours | P3       | ⬚ Pending   | Ready for public use               |
| Phase 6: Release Infrastructure | 2 hours   | P3       | ⬚ Pending   | Automated releases                 |

**Progress**: 3 of 6 phases complete (~5.5 hours done, ~8-11 hours remaining)

**Current state**: Provider is buildable and ready for local testing

---

## Quick Start

The provider is ready for local testing:

```bash
# 1. Build the provider
cd sdks/terraform
go build -o terraform-provider-chatbotkit

# 2. Set up dev override in ~/.terraformrc
cat >> ~/.terraformrc << 'EOF'
provider_installation {
  dev_overrides {
    "chatbotkit/chatbotkit" = "/path/to/cbk-platform/sdks/terraform"
  }
  direct {}
}
EOF

# 3. Test with example configuration
cd examples/basic
export CHATBOTKIT_API_KEY="your-api-key"
terraform plan
terraform apply
```

> **Note**: With dev_overrides, you don't need to run `terraform init`.
