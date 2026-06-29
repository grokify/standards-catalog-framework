# Getting Started

This guide walks you through creating your first standards catalog.

## Installation

### CLI Tool

```bash
go install github.com/grokify/standards-catalog-framework/cmd/standards-catalog@latest
```

### Go Library

```bash
go get github.com/grokify/standards-catalog-framework
```

## Creating a Catalog

### Step 1: Create the Catalog File

Create a file named `catalog.yaml`:

```yaml
apiVersion: standards-catalog/v1
kind: Catalog
metadata:
  name: my-standards-catalog
  version: "1.0.0"
  description: My collection of standards
  maintainers:
    - name: Your Name
      email: you@example.com

standards: []
```

### Step 2: Add Standards

Add standards to the `standards` array:

```yaml
standards:
  - id: oauth2
    name: OAuth 2.0
    version: "2.1"
    status: adopted
    organization: IETF
    specUrl: https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/
    category: authorization
    layer: human
    description: Industry-standard protocol for authorization
    protocols:
      - HTTP
      - TLS
    tags:
      - authorization
      - foundational

  - id: oidc
    name: OpenID Connect
    version: "1.0"
    status: adopted
    organization: OpenID Foundation
    specUrl: https://openid.net/specs/openid-connect-core-1_0.html
    category: authentication
    layer: human
    description: Identity layer on top of OAuth 2.0
    compatibleWith:
      - oauth2
    tags:
      - authentication
      - identity
```

### Step 3: Validate

```bash
standards-catalog validate catalog.yaml
```

Expected output:

```
Catalog 'my-standards-catalog' is valid (2 standards)
```

### Step 4: Explore

```bash
# List all standards
standards-catalog list catalog.yaml

# Show statistics
standards-catalog stats catalog.yaml

# Query by layer
standards-catalog query catalog.yaml --layer=human

# Show specific standard
standards-catalog show catalog.yaml oauth2
```

## Required Fields

Every standard must have these fields:

| Field | Description | Example |
|-------|-------------|---------|
| `id` | Unique identifier (lowercase, hyphens) | `oauth2`, `id-jag` |
| `name` | Human-readable name | `OAuth 2.0` |
| `version` | Standard version | `2.1`, `RFC 7519` |
| `status` | Maturity status | `draft`, `adopted` |
| `organization` | Issuing body | `IETF`, `CNCF` |
| `specUrl` | Link to specification | `https://...` |
| `category` | Functional category | `authorization` |
| `layer` | Identity layer | `human`, `agent` |

## Optional Fields

| Field | Description |
|-------|-------------|
| `description` | Brief description of the standard |
| `protocols` | Related protocols (HTTP, TLS, gRPC) |
| `implementations` | Known implementations |
| `compatibleWith` | IDs of compatible standards |
| `supersedes` | IDs of superseded standards |
| `supersededBy` | IDs of standards that supersede this |
| `tags` | Arbitrary labels for filtering |
| `metadata` | Custom key-value pairs |

## Adding Implementations

Document known implementations:

```yaml
standards:
  - id: spiffe
    name: SPIFFE
    # ... required fields ...
    implementations:
      - name: SPIRE
        url: https://github.com/spiffe/spire
        status: stable
      - name: go-spiffe
        url: https://github.com/spiffe/go-spiffe
        language: Go
        status: stable
```

## Documenting Relationships

### Compatibility

Standards that work together:

```yaml
- id: oidc
  compatibleWith:
    - oauth2
```

### Succession

When one standard replaces another:

```yaml
- id: oauth2-v3
  supersedes:
    - oauth2
```

## Generating Visualization

Create an interactive HTML visualization:

```bash
standards-catalog generate-web catalog.yaml -o standards-graph.html
```

Open `standards-graph.html` in a browser to explore your catalog visually.

## JSON Format

Catalogs can also be written in JSON:

```json
{
  "apiVersion": "standards-catalog/v1",
  "kind": "Catalog",
  "metadata": {
    "name": "my-standards-catalog",
    "version": "1.0.0"
  },
  "standards": [
    {
      "id": "oauth2",
      "name": "OAuth 2.0",
      "version": "2.1",
      "status": "adopted",
      "organization": "IETF",
      "specUrl": "https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/",
      "category": "authorization",
      "layer": "human"
    }
  ]
}
```

## Next Steps

- [CLI Reference](cli-reference.md) - Complete command documentation
- [Library Usage](library.md) - Use the Go library
- [Schema Reference](schema.md) - Full schema details
