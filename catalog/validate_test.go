package catalog

import (
	"strings"
	"testing"
)

func TestValidationError(t *testing.T) {
	e := &ValidationError{
		Field:   "standards[0].id",
		Message: "id is required",
	}
	if e.Error() != "standards[0].id: id is required" {
		t.Errorf("Error() = %q, want %q", e.Error(), "standards[0].id: id is required")
	}
}

func TestValidationErrors(t *testing.T) {
	var errors ValidationErrors

	// Empty errors
	if errors.HasErrors() {
		t.Error("Empty errors should not have errors")
	}
	if errors.Error() != "" {
		t.Error("Empty errors should return empty string")
	}

	// Add errors
	errors = append(errors, ValidationError{Field: "a", Message: "error 1"})
	errors = append(errors, ValidationError{Field: "b", Message: "error 2"})

	if !errors.HasErrors() {
		t.Error("Should have errors")
	}
	if !strings.Contains(errors.Error(), "a: error 1") {
		t.Error("Error string should contain first error")
	}
	if !strings.Contains(errors.Error(), "b: error 2") {
		t.Error("Error string should contain second error")
	}
}

func TestValidateValidCatalog(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata: CatalogMetadata{
			Name:    "Test",
			Version: "1.0",
		},
		Standards: []Standard{
			{
				ID:           "test-standard",
				Name:         "Test Standard",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test Org",
				SpecURL:      "https://example.com/spec",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if errors.HasErrors() {
		t.Errorf("Valid catalog should have no errors: %v", errors)
	}
}

func TestValidateAPIVersion(t *testing.T) {
	c := &Catalog{
		APIVersion: "wrong-version",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for wrong API version")
	}

	hasAPIVersionError := false
	for _, e := range errors {
		if e.Field == "apiVersion" {
			hasAPIVersionError = true
			break
		}
	}
	if !hasAPIVersionError {
		t.Error("Expected apiVersion error")
	}
}

func TestValidateKind(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "WrongKind",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for wrong Kind")
	}
}

func TestValidateMetadata(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "", Version: ""},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected errors for missing metadata")
	}

	hasNameError := false
	hasVersionError := false
	for _, e := range errors {
		if e.Field == "metadata.name" {
			hasNameError = true
		}
		if e.Field == "metadata.version" {
			hasVersionError = true
		}
	}
	if !hasNameError {
		t.Error("Expected metadata.name error")
	}
	if !hasVersionError {
		t.Error("Expected metadata.version error")
	}
}

func TestValidateStandardID(t *testing.T) {
	tests := []struct {
		id      string
		wantErr bool
	}{
		{"valid-id", false},
		{"oauth2", false},
		{"test123", false},
		{"a", false},
		{"INVALID", true},    // uppercase
		{"Invalid-Id", true}, // mixed case
		{"123invalid", true}, // starts with number
		{"-invalid", true},   // starts with hyphen
		{"invalid_id", true}, // underscore
		{"invalid id", true}, // space
		{"", true},           // empty
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			c := &Catalog{
				APIVersion: "standards-catalog/v1",
				Kind:       "Catalog",
				Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
				Standards: []Standard{
					{
						ID:           tt.id,
						Name:         "Test",
						Version:      "1.0",
						Status:       StatusAdopted,
						Organization: "Test",
						SpecURL:      "https://example.com",
						Category:     CategoryAuthentication,
						Layer:        LayerService,
					},
				},
			}

			v := NewValidator()
			errors := v.Validate(c)
			hasIDError := false
			for _, e := range errors {
				if strings.Contains(e.Field, ".id") {
					hasIDError = true
					break
				}
			}
			if hasIDError != tt.wantErr {
				t.Errorf("id=%q: hasIDError=%v, wantErr=%v", tt.id, hasIDError, tt.wantErr)
			}
		})
	}
}

func TestValidateDuplicateID(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "duplicate",
				Name:         "Test 1",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
			{
				ID:           "duplicate",
				Name:         "Test 2",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for duplicate ID")
	}

	hasDuplicateError := false
	for _, e := range errors {
		if strings.Contains(e.Message, "duplicate") {
			hasDuplicateError = true
			break
		}
	}
	if !hasDuplicateError {
		t.Error("Expected duplicate ID error")
	}
}

