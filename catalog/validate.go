package catalog

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with context.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// HasErrors returns true if there are any validation errors.
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Validator validates catalog structures.
type Validator struct {
	// RequireImplementations requires at least one implementation per standard.
	RequireImplementations bool

	// RequireDescription requires descriptions for standards.
	RequireDescription bool

	// AllowedOrganizations restricts which organizations are valid.
	AllowedOrganizations []string

	// AllowedCategories restricts which categories are valid.
	AllowedCategories []Category

	// AllowedLayers restricts which layers are valid.
	AllowedLayers []IdentityLayer
}

// NewValidator creates a new validator with default settings.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates a catalog and returns all errors found.
func (v *Validator) Validate(c *Catalog) ValidationErrors {
	var errors ValidationErrors

	// Validate catalog metadata
	errors = append(errors, v.validateCatalogMetadata(c)...)

	// Track IDs for uniqueness check
	seenIDs := make(map[string]bool)

	// Validate each standard
	for i, std := range c.Standards {
		prefix := fmt.Sprintf("standards[%d]", i)

		// Check ID uniqueness
		if seenIDs[std.ID] {
			errors = append(errors, ValidationError{
				Field:   prefix + ".id",
				Message: fmt.Sprintf("duplicate ID: %s", std.ID),
			})
		}
		seenIDs[std.ID] = true

		errors = append(errors, v.validateStandard(prefix, &std)...)
	}

	// Validate cross-references
	errors = append(errors, v.validateCrossReferences(c, seenIDs)...)

	return errors
}

func (v *Validator) validateCatalogMetadata(c *Catalog) ValidationErrors {
	var errors ValidationErrors

	if c.APIVersion != "standards-catalog/v1" {
		errors = append(errors, ValidationError{
			Field:   "apiVersion",
			Message: fmt.Sprintf("expected 'standards-catalog/v1', got '%s'", c.APIVersion),
		})
	}

	if c.Kind != "Catalog" {
		errors = append(errors, ValidationError{
			Field:   "kind",
			Message: fmt.Sprintf("expected 'Catalog', got '%s'", c.Kind),
		})
	}

	if c.Metadata.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "metadata.name",
			Message: "name is required",
		})
	}

	if c.Metadata.Version == "" {
		errors = append(errors, ValidationError{
			Field:   "metadata.version",
			Message: "version is required",
		})
	}

	return errors
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (v *Validator) validateStandard(prefix string, s *Standard) ValidationErrors {
	var errors ValidationErrors

	// Required fields
	if s.ID == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".id",
			Message: "id is required",
		})
	} else if !idPattern.MatchString(s.ID) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".id",
			Message: "id must start with lowercase letter and contain only lowercase letters, numbers, and hyphens",
		})
	}

	if s.Name == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".name",
			Message: "name is required",
		})
	}

	if s.Version == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".version",
			Message: "version is required",
		})
	}

	if s.Organization == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".organization",
			Message: "organization is required",
		})
	} else if len(v.AllowedOrganizations) > 0 {
		if !containsString(v.AllowedOrganizations, s.Organization) {
			errors = append(errors, ValidationError{
				Field:   prefix + ".organization",
				Message: fmt.Sprintf("organization '%s' not in allowed list", s.Organization),
			})
		}
	}

	if s.SpecURL == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".specUrl",
			Message: "specUrl is required",
		})
	} else if _, err := url.ParseRequestURI(s.SpecURL); err != nil {
		errors = append(errors, ValidationError{
			Field:   prefix + ".specUrl",
			Message: fmt.Sprintf("invalid URL: %v", err),
		})
	}

	// Validate status
	if !isValidStatus(s.Status) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".status",
			Message: fmt.Sprintf("invalid status: %s", s.Status),
		})
	}

	// Validate category
	if !isValidCategory(s.Category) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".category",
			Message: fmt.Sprintf("invalid category: %s", s.Category),
		})
	} else if len(v.AllowedCategories) > 0 {
		if !containsCategory(v.AllowedCategories, s.Category) {
			errors = append(errors, ValidationError{
				Field:   prefix + ".category",
				Message: fmt.Sprintf("category '%s' not in allowed list", s.Category),
			})
		}
	}

	// Validate layer
	if !isValidLayer(s.Layer) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".layer",
			Message: fmt.Sprintf("invalid layer: %s", s.Layer),
		})
	} else if len(v.AllowedLayers) > 0 {
		if !containsLayer(v.AllowedLayers, s.Layer) {
			errors = append(errors, ValidationError{
				Field:   prefix + ".layer",
				Message: fmt.Sprintf("layer '%s' not in allowed list", s.Layer),
			})
		}
	}

	// Optional validations
	if v.RequireDescription && s.Description == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".description",
			Message: "description is required",
		})
	}

	if v.RequireImplementations && len(s.Implementations) == 0 {
		errors = append(errors, ValidationError{
			Field:   prefix + ".implementations",
			Message: "at least one implementation is required",
		})
	}

	// Validate implementations
	for i, impl := range s.Implementations {
		implPrefix := fmt.Sprintf("%s.implementations[%d]", prefix, i)
		errors = append(errors, v.validateImplementation(implPrefix, &impl)...)
	}

	return errors
}

