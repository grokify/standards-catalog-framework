# Schema Reference

Complete reference for the standards catalog schema.

## Catalog Structure

```yaml
apiVersion: standards-catalog/v1  # Required
kind: Catalog                      # Required
metadata:                          # Required
  name: string                     # Required
  version: string                  # Required
  description: string              # Optional
  maintainers:                     # Optional
    - name: string                 # Required
      email: string                # Optional
      url: string                  # Optional
  lastUpdated: timestamp           # Optional
standards:                         # Required
  - <Standard>
```

## Standard

```yaml
# Required fields
id: string              # Unique identifier
name: string            # Human-readable name
version: string         # Standard's version
status: StandardStatus  # Maturity status
organization: string    # Issuing body
specUrl: string         # URL to specification
category: Category      # Functional category
layer: IdentityLayer    # Identity layer

# Optional fields
description: string                 # Brief description
protocols: [string]                 # Related protocols
implementations: [Implementation]   # Known implementations
compatibleWith: [string]           # Compatible standard IDs
supersedes: [string]               # Superseded standard IDs
supersededBy: [string]             # Superseding standard IDs
tags: [string]                     # Arbitrary labels
metadata: {string: string}         # Custom key-value pairs
```

## Implementation

```yaml
name: string                    # Required - Implementation name
url: string                     # Required - URL to implementation
language: string                # Optional - Programming language
status: ImplementationStatus    # Optional - Maturity status
description: string             # Optional - Additional details
```

## Enumerations

### StandardStatus

| Value | Description |
|-------|-------------|
| `draft` | In development, not yet stable |
| `proposed` | Proposed for adoption, seeking feedback |
| `adopted` | Widely adopted and stable |
| `deprecated` | Deprecated, but still in use |
| `retired` | No longer used or maintained |

### Category

| Value | Description |
|-------|-------------|
| `authentication` | Identity verification protocols |
| `authorization` | Access control and permissions |
| `identity` | Identity representation and management |
| `provisioning` | Lifecycle management (SCIM, etc.) |
| `communication` | Protocol communication (A2A, MCP) |
| `discovery` | Service/agent discovery mechanisms |
| `governance` | Policy and compliance standards |

### IdentityLayer

| Value | Description | Examples |
|-------|-------------|----------|
| `human` | Human identity and delegation | OIDC, SAML, ID-JAG |
| `agent` | AI agent identity | AAuth, A2A, MCP |
| `workload` | Workload/service identity | SPIFFE, WIMSE |
| `service` | Service-to-service auth | OAuth 2.0, mTLS, JWT |

### ImplementationStatus

| Value | Description |
|-------|-------------|
| `experimental` | Early stage, API may change |
| `beta` | Testing stage, approaching stability |
| `stable` | Production-ready |
| `deprecated` | No longer maintained |

## Validation Rules

### ID Format

- Must start with a lowercase letter
- Can contain lowercase letters, numbers, and hyphens
- Pattern: `^[a-z][a-z0-9-]*$`

**Valid:** `oauth2`, `id-jag`, `a2a`, `http-signatures`

**Invalid:** `OAuth2`, `ID-JAG`, `2fa`, `-auth`

### URL Fields

- `specUrl` must be a valid URL
- Implementation `url` must be a valid URL

### Cross-References

- `compatibleWith` must reference existing standard IDs
- `supersedes` must reference existing standard IDs
- `supersededBy` must reference existing standard IDs

### Uniqueness

- Standard `id` must be unique within the catalog

## Complete Example

```yaml
apiVersion: standards-catalog/v1
kind: Catalog
metadata:
  name: AI Agent Standards Catalog
  version: "1.0.0"
  description: Catalog of identity and authorization standards for AI agents
  maintainers:
    - name: John Wang
      email: john@example.com
      url: https://example.com
  lastUpdated: 2024-01-15T10:30:00Z

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
    implementations:
      - name: golang-oauth2
        url: https://github.com/golang/oauth2
        language: Go
        status: stable
      - name: authlib
        url: https://github.com/lepture/authlib
        language: Python
        status: stable
    tags:
      - authorization
      - foundational
      - token
    metadata:
      rfc: "6749"
      maintainer: IETF OAuth WG

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
    protocols:
      - HTTP
      - TLS
    tags:
      - authentication
      - identity
      - sso

  - id: aauth
    name: AAuth Protocol
    version: "02"
    status: draft
    organization: IETF
    specUrl: https://datatracker.ietf.org/doc/draft-hardt-oauth-aauth-protocol/
    category: authorization
    layer: agent
    description: Agent authorization with mission-based consent
    compatibleWith:
      - id-jag
      - spiffe
    protocols:
      - HTTP
      - http-signatures
      - JWT
    implementations:
      - name: agent-protocols
        url: https://github.com/aistandardsio/agent-protocols
        language: Go
        status: beta
    tags:
      - agent
      - authorization
      - mission
```

## JSON Schema

The catalog format can be validated against a JSON Schema. See the [schema directory](https://github.com/grokify/standards-catalog-framework/tree/main/schema) for schema files.
