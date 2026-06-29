# Library Usage

The Go library provides programmatic access to standards catalogs.

## Installation

```bash
go get github.com/grokify/standards-catalog-framework
```

## Loading Catalogs

### Load Only

```go
package main

import (
    "fmt"
    "log"

    "github.com/grokify/standards-catalog-framework/catalog"
)

func main() {
    c, err := catalog.Load("catalog.yaml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Loaded %d standards\n", len(c.Standards))
}
```

### Load and Validate

```go
c, errors, err := catalog.LoadAndValidate("catalog.yaml")
if err != nil {
    log.Fatal(err)
}

if errors.HasErrors() {
    for _, e := range errors {
        fmt.Printf("Validation error: %s\n", e)
    }
    os.Exit(1)
}
```

## Querying Standards

### Find by ID

```go
standard := c.FindByID("oauth2")
if standard != nil {
    fmt.Printf("Found: %s\n", standard.Name)
}
```

### Find by Layer

```go
agentStandards := c.FindByLayer(catalog.LayerAgent)
for _, s := range agentStandards {
    fmt.Printf("%s: %s\n", s.ID, s.Name)
}
```

### Find by Category

```go
authStandards := c.FindByCategory(catalog.CategoryAuthorization)
```

### Find by Status

```go
adoptedStandards := c.FindByStatus(catalog.StatusAdopted)
```

### Find by Organization

```go
ietfStandards := c.FindByOrganization("IETF")
```

### Find by Tag

```go
foundationalStandards := c.FindByTag("foundational")
```

## Aggregation

### Get All Organizations

```go
orgs := c.GetOrganizations()
fmt.Printf("Organizations: %v\n", orgs)
// Output: Organizations: [CNCF IETF OpenID Foundation]
```

### Get All Tags

```go
tags := c.GetTags()
fmt.Printf("Tags: %v\n", tags)
```

## Validation

### Basic Validation

```go
validator := catalog.NewValidator()
errors := validator.Validate(c)

if errors.HasErrors() {
    for _, e := range errors {
        fmt.Println(e)
    }
}
```

### Custom Validation Options

```go
validator := catalog.NewValidator(
    catalog.WithRequireImplementations(true),
    catalog.WithRequireDescription(true),
    catalog.WithAllowedOrganizations([]string{"IETF", "CNCF"}),
    catalog.WithAllowedCategories([]catalog.Category{
        catalog.CategoryAuthorization,
        catalog.CategoryAuthentication,
    }),
    catalog.WithAllowedLayers([]catalog.IdentityLayer{
        catalog.LayerAgent,
        catalog.LayerHuman,
    }),
)

errors := validator.Validate(c)
```

## Saving Catalogs

### Save as YAML

```go
err := catalog.Save(c, "output.yaml")
```

### Save as JSON

```go
err := catalog.Save(c, "output.json")
```

## Working with Types

### Standard Status

```go
const (
    StatusDraft      StandardStatus = "draft"
    StatusProposed   StandardStatus = "proposed"
    StatusAdopted    StandardStatus = "adopted"
    StatusDeprecated StandardStatus = "deprecated"
    StatusRetired    StandardStatus = "retired"
)
```

### Categories

```go
const (
    CategoryAuthentication Category = "authentication"
    CategoryAuthorization  Category = "authorization"
    CategoryIdentity       Category = "identity"
    CategoryProvisioning   Category = "provisioning"
    CategoryCommunication  Category = "communication"
    CategoryDiscovery      Category = "discovery"
    CategoryGovernance     Category = "governance"
)
```

### Identity Layers

```go
const (
    LayerHuman    IdentityLayer = "human"
    LayerAgent    IdentityLayer = "agent"
    LayerWorkload IdentityLayer = "workload"
    LayerService  IdentityLayer = "service"
)
```

### Implementation Status

```go
const (
    ImplStatusExperimental ImplementationStatus = "experimental"
    ImplStatusBeta         ImplementationStatus = "beta"
    ImplStatusStable       ImplementationStatus = "stable"
    ImplStatusDeprecated   ImplementationStatus = "deprecated"
)
```

## Complete Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/grokify/standards-catalog-framework/catalog"
)

func main() {
    // Load and validate
    c, errors, err := catalog.LoadAndValidate("catalog.yaml")
    if err != nil {
        log.Fatal(err)
    }
    if errors.HasErrors() {
        log.Fatal(errors)
    }

    // Print statistics
    fmt.Printf("Catalog: %s v%s\n", c.Metadata.Name, c.Metadata.Version)
    fmt.Printf("Total standards: %d\n\n", len(c.Standards))

    // Group by layer
    layers := []catalog.IdentityLayer{
        catalog.LayerHuman,
        catalog.LayerAgent,
        catalog.LayerWorkload,
        catalog.LayerService,
    }

    for _, layer := range layers {
        standards := c.FindByLayer(layer)
        fmt.Printf("%s layer (%d):\n", layer, len(standards))
        for _, s := range standards {
            fmt.Printf("  - %s (%s)\n", s.Name, s.Status)
        }
    }

    // Export specific standard as JSON
    if s := c.FindByID("oauth2"); s != nil {
        data, _ := json.MarshalIndent(s, "", "  ")
        fmt.Printf("\nOAuth 2.0 details:\n%s\n", data)
    }
}
```

## API Reference

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/grokify/standards-catalog-framework/catalog).
