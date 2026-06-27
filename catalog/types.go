// Package catalog provides types and utilities for managing standards catalogs.
package catalog

import (
	"time"
)

// Catalog represents a collection of standards.
type Catalog struct {
	// APIVersion identifies the catalog schema version.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" jsonschema:"required,enum=standards-catalog/v1"`

	// Kind identifies this as a Catalog resource.
	Kind string `json:"kind" yaml:"kind" jsonschema:"required,enum=Catalog"`

	// Metadata contains catalog metadata.
	Metadata CatalogMetadata `json:"metadata" yaml:"metadata" jsonschema:"required"`

	// Standards is the list of standards in this catalog.
	Standards []Standard `json:"standards" yaml:"standards" jsonschema:"required"`
}

// CatalogMetadata contains metadata about the catalog itself.
type CatalogMetadata struct {
	// Name is the catalog name.
	Name string `json:"name" yaml:"name" jsonschema:"required"`

	// Version is the catalog version (not the schema version).
	Version string `json:"version" yaml:"version" jsonschema:"required"`

	// Description describes the catalog's purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Maintainers lists the catalog maintainers.
	Maintainers []Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`

	// LastUpdated is when the catalog was last modified.
	LastUpdated time.Time `json:"lastUpdated,omitempty" yaml:"lastUpdated,omitempty"`
}

// Maintainer represents a catalog maintainer.
type Maintainer struct {
	Name  string `json:"name" yaml:"name" jsonschema:"required"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Standard represents a single standard in the catalog.
type Standard struct {
	// ID is the unique identifier for this standard.
	ID string `json:"id" yaml:"id" jsonschema:"required,pattern=^[a-z][a-z0-9-]*$"`

	// Name is the human-readable name.
	Name string `json:"name" yaml:"name" jsonschema:"required"`

	// Version is the standard's version.
	Version string `json:"version" yaml:"version" jsonschema:"required"`

	// Status indicates the standard's maturity.
	Status StandardStatus `json:"status" yaml:"status" jsonschema:"required"`

	// Organization is the issuing body.
	Organization string `json:"organization" yaml:"organization" jsonschema:"required"`

	// SpecURL is the link to the specification.
	SpecURL string `json:"specUrl" yaml:"specUrl" jsonschema:"required,format=uri"`

	// Category classifies the standard's purpose.
	Category Category `json:"category" yaml:"category" jsonschema:"required"`

	// Layer identifies the identity layer this standard operates at.
	Layer IdentityLayer `json:"layer" yaml:"layer" jsonschema:"required"`

	// Description provides a brief description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Protocols lists related protocols or standards.
	Protocols []string `json:"protocols,omitempty" yaml:"protocols,omitempty"`

	// Implementations lists known implementations.
	Implementations []Implementation `json:"implementations,omitempty" yaml:"implementations,omitempty"`

	// CompatibleWith lists compatible standards by ID.
	CompatibleWith []string `json:"compatibleWith,omitempty" yaml:"compatibleWith,omitempty"`

	// Supersedes lists standards this one supersedes by ID.
	Supersedes []string `json:"supersedes,omitempty" yaml:"supersedes,omitempty"`

	// SupersededBy lists standards that supersede this one by ID.
	SupersededBy []string `json:"supersededBy,omitempty" yaml:"supersededBy,omitempty"`

	// Tags are arbitrary labels for filtering.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Metadata contains additional key-value metadata.
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// StandardStatus represents the maturity status of a standard.
type StandardStatus string

// Standard status values.
const (
	StatusDraft      StandardStatus = "draft"
	StatusProposed   StandardStatus = "proposed"
	StatusAdopted    StandardStatus = "adopted"
	StatusDeprecated StandardStatus = "deprecated"
	StatusRetired    StandardStatus = "retired"
)

// Category classifies standards by their primary purpose.
type Category string

// Category values.
const (
	CategoryAuthentication Category = "authentication"
	CategoryAuthorization  Category = "authorization"
	CategoryIdentity       Category = "identity"
	CategoryProvisioning   Category = "provisioning"
	CategoryCommunication  Category = "communication"
	CategoryDiscovery      Category = "discovery"
	CategoryGovernance     Category = "governance"
)

// IdentityLayer identifies which identity layer a standard operates at.
type IdentityLayer string

// Identity layer values.
const (
	LayerHuman    IdentityLayer = "human"
	LayerAgent    IdentityLayer = "agent"
	LayerWorkload IdentityLayer = "workload"
	LayerService  IdentityLayer = "service"
)

// Implementation represents a known implementation of a standard.
type Implementation struct {
	// Name is the implementation name.
	Name string `json:"name" yaml:"name" jsonschema:"required"`

	// URL is the link to the implementation.
	URL string `json:"url" yaml:"url" jsonschema:"required,format=uri"`

	// Language is the programming language (if applicable).
	Language string `json:"language,omitempty" yaml:"language,omitempty"`

	// Status indicates the implementation's maturity.
	Status ImplementationStatus `json:"status,omitempty" yaml:"status,omitempty"`

	// Description provides additional details.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// ImplementationStatus represents the maturity of an implementation.
type ImplementationStatus string

// Implementation status values.
const (
	ImplStatusExperimental ImplementationStatus = "experimental"
	ImplStatusBeta         ImplementationStatus = "beta"
	ImplStatusStable       ImplementationStatus = "stable"
	ImplStatusDeprecated   ImplementationStatus = "deprecated"
)