func TestValidateRequiredFields(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				// All required fields missing
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected errors for missing required fields")
	}

	requiredFields := []string{"id", "name", "version", "organization", "specUrl", "status", "category", "layer"}
	for _, field := range requiredFields {
		found := false
		for _, e := range errors {
			if strings.Contains(e.Field, field) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error for missing %s", field)
		}
	}
}

func TestValidateInvalidURL(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "not-a-url",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for invalid URL")
	}

	hasURLError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "specUrl") {
			hasURLError = true
			break
		}
	}
	if !hasURLError {
		t.Error("Expected specUrl error")
	}
}

func TestValidateInvalidStatus(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StandardStatus("invalid"),
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for invalid status")
	}
}

func TestValidateInvalidCategory(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     Category("invalid"),
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for invalid category")
	}
}

func TestValidateInvalidLayer(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        IdentityLayer("invalid"),
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if !errors.HasErrors() {
		t.Error("Expected error for invalid layer")
	}
}

func TestValidatorRequireDescription(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
				Description:  "", // Missing description
			},
		},
	}

	// Without RequireDescription
	v := NewValidator()
	errors := v.Validate(c)
	hasDescError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "description") {
			hasDescError = true
		}
	}
	if hasDescError {
		t.Error("Should not require description by default")
	}

	// With RequireDescription
	v.RequireDescription = true
	errors = v.Validate(c)
	hasDescError = false
	for _, e := range errors {
		if strings.Contains(e.Field, "description") {
			hasDescError = true
		}
	}
	if !hasDescError {
		t.Error("Should require description when RequireDescription is true")
	}
}

func TestValidatorRequireImplementations(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:              "test",
				Name:            "Test",
				Version:         "1.0",
				Status:          StatusAdopted,
				Organization:    "Test",
				SpecURL:         "https://example.com",
				Category:        CategoryAuthentication,
				Layer:           LayerService,
				Implementations: nil, // No implementations
			},
		},
	}

	// Without RequireImplementations
	v := NewValidator()
	errors := v.Validate(c)
	hasImplError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "implementations") {
			hasImplError = true
		}
	}
	if hasImplError {
		t.Error("Should not require implementations by default")
	}

	// With RequireImplementations
	v.RequireImplementations = true
	errors = v.Validate(c)
	hasImplError = false
	for _, e := range errors {
		if strings.Contains(e.Field, "implementations") {
			hasImplError = true
		}
	}
	if !hasImplError {
		t.Error("Should require implementations when RequireImplementations is true")
	}
}

func TestValidatorAllowedOrganizations(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Unknown Org",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	v.AllowedOrganizations = []string{"IETF", "W3C"}
	errors := v.Validate(c)

	hasOrgError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "organization") && strings.Contains(e.Message, "not in allowed list") {
			hasOrgError = true
		}
	}
	if !hasOrgError {
		t.Error("Should reject organization not in allowed list")
	}
}

func TestValidatorAllowedCategories(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryGovernance,
				Layer:        LayerService,
			},
		},
	}

	v := NewValidator()
	v.AllowedCategories = []Category{CategoryAuthentication, CategoryAuthorization}
	errors := v.Validate(c)

	hasCatError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "category") && strings.Contains(e.Message, "not in allowed list") {
			hasCatError = true
		}
	}
	if !hasCatError {
		t.Error("Should reject category not in allowed list")
	}
}

func TestValidatorAllowedLayers(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerWorkload,
			},
		},
	}

	v := NewValidator()
	v.AllowedLayers = []IdentityLayer{LayerHuman, LayerAgent}
	errors := v.Validate(c)

	hasLayerError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "layer") && strings.Contains(e.Message, "not in allowed list") {
			hasLayerError = true
		}
	}
	if !hasLayerError {
		t.Error("Should reject layer not in allowed list")
	}
}

func TestValidateImplementation(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "test",
				Name:         "Test",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
				Implementations: []Implementation{
					{
						Name:   "",          // Missing name
						URL:    "not-a-url", // Invalid URL
						Status: "invalid",   // Invalid status
					},
				},
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)

	hasNameError := false
	hasURLError := false
	hasStatusError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "implementations[0].name") {
			hasNameError = true
		}
		if strings.Contains(e.Field, "implementations[0].url") {
			hasURLError = true
		}
		if strings.Contains(e.Field, "implementations[0].status") {
			hasStatusError = true
		}
	}
	if !hasNameError {
		t.Error("Expected implementation name error")
	}
	if !hasURLError {
		t.Error("Expected implementation url error")
	}
	if !hasStatusError {
		t.Error("Expected implementation status error")
	}
}

