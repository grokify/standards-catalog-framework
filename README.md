# Standards Catalog Framework

A generic Go framework for creating, validating, and managing standards catalogs.

## Overview

Standards Catalog Framework provides:

- **Schema definitions** for standards catalogs
- **Validation tools** for catalog entries
- **CLI tool** for catalog management
- **Go library** for programmatic access

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

- **human**: Standards for human identity (OIDC, SAML, ID-JAG)
- **agent**: Standards for AI agent identity (AAuth)
- **workload**: Standards for workload/service identity (SPIFFE)
- **service**: Standards for service-to-service auth (OAuth 2.0)

## Validation

The framework validates:

- Required fields presence
- ID uniqueness and format
- URL validity
- Status/category/layer enum values
- Cross-references (compatibleWith, supersedes)

## License

MIT License
