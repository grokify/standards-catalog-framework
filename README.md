# Standards Catalog Framework

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/standards-catalog-framework/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/standards-catalog-framework
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/standards-catalog-framework
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/standards-catalog-framework
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/standards-catalog-framework
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fstandards-catalog-framework
 [loc-svg]: https://tokei.rs/b1/github/grokify/standards-catalog-framework
 [repo-url]: https://github.com/grokify/standards-catalog-framework
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/standards-catalog-framework/blob/main/LICENSE

A generic Go framework for creating, validating, and managing standards catalogs.

## Overview

Standards Catalog Framework provides:

- 📋 **Schema definitions** for standards catalogs
- ✅ **Validation tools** for catalog entries
- 🖥️ **CLI tool** for catalog management
- 📚 **Go library** for programmatic access
- 🌐 **Web visualization** for interactive exploration

## Ecosystem Position

```
Standards Catalog Framework (SCF)     ← You are here
        │
        ▼
Agent Standards Catalog (ASC)
        │
        ▼
Open Agent Internet Architecture Framework (OAIAF)
        │
        ▼
agent-protocols
        │
        ▼
Generated protocol artifacts
```

| Project | Purpose |
|---------|---------|
| **Standards Catalog Framework** | Generic framework for any standards catalog |
| [Agent Standards Catalog](https://github.com/aistandardsio/agent-standards-catalog) | Catalog of AI agent standards built on SCF |
| [OAIAF](https://github.com/aistandardsio/oaiaf) | Reference architecture using cataloged standards |
| [agent-protocols](https://github.com/aistandardsio/agent-protocols) | Go implementations of cataloged protocols |

## Installation

```bash
go install github.com/grokify/standards-catalog-framework/cmd/standards-catalog@latest
```

## Usage

### CLI Commands

```bash
# Validate a catalog file
standards-catalog validate catalog.yaml

# List all standards
standards-catalog list catalog.yaml

# Show details of a specific standard
standards-catalog show catalog.yaml aauth

# Query with filters
standards-catalog query catalog.yaml --category=authorization --layer=agent

# Show statistics
standards-catalog stats catalog.yaml

# Export catalog as JSON
standards-catalog export catalog.yaml > catalog.json

# Generate interactive web visualization
standards-catalog generate-web catalog.yaml -o standards-graph.html
```

### Web Visualization

Generate interactive web visualizations that can be embedded in MkDocs or any HTML page:

```bash
# Generate HTML visualization with embedded data
standards-catalog generate-web catalog.yaml -o docs/standards-graph.html
```

Features:

- **Graph View**: Force-directed network showing standards and relationships
- **Table View**: Sortable table with all standard details
- **Filtering**: Filter by layer, organization, category, and status
- **Color-coded**: Standards colored by identity layer (Human, Agent, Workload, Service)
- **Interactive**: Click nodes to open specifications, drag to rearrange

To embed in MkDocs, copy the generated HTML to your `docs/` folder and link to it:

```markdown
[View Standards Graph](standards-graph.html)
```

Or use an iframe:

```markdown
<iframe src="standards-graph.html" width="100%" height="800" style="border:none;"></iframe>
```

### Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/grokify/standards-catalog-framework/catalog"
)

func main() {
    // Load and validate a catalog
    c, errors, err := catalog.LoadAndValidate("catalog.yaml")
    if err != nil {
        log.Fatal(err)
    }
    if errors.HasErrors() {
        log.Fatal(errors)
    }

    // Query standards
    agentStandards := c.FindByLayer(catalog.LayerAgent)
    for _, s := range agentStandards {
        fmt.Printf("%s: %s\n", s.ID, s.Name)
    }
}
```

## Catalog Format

Catalogs can be written in JSON or YAML:

```yaml
apiVersion: standards-catalog/v1
kind: Catalog
metadata:
  name: agent-standards
  version: "1.0.0"
  description: Catalog of AI agent-related standards

standards:
  - id: aauth
    name: AAuth Protocol
    version: "02"
    status: draft
    organization: IETF
    specUrl: https://datatracker.ietf.org/doc/draft-hardt-oauth-aauth-protocol/
    category: authorization
    layer: agent
    description: Agent authorization protocol with mission-based consent
    protocols:
      - oauth2
      - http-signatures
    tags:
      - agent
      - authorization
```

## Schema

### Standard Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (lowercase, hyphens allowed) |
| `name` | string | Yes | Human-readable name |
| `version` | string | Yes | Standard's version |
| `status` | enum | Yes | draft, proposed, adopted, deprecated, retired |
| `organization` | string | Yes | Issuing body |
| `specUrl` | URL | Yes | Link to specification |
| `category` | enum | Yes | authentication, authorization, identity, provisioning, communication, discovery, governance |
| `layer` | enum | Yes | human, agent, workload, service |
| `description` | string | No | Brief description |
| `protocols` | []string | No | Related protocols |
| `implementations` | []Implementation | No | Known implementations |
| `compatibleWith` | []string | No | Compatible standard IDs |
| `supersedes` | []string | No | Standards this supersedes |
| `tags` | []string | No | Arbitrary labels |

### Identity Layers

| Layer | Description | Example Standards |
|-------|-------------|-------------------|
| **human** | Human identity and delegation | OIDC, SAML, ID-JAG |
| **agent** | AI agent identity | AAuth, A2A |
| **workload** | Workload/service identity | SPIFFE, WIMSE |
| **service** | Service-to-service auth | OAuth 2.0, mTLS |

### Categories

| Category | Description | Example Standards |
|----------|-------------|-------------------|
| **authentication** | Identity verification | AAuth, OIDC, SAML |
| **authorization** | Access control decisions | AuthZEN, Cedar, OpenFGA |
| **identity** | Identity representation | SPIFFE, X.509 |
| **provisioning** | Lifecycle management | SCIM |
| **communication** | Protocol communication | A2A, MCP |
| **discovery** | Service/agent discovery | A2A Agent Cards |
| **governance** | Policy and compliance | - |

## Standards Organizations

The framework supports standards from various organizations:

| Organization | Abbreviation | Focus Area |
|--------------|--------------|------------|
| [Internet Engineering Task Force](https://ietf.org/) | IETF | Internet protocols, OAuth, HTTP |
| [OpenID Foundation](https://openid.net/) | OIDF | Identity protocols, AuthZEN |
| [Cloud Native Computing Foundation](https://cncf.io/) | CNCF | Cloud-native infrastructure, SPIFFE |
| [Linux Foundation](https://linuxfoundation.org/) | LF | Open source projects, A2A |
| [World Wide Web Consortium](https://w3.org/) | W3C | Web standards |
| [OASIS](https://oasis-open.org/) | OASIS | Enterprise standards, SAML |
| [ISO](https://iso.org/) | ISO | International standards |

## Example: AI Agent Standards

Standards commonly used in AI agent architectures:

### Identity & Authentication

| Standard | Organization | Status | Layer |
|----------|--------------|--------|-------|
| [AAuth](https://datatracker.ietf.org/doc/draft-hardt-oauth-aauth-protocol/) | IETF | Draft | Agent |
| [ID-JAG](https://datatracker.ietf.org/doc/draft-ietf-oauth-identity-assertion-authz-grant/) | IETF | Draft | Human |
| [SPIFFE](https://spiffe.io/) | CNCF | Adopted | Workload |
| [WIMSE](https://datatracker.ietf.org/wg/wimse/about/) | IETF | Draft | Workload |
| [SCIM Agent Resource](https://datatracker.ietf.org/doc/draft-wzdk-scim-agent-resource/) | IETF | Draft | Agent |

### Authorization

| Standard | Organization | Status | Description |
|----------|--------------|--------|-------------|
| [AuthZEN](https://openid.github.io/authzen/) | OIDF | Draft | PEP-PDP communication API |
| [Cedar](https://www.cedarpolicy.com/) | AWS | Adopted | ABAC policy language |
| [OpenFGA](https://openfga.dev/) | CNCF | Adopted | ReBAC authorization |

### Interoperability

| Standard | Organization | Status | Description |
|----------|--------------|--------|-------------|
| [A2A](https://google.github.io/A2A/) | LF | Draft | Agent-to-Agent protocol |
| [MCP](https://spec.modelcontextprotocol.io/) | Anthropic | Draft | Model Context Protocol |

### Foundational

| Standard | Organization | Status | Description |
|----------|--------------|--------|-------------|
| [OAuth 2.1](https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/) | IETF | Draft | Authorization framework |
| [RFC 9421](https://www.rfc-editor.org/rfc/rfc9421) | IETF | Adopted | HTTP Message Signatures |
| [RFC 8693](https://www.rfc-editor.org/rfc/rfc8693) | IETF | Adopted | Token Exchange |

## Validation

The framework validates:

- ✅ Required fields presence
- ✅ ID uniqueness and format
- ✅ URL validity
- ✅ Status/category/layer enum values
- ✅ Cross-references (compatibleWith, supersedes)

## License

MIT License
