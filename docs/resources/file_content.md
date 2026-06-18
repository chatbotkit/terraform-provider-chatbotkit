---
page_title: "chatbotkit_file_content Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Uploads content to an existing ChatBotKit File.
---

# chatbotkit_file_content (Resource)

Uploads content to an existing ChatBotKit File. A `chatbotkit_file` only creates the file record; use `chatbotkit_file_content` to store the actual bytes against it.

Provide exactly one of `content`, `source`, or `source_url`. Content up to ~4.5MB is uploaded directly through the API; larger files are transparently uploaded to storage using a presigned request.

## Example Usage

### Inline Content

```terraform
resource "chatbotkit_file" "notes" {
  name = "notes.txt"
}

resource "chatbotkit_file_content" "notes" {
  file_id      = chatbotkit_file.notes.id
  content      = "Hello from Terraform"
  content_type = "text/plain"
}
```

### From a Local File

```terraform
resource "chatbotkit_file" "report" {
  name = "report.pdf"
}

resource "chatbotkit_file_content" "report" {
  file_id     = chatbotkit_file.report.id
  source      = "${path.module}/files/report.pdf"
  source_hash = filesha256("${path.module}/files/report.pdf")
}
```

### From a URL (fetched server-side)

```terraform
resource "chatbotkit_file" "remote" {
  name = "data.csv"
}

resource "chatbotkit_file_content" "remote" {
  file_id    = chatbotkit_file.remote.id
  source_url = "https://example.com/data.csv"
}
```

## Argument Reference

The following arguments are supported:

- `file_id` - (Required, Forces new resource) The ID of the file to upload content to.
- `content` - (Optional) Inline content to upload as the file body. Conflicts with `source` and `source_url`.
- `source` - (Optional) Path to a local file whose contents are uploaded. Conflicts with `content` and `source_url`.
- `source_url` - (Optional) An HTTP(S) or data URL that the platform fetches and stores server-side. Conflicts with `content` and `source`.
- `source_hash` - (Optional) A hash of the source content (e.g. `filesha256(...)`). Provide this alongside `source` so Terraform detects changes to a local file and re-uploads it.
- `content_type` - (Optional) The MIME type of the content. When omitted it is detected from the file/extension (falling back to content sniffing). Ignored for `source_url`, where the platform determines the type while fetching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The identifier of the resource (equal to `file_id`).
- `content_sha256` - The SHA-256 digest of the uploaded content. For `source_url`, this is the digest of the URL string.

~> **Note** File content is write-mostly. The provider does not read stored bytes back during refresh; changes are driven by configuration (and `source_hash`). Destroying this resource does not remove the file's content — it is replaced on the next upload and removed when the underlying `chatbotkit_file` is deleted.

## Import

File content can be imported using the file ID:

```bash
terraform import chatbotkit_file_content.notes file_abc123def456
```

After import you must set one of `content`, `source`, or `source_url` in configuration; on the next apply the content is re-uploaded.
