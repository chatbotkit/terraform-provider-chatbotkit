# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of ChatBotKit Terraform Provider
- Support for managing bots with full CRUD operations
- Support for managing datasets with full CRUD operations
- Support for managing skillsets with full CRUD operations
- Support for managing files with full CRUD operations
- Support for managing integrations with full CRUD operations
- Support for managing secrets with full CRUD operations
- Data sources for fetching individual resources
- Data sources for listing all resources
- API client implementation with Bearer token authentication
- Examples for all supported resources
- Comprehensive documentation
- API synchronization validation tool
- Makefile for common development tasks

### Features
- **Resources**: bot, dataset, skillset, file, integration, secret
- **Data Sources**: Single resource fetch and list operations for all resources
- **Authentication**: Bearer token via configuration or environment variable
- **API Version**: v1
- **Framework**: Built with Terraform Plugin Framework for modern provider development

### Excluded Resources
As per design specification, the following resources are intentionally excluded:
- Contacts
- Conversations
- Tasks
- Memory
- Spaces
- Ratings

## [0.1.0] - TBD

### Initial Release
- First public release of the ChatBotKit Terraform Provider