func (v *Validator) validateImplementation(prefix string, impl *Implementation) ValidationErrors {
	var errors ValidationErrors

	if impl.Name == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".name",
			Message: "name is required",
		})
	}

	if impl.URL == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".url",
			Message: "url is required",
		})
	} else if _, err := url.ParseRequestURI(impl.URL); err != nil {
		errors = append(errors, ValidationError{
			Field:   prefix + ".url",
			Message: fmt.Sprintf("invalid URL: %v", err),
		})
	}

	if impl.Status != "" && !isValidImplStatus(impl.Status) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".status",
			Message: fmt.Sprintf("invalid status: %s", impl.Status),
		})
	}

	return errors
}

func (v *Validator) validateCrossReferences(c *Catalog, validIDs map[string]bool) ValidationErrors {
	var errors ValidationErrors

	for i, std := range c.Standards {
		prefix := fmt.Sprintf("standards[%d]", i)

		// Check compatibleWith references
		for j, ref := range std.CompatibleWith {
			if !validIDs[ref] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("%s.compatibleWith[%d]", prefix, j),
					Message: fmt.Sprintf("references unknown standard: %s", ref),
				})
			}
		}

		// Check supersedes references
		for j, ref := range std.Supersedes {
			if !validIDs[ref] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("%s.supersedes[%d]", prefix, j),
					Message: fmt.Sprintf("references unknown standard: %s", ref),
				})
			}
		}

		// Check supersededBy references
		for j, ref := range std.SupersededBy {
			if !validIDs[ref] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("%s.supersededBy[%d]", prefix, j),
					Message: fmt.Sprintf("references unknown standard: %s", ref),
				})
			}
		}
	}

	return errors
}

func isValidStatus(s StandardStatus) bool {
	switch s {
	case StatusDraft, StatusProposed, StatusAdopted, StatusDeprecated, StatusRetired:
		return true
	}
	return false
}

func isValidCategory(c Category) bool {
	switch c {
	case CategoryAuthentication, CategoryAuthorization, CategoryIdentity,
		CategoryProvisioning, CategoryCommunication, CategoryDiscovery, CategoryGovernance:
		return true
	}
	return false
}

func isValidLayer(l IdentityLayer) bool {
	switch l {
	case LayerHuman, LayerAgent, LayerWorkload, LayerService:
		return true
	}
	return false
}

func isValidImplStatus(s ImplementationStatus) bool {
	switch s {
	case ImplStatusExperimental, ImplStatusBeta, ImplStatusStable, ImplStatusDeprecated:
		return true
	}
	return false
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func containsCategory(slice []Category, c Category) bool {
	for _, v := range slice {
		if v == c {
			return true
		}
	}
	return false
}

func containsLayer(slice []IdentityLayer, l IdentityLayer) bool {
	for _, v := range slice {
		if v == l {
			return true
		}
	}
	return false
}
