# Changelog

All notable changes to the ChatBotKit Terraform Provider are documented in this
file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.3.0]: https://github.com/chatbotkit/terraform-provider-chatbotkit/releases/tag/v1.3.0
