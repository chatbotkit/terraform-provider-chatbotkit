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

# --- File content -----------------------------------------------------------

# A file record only; content is uploaded separately.
resource "chatbotkit_file" "notes" {
  name        = "notes.txt"
  description = "Notes managed by Terraform"
}

# Upload inline content to the file.
resource "chatbotkit_file_content" "notes" {
  file_id      = chatbotkit_file.notes.id
  content      = "Hello from Terraform"
  content_type = "text/plain"
}

# Upload a local file (use source_hash so changes are detected).
resource "chatbotkit_file" "report" {
  name = "report.csv"
}

resource "chatbotkit_file_content" "report" {
  file_id     = chatbotkit_file.report.id
  source      = "${path.module}/files/report.csv"
  source_hash = filesha256("${path.module}/files/report.csv")
}

# --- Space storage files ----------------------------------------------------

resource "chatbotkit_space" "workspace" {
  name        = "Team Workspace"
  description = "Workspace managed by Terraform"
}

# Inline content at a path within the space's storage.
resource "chatbotkit_space_storage_file" "readme" {
  space_id     = chatbotkit_space.workspace.id
  path         = "docs/README.md"
  content      = "# Team Workspace\n\nShared files live here."
  content_type = "text/markdown"
}

# Import a remote file; the platform fetches it server-side.
resource "chatbotkit_space_storage_file" "dataset" {
  space_id   = chatbotkit_space.workspace.id
  path       = "imports/dataset.csv"
  source_url = "https://example.com/dataset.csv"
}