func TestValidateCrossReferences(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:             "test",
				Name:           "Test",
				Version:        "1.0",
				Status:         StatusAdopted,
				Organization:   "Test",
				SpecURL:        "https://example.com",
				Category:       CategoryAuthentication,
				Layer:          LayerService,
				CompatibleWith: []string{"nonexistent1"},
				Supersedes:     []string{"nonexistent2"},
				SupersededBy:   []string{"nonexistent3"},
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)

	hasCompatibleError := false
	hasSupersedesError := false
	hasSupersededByError := false
	for _, e := range errors {
		if strings.Contains(e.Field, "compatibleWith") && strings.Contains(e.Message, "nonexistent1") {
			hasCompatibleError = true
		}
		if strings.Contains(e.Field, "supersedes") && strings.Contains(e.Message, "nonexistent2") {
			hasSupersedesError = true
		}
		if strings.Contains(e.Field, "supersededBy") && strings.Contains(e.Message, "nonexistent3") {
			hasSupersededByError = true
		}
	}
	if !hasCompatibleError {
		t.Error("Expected compatibleWith cross-reference error")
	}
	if !hasSupersedesError {
		t.Error("Expected supersedes cross-reference error")
	}
	if !hasSupersededByError {
		t.Error("Expected supersededBy cross-reference error")
	}
}

func TestValidateCrossReferencesValid(t *testing.T) {
	c := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata:   CatalogMetadata{Name: "Test", Version: "1.0"},
		Standards: []Standard{
			{
				ID:           "standard-a",
				Name:         "Standard A",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test",
				SpecURL:      "https://example.com/a",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
				SupersededBy: []string{"standard-b"},
			},
			{
				ID:             "standard-b",
				Name:           "Standard B",
				Version:        "2.0",
				Status:         StatusAdopted,
				Organization:   "Test",
				SpecURL:        "https://example.com/b",
				Category:       CategoryAuthentication,
				Layer:          LayerService,
				Supersedes:     []string{"standard-a"},
				CompatibleWith: []string{"standard-a"},
			},
		},
	}

	v := NewValidator()
	errors := v.Validate(c)
	if errors.HasErrors() {
		t.Errorf("Valid cross-references should have no errors: %v", errors)
	}
}

func TestIsValidStatus(t *testing.T) {
	validStatuses := []StandardStatus{StatusDraft, StatusProposed, StatusAdopted, StatusDeprecated, StatusRetired}
	for _, s := range validStatuses {
		if !isValidStatus(s) {
			t.Errorf("Status %q should be valid", s)
		}
	}
	if isValidStatus(StandardStatus("invalid")) {
		t.Error("Invalid status should not be valid")
	}
}

func TestIsValidCategory(t *testing.T) {
	validCategories := []Category{
		CategoryAuthentication, CategoryAuthorization, CategoryIdentity,
		CategoryProvisioning, CategoryCommunication, CategoryDiscovery, CategoryGovernance,
	}
	for _, c := range validCategories {
		if !isValidCategory(c) {
			t.Errorf("Category %q should be valid", c)
		}
	}
	if isValidCategory(Category("invalid")) {
		t.Error("Invalid category should not be valid")
	}
}

func TestIsValidLayer(t *testing.T) {
	validLayers := []IdentityLayer{LayerHuman, LayerAgent, LayerWorkload, LayerService}
	for _, l := range validLayers {
		if !isValidLayer(l) {
			t.Errorf("Layer %q should be valid", l)
		}
	}
	if isValidLayer(IdentityLayer("invalid")) {
		t.Error("Invalid layer should not be valid")
	}
}

func TestIsValidImplStatus(t *testing.T) {
	validStatuses := []ImplementationStatus{ImplStatusExperimental, ImplStatusBeta, ImplStatusStable, ImplStatusDeprecated}
	for _, s := range validStatuses {
		if !isValidImplStatus(s) {
			t.Errorf("ImplStatus %q should be valid", s)
		}
	}
	if isValidImplStatus(ImplementationStatus("invalid")) {
		t.Error("Invalid impl status should not be valid")
	}
}
