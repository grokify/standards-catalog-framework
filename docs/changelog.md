# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

No unreleased changes.

## [0.2.0] - 2026-06-29

### Added

- `export` CLI command for JSON output
- `generate-web` CLI command for interactive HTML visualization
- Interactive D3.js visualization with graph and table views
- MkDocs documentation site with Material theme
- Filter controls for layer, organization, category, status
- Color-coded identity layers in visualization

### Changed

- README updated with web visualization documentation

## [0.1.0] - 2026-06-28

### Added

- Initial release
- Core catalog types: `Catalog`, `Standard`, `Implementation`
- Enumeration types: `StandardStatus`, `Category`, `IdentityLayer`, `ImplementationStatus`
- YAML and JSON file format support
- Validation system with configurable options
- CLI commands: `validate`, `list`, `show`, `query`, `stats`
- Go library with query methods:
  - `FindByID`
  - `FindByCategory`
  - `FindByLayer`
  - `FindByStatus`
  - `FindByTag`
  - `FindByOrganization`
  - `GetOrganizations`
  - `GetTags`
- Cross-reference validation for `compatibleWith`, `supersedes`, `supersededBy`
- Example catalog in `testdata/`

### Validation Features

- Required field validation
- ID format validation (lowercase, hyphens)
- URL validation
- Enum value validation
- ID uniqueness validation
- Cross-reference resolution

### CLI Features

- Tabular output for `list` and `query`
- JSON output for `show`
- Statistics by status, category, layer, organization
- Filter flags for `query` command

---

For detailed release notes, see the Releases section.

[Unreleased]: https://github.com/grokify/standards-catalog-framework/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/grokify/standards-catalog-framework/releases/tag/v0.2.0
[0.1.0]: https://github.com/grokify/standards-catalog-framework/releases/tag/v0.1.0
