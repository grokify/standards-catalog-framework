# CLI Reference

The `standards-catalog` CLI provides commands for managing standards catalogs.

## Installation

```bash
go install github.com/grokify/standards-catalog-framework/cmd/standards-catalog@latest
```

## Commands

### validate

Validate a catalog file for correctness.

```bash
standards-catalog validate <catalog-file>
```

**Arguments:**

- `<catalog-file>` - Path to YAML or JSON catalog file

**Example:**

```bash
$ standards-catalog validate catalog.yaml
Catalog 'AI Agent Standards' is valid (12 standards)
```

**Validation checks:**

- Required fields present
- ID format (lowercase, hyphens only)
- URL validity
- Enum values (status, category, layer)
- ID uniqueness
- Cross-reference validity (compatibleWith, supersedes)

---

### list

List all standards in a catalog.

```bash
standards-catalog list <catalog-file>
```

**Example:**

```bash
$ standards-catalog list catalog.yaml
ID        NAME                  VERSION  STATUS   CATEGORY        LAYER
--        ----                  -------  ------   --------        -----
oauth2    OAuth 2.0             2.1      adopted  authorization   human
oidc      OpenID Connect        1.0      adopted  authentication  human
aauth     AAuth Protocol        02       draft    authorization   agent
spiffe    SPIFFE                1.0      adopted  identity        workload
```

---

### show

Show details of a specific standard as JSON.

```bash
standards-catalog show <catalog-file> <standard-id>
```

**Arguments:**

- `<catalog-file>` - Path to catalog file
- `<standard-id>` - ID of the standard to show

**Example:**

```bash
$ standards-catalog show catalog.yaml oauth2
{
  "id": "oauth2",
  "name": "OAuth 2.0",
  "version": "2.1",
  "status": "adopted",
  "organization": "IETF",
  "specUrl": "https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/",
  "category": "authorization",
  "layer": "human",
  "description": "Industry-standard protocol for authorization",
  "protocols": ["HTTP", "TLS"],
  "tags": ["authorization", "foundational"]
}
```

---

### query

Query standards with filters.

```bash
standards-catalog query <catalog-file> [flags]
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--category` | Filter by category | `--category=authorization` |
| `--layer` | Filter by identity layer | `--layer=agent` |
| `--status` | Filter by status | `--status=adopted` |
| `--org` | Filter by organization | `--org=IETF` |
| `--tag` | Filter by tag | `--tag=foundational` |

**Examples:**

```bash
# Find all agent-layer standards
standards-catalog query catalog.yaml --layer=agent

# Find IETF authorization standards
standards-catalog query catalog.yaml --org=IETF --category=authorization

# Find adopted standards
standards-catalog query catalog.yaml --status=adopted
```

---

### stats

Show catalog statistics.

```bash
standards-catalog stats <catalog-file>
```

**Example:**

```bash
$ standards-catalog stats catalog.yaml
Catalog: AI Agent Standards (v1.0.0)
Total standards: 12

By Status:
  adopted: 5
  draft: 6
  proposed: 1

By Category:
  authorization: 4
  authentication: 3
  identity: 2
  communication: 2
  provisioning: 1

By Layer:
  human: 3
  agent: 5
  workload: 2
  service: 2

Organizations: IETF, OpenID Foundation, CNCF, Google, Anthropic
Tags: authorization, authentication, identity, agent, foundational
```

---

### export

Export catalog as formatted JSON.

```bash
standards-catalog export <catalog-file>
```

**Example:**

```bash
# Export to stdout
standards-catalog export catalog.yaml

# Export to file
standards-catalog export catalog.yaml > catalog.json
```

---

### generate-web

Generate an interactive HTML visualization.

```bash
standards-catalog generate-web <catalog-file> [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output HTML file path | `standards-graph.html` |

**Example:**

```bash
# Generate with default filename
standards-catalog generate-web catalog.yaml

# Generate with custom filename
standards-catalog generate-web catalog.yaml -o docs/visualization.html
```

**Features of generated visualization:**

- Force-directed graph view
- Table view
- Filter by layer, organization, category, status
- Color-coded by identity layer
- Click nodes to open specifications
- Zoom and pan support

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (validation failed, file not found, etc.) |

## File Formats

The CLI auto-detects format by file extension:

| Extension | Format |
|-----------|--------|
| `.yaml`, `.yml` | YAML |
| `.json` | JSON |
