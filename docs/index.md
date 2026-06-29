# Standards Catalog Framework

A Go framework for creating, validating, and managing standards catalogs.

[![Go CI](https://github.com/grokify/standards-catalog-framework/actions/workflows/go-ci.yaml/badge.svg?branch=main)](https://github.com/grokify/standards-catalog-framework/actions/workflows/go-ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/standards-catalog-framework)](https://goreportcard.com/report/github.com/grokify/standards-catalog-framework)
[![GoDoc](https://pkg.go.dev/badge/github.com/grokify/standards-catalog-framework)](https://pkg.go.dev/github.com/grokify/standards-catalog-framework)

## Overview

Standards Catalog Framework provides:

- :material-file-document: **Schema definitions** for standards catalogs
- :material-check-circle: **Validation tools** for catalog entries
- :material-console: **CLI tool** for catalog management
- :material-book: **Go library** for programmatic access
- :material-graph: **Web visualization** for interactive exploration

## Quick Start

### Installation

```bash
go install github.com/grokify/standards-catalog-framework/cmd/standards-catalog@latest
```

### Create a Catalog

```yaml
# catalog.yaml
apiVersion: standards-catalog/v1
kind: Catalog
metadata:
  name: my-standards
  version: "1.0.0"

standards:
  - id: oauth2
    name: OAuth 2.0
    version: "2.1"
    status: adopted
    organization: IETF
    specUrl: https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/
    category: authorization
    layer: human
```

### Validate and List

```bash
# Validate catalog
standards-catalog validate catalog.yaml

# List all standards
standards-catalog list catalog.yaml

# Generate visualization
standards-catalog generate-web catalog.yaml -o standards-graph.html
```

## Ecosystem

```
Standards Catalog Framework (SCF)     ← You are here
        │
        ▼
Agent Standards Catalog (ASC)
        │
        ▼
agent-protocols (Go implementations)
```

| Project | Purpose |
|---------|---------|
| **Standards Catalog Framework** | Generic framework for any standards catalog |
| [Agent Standards Catalog](https://github.com/aistandardsio/agent-standards-catalog) | Catalog of AI agent standards |
| [agent-protocols](https://github.com/aistandardsio/agent-protocols) | Go implementations of protocols |

## Features

### Identity Layers

Standards are organized by identity layer:

| Layer | Color | Description | Examples |
|-------|-------|-------------|----------|
| **Human** | :material-circle:{ style="color: #00d4ff" } Cyan | Human identity and delegation | OIDC, SAML, ID-JAG |
| **Agent** | :material-circle:{ style="color: #ff6b00" } Orange | AI agent identity | AAuth, A2A, MCP |
| **Workload** | :material-circle:{ style="color: #a855f7" } Purple | Workload/service identity | SPIFFE, WIMSE |
| **Service** | :material-circle:{ style="color: #22c55e" } Green | Service-to-service auth | OAuth 2.0, mTLS |

### Categories

| Category | Description |
|----------|-------------|
| `authentication` | Identity verification |
| `authorization` | Access control decisions |
| `identity` | Identity representation |
| `provisioning` | Lifecycle management |
| `communication` | Protocol communication |
| `discovery` | Service/agent discovery |
| `governance` | Policy and compliance |

### Status Values

| Status | Description |
|--------|-------------|
| `draft` | In development |
| `proposed` | Proposed for adoption |
| `adopted` | Widely adopted |
| `deprecated` | Deprecated but still used |
| `retired` | No longer used |

## Next Steps

- [Getting Started](getting-started.md) - Detailed setup guide
- [CLI Reference](cli-reference.md) - All CLI commands
- [Library Usage](library.md) - Go library API
- [Visualization](visualization.md) - Interactive standards graph
