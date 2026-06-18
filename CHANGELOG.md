# Changelog

All notable changes to the ChatBotKit Terraform Provider are documented in this
file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.0] - 2026-06-18

### Added

- `chatbotkit_file_content` resource for uploading content to an existing
  `chatbotkit_file`. Accepts exactly one of `content` (inline), `source` (a local
  file), or `source_url` (an HTTP or data URL fetched server-side). Files up to
  ~4.5MB are uploaded directly; larger files are uploaded to storage via a
  presigned request automatically.
- `chatbotkit_space_storage_file` resource for managing a file at a path within a
  space's storage. Supports the same `content` / `source` / `source_url` inputs
  and transparent presigned uploads for large files, and removes the stored file
  on destroy.
- Both resources support an optional `content_type` (auto-detected from the file
  extension, falling back to content sniffing) and an optional `source_hash`
  (e.g. `filesha256(...)`) so changes to a local `source` file trigger
  re-uploads.

### Changed

- The provider registration now appends a hand-maintained set of resources
  (`manualResources()` in `resources_manual.go`) to the generated list, allowing
  resources backed by REST endpoints that are not expressed in the GraphQL schema
  (such as the new content uploads) to coexist with generated resources across
  regeneration.

## [1.3.0] - 2026-06-16

### Added

- `chatbotkit_googlechat_integration` resource for managing Google Chat integrations.
- `chatbotkit_instagram_integration` resource for managing Instagram integrations.
- `chatbotkit_microsoftteams_integration` resource for managing Microsoft Teams integrations.
- `chatbotkit_policy` resource for managing retention/usage policies.
- `chatbotkit_space` resource for managing spaces.
- `chatbotkit_support_integration` resource for managing support-email escalation.
- `chatbotkit_widget_integration` resource for managing the embeddable chat widget.
- New arguments surfaced on existing resources after regenerating from the current
  GraphQL schema:
  - `alias` on `chatbotkit_bot`, `chatbotkit_blueprint`, `chatbotkit_dataset`,
    `chatbotkit_file`, and `chatbotkit_skillset`.
  - `contact_collection` on `chatbotkit_messenger_integration`.
  - additional fields on `chatbotkit_extract_integration` and `chatbotkit_trigger_integration`.
  - matching read-only attributes on the `bot`, `dataset`, `skillset`, and `blueprint`
    data sources.

### Security

- The Google Chat `service_account_key` argument is now marked `Sensitive`, so the
  service-account credential is redacted from Terraform plan and state output.

### Changed

- Added `policy`, `space`, `widget`, and `support` create/update/delete operations to
  the ChatBotKit GraphQL API (they previously existed only on the REST API), then
  regenerated the provider (GraphQL client, resources, data sources, and provider
  registration). The provider now covers all 27 createable API resources, with no
  createable resource left unimplemented.

### Changed (breaking)

- Renamed the WhatsApp resource type from `chatbotkit_whats_app_integration` to
  `chatbotkit_whatsapp_integration` to match the documented name. Existing
  configurations and state referencing the old type must be updated (e.g. via
  `terraform state mv`).

[1.4.0]: https://github.com/chatbotkit/terraform-provider-chatbotkit/releases/tag/v1.4.0
[1.3.0]: https://github.com/chatbotkit/terraform-provider-chatbotkit/releases/tag/v1.3.0
