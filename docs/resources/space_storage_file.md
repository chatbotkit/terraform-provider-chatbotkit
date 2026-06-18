---
page_title: "chatbotkit_space_storage_file Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a file within a ChatBotKit Space's storage.
---

# chatbotkit_space_storage_file (Resource)

Manages a file stored at a path within a ChatBotKit Space's storage. Each space has an isolated, path-addressed storage area; this resource manages a single file within it.

Provide exactly one of `content`, `source`, or `source_url`. Content up to ~4.5MB is uploaded directly through the API; larger files are transparently uploaded to storage using a presigned request.

## Example Usage

### Inline Content

```terraform
resource "chatbotkit_space" "workspace" {
  name = "Team Workspace"
}

resource "chatbotkit_space_storage_file" "readme" {
  space_id     = chatbotkit_space.workspace.id
  path         = "docs/README.md"
  content      = "# Team Workspace\n\nShared files live here."
  content_type = "text/markdown"
}
```

### From a Local File

```terraform
resource "chatbotkit_space_storage_file" "logo" {
  space_id    = chatbotkit_space.workspace.id
  path        = "assets/logo.png"
  source      = "${path.module}/assets/logo.png"
  source_hash = filesha256("${path.module}/assets/logo.png")
}
```

### From a URL (fetched server-side)

```terraform
resource "chatbotkit_space_storage_file" "import" {
  space_id   = chatbotkit_space.workspace.id
  path       = "imports/dataset.csv"
  source_url = "https://example.com/dataset.csv"
}
```

## Argument Reference

The following arguments are supported:

- `space_id` - (Required, Forces new resource) The ID of the space that owns the storage.
- `path` - (Required, Forces new resource) The storage path of the file (e.g. `docs/report.pdf`).
- `content` - (Optional) Inline content to upload as the file body. Conflicts with `source` and `source_url`.
- `source` - (Optional) Path to a local file whose contents are uploaded. Conflicts with `content` and `source_url`.
- `source_url` - (Optional) An HTTP(S) or data URL that the platform fetches and stores server-side. Conflicts with `content` and `source`.
- `source_hash` - (Optional) A hash of the source content (e.g. `filesha256(...)`). Provide this alongside `source` so Terraform detects changes to a local file and re-uploads it.
- `content_type` - (Optional) The MIME type of the content. When omitted it is detected from the path extension (falling back to content sniffing). Ignored for `source_url`, where the platform determines the type while fetching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The identifier of the resource (`<space_id>/<path>`).
- `content_sha256` - The SHA-256 digest of the uploaded content. For `source_url`, this is the digest of the URL string.

~> **Note** File content is write-mostly. The provider does not read stored bytes back during refresh; changes are driven by configuration (and `source_hash`). Destroying this resource removes the file from space storage.

## Import

Storage files can be imported using a `<space_id>/<path>` identifier:

```bash
terraform import chatbotkit_space_storage_file.readme space_abc123/docs/README.md
```

After import you must set one of `content`, `source`, or `source_url` in configuration; on the next apply the content is re-uploaded.
