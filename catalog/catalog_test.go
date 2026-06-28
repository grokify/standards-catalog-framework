package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if c.APIVersion != "standards-catalog/v1" {
		t.Errorf("APIVersion = %q, want %q", c.APIVersion, "standards-catalog/v1")
	}
	if c.Kind != "Catalog" {
		t.Errorf("Kind = %q, want %q", c.Kind, "Catalog")
	}
	if c.Metadata.Name != "AI Agent Standards Catalog" {
		t.Errorf("Metadata.Name = %q, want %q", c.Metadata.Name, "AI Agent Standards Catalog")
	}
	if len(c.Standards) != 5 {
		t.Errorf("len(Standards) = %d, want 5", len(c.Standards))
	}
}

func TestLoadJSON(t *testing.T) {
	c, err := Load("../testdata/example-catalog.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if c.Metadata.Name != "Minimal Catalog" {
		t.Errorf("Metadata.Name = %q, want %q", c.Metadata.Name, "Minimal Catalog")
	}
	if len(c.Standards) != 1 {
		t.Errorf("len(Standards) = %d, want 1", len(c.Standards))
	}
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Load should fail for nonexistent file")
	}
}

func TestLoadUnsupportedFormat(t *testing.T) {
	// Create temp file with unsupported extension
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.Write([]byte("test"))
	_ = tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Load should fail for unsupported format")
	}
}

func TestParse(t *testing.T) {
	jsonData := []byte(`{
		"apiVersion": "standards-catalog/v1",
		"kind": "Catalog",
		"metadata": {"name": "Test", "version": "1.0"},
		"standards": []
	}`)

	c, err := Parse(jsonData, ".json")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if c.Metadata.Name != "Test" {
		t.Errorf("Metadata.Name = %q, want %q", c.Metadata.Name, "Test")
	}
}

func TestParseYAML(t *testing.T) {
	yamlData := []byte(`
apiVersion: standards-catalog/v1
kind: Catalog
metadata:
  name: Test
  version: "1.0"
standards: []
`)

	c, err := Parse(yamlData, ".yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if c.Metadata.Name != "Test" {
		t.Errorf("Metadata.Name = %q, want %q", c.Metadata.Name, "Test")
	}

	// Test .yml extension too
	c, err = Parse(yamlData, ".yml")
	if err != nil {
		t.Fatalf("Parse .yml failed: %v", err)
	}
	if c.Metadata.Name != "Test" {
		t.Errorf("Metadata.Name = %q, want %q", c.Metadata.Name, "Test")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse([]byte("invalid json"), ".json")
	if err == nil {
		t.Error("Parse should fail for invalid JSON")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := Parse([]byte(":\n  invalid: yaml:"), ".yaml")
	if err == nil {
		t.Error("Parse should fail for invalid YAML")
	}
}

func TestSaveAndLoad(t *testing.T) {
	original := &Catalog{
		APIVersion: "standards-catalog/v1",
		Kind:       "Catalog",
		Metadata: CatalogMetadata{
			Name:    "Test Catalog",
			Version: "1.0.0",
		},
		Standards: []Standard{
			{
				ID:           "test-1",
				Name:         "Test Standard",
				Version:      "1.0",
				Status:       StatusAdopted,
				Organization: "Test Org",
				SpecURL:      "https://example.com",
				Category:     CategoryAuthentication,
				Layer:        LayerService,
			},
		},
	}

	// Test JSON save/load
	jsonPath := filepath.Join(t.TempDir(), "test.json")
	if err := Save(original, jsonPath); err != nil {
		t.Fatalf("Save JSON failed: %v", err)
	}

	loaded, err := Load(jsonPath)
	if err != nil {
		t.Fatalf("Load JSON failed: %v", err)
	}
	if loaded.Metadata.Name != original.Metadata.Name {
		t.Errorf("Loaded name = %q, want %q", loaded.Metadata.Name, original.Metadata.Name)
	}
	if len(loaded.Standards) != 1 {
		t.Errorf("Loaded standards count = %d, want 1", len(loaded.Standards))
	}

	// Test YAML save/load
	yamlPath := filepath.Join(t.TempDir(), "test.yaml")
	if err := Save(original, yamlPath); err != nil {
		t.Fatalf("Save YAML failed: %v", err)
	}

	loaded, err = Load(yamlPath)
	if err != nil {
		t.Fatalf("Load YAML failed: %v", err)
	}
	if loaded.Metadata.Name != original.Metadata.Name {
		t.Errorf("Loaded name = %q, want %q", loaded.Metadata.Name, original.Metadata.Name)
	}
}

func TestSaveUnsupportedFormat(t *testing.T) {
	c := &Catalog{}
	err := Save(c, "test.txt")
	if err == nil {
		t.Error("Save should fail for unsupported format")
	}
}

func TestFindByID(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Find existing
	s := c.FindByID("oauth2")
	if s == nil {
		t.Fatal("FindByID should find oauth2")
	}
	if s.Name != "OAuth 2.0" {
		t.Errorf("Name = %q, want %q", s.Name, "OAuth 2.0")
	}

	// Find non-existing
	s = c.FindByID("nonexistent")
	if s != nil {
		t.Error("FindByID should return nil for nonexistent")
	}
}

func TestFindByCategory(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	authStandards := c.FindByCategory(CategoryAuthentication)
	if len(authStandards) != 2 {
		t.Errorf("len(auth standards) = %d, want 2", len(authStandards))
	}

	// Check we got the right ones
	ids := make(map[string]bool)
	for _, s := range authStandards {
		ids[s.ID] = true
	}
	if !ids["oidc"] || !ids["aauth"] {
		t.Error("Expected to find oidc and aauth in authentication category")
	}
}

func TestFindByLayer(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	humanLayer := c.FindByLayer(LayerHuman)
	if len(humanLayer) != 3 {
		t.Errorf("len(human layer) = %d, want 3", len(humanLayer))
	}

	agentLayer := c.FindByLayer(LayerAgent)
	if len(agentLayer) != 1 {
		t.Errorf("len(agent layer) = %d, want 1", len(agentLayer))
	}
	if agentLayer[0].ID != "aauth" {
		t.Errorf("agent layer standard ID = %q, want %q", agentLayer[0].ID, "aauth")
	}
}

func TestFindByStatus(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	adopted := c.FindByStatus(StatusAdopted)
	if len(adopted) != 3 {
		t.Errorf("len(adopted) = %d, want 3", len(adopted))
	}

	draft := c.FindByStatus(StatusDraft)
	if len(draft) != 1 {
		t.Errorf("len(draft) = %d, want 1", len(draft))
	}
}

func TestFindByTag(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	tokenTag := c.FindByTag("token")
	if len(tokenTag) != 1 {
		t.Errorf("len(token tag) = %d, want 1", len(tokenTag))
	}
	if tokenTag[0].ID != "oauth2" {
		t.Errorf("token tag standard ID = %q, want %q", tokenTag[0].ID, "oauth2")
	}
}

func TestFindByOrganization(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	ietf := c.FindByOrganization("IETF")
	if len(ietf) != 2 {
		t.Errorf("len(IETF) = %d, want 2", len(ietf))
	}
}

func TestGetOrganizations(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	orgs := c.GetOrganizations()
	if len(orgs) != 4 {
		t.Errorf("len(orgs) = %d, want 4", len(orgs))
	}

	orgSet := make(map[string]bool)
	for _, org := range orgs {
		orgSet[org] = true
	}
	expected := []string{"IETF", "OpenID Foundation", "CNCF", "OAIAF"}
	for _, e := range expected {
		if !orgSet[e] {
			t.Errorf("Expected org %q not found", e)
		}
	}
}

func TestGetTags(t *testing.T) {
	c, err := Load("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	tags := c.GetTags()
	if len(tags) == 0 {
		t.Error("Expected some tags")
	}

	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}
	if !tagSet["token"] {
		t.Error("Expected 'token' tag")
	}
	if !tagSet["agent"] {
		t.Error("Expected 'agent' tag")
	}
}

func TestLoadAndValidate(t *testing.T) {
	// Valid catalog
	c, errors, err := LoadAndValidate("../testdata/example-catalog.yaml")
	if err != nil {
		t.Fatalf("LoadAndValidate failed: %v", err)
	}
	if errors.HasErrors() {
		t.Errorf("Expected no validation errors, got: %v", errors)
	}
	if c.Metadata.Name != "AI Agent Standards Catalog" {
		t.Errorf("Name = %q, want %q", c.Metadata.Name, "AI Agent Standards Catalog")
	}

	// Invalid catalog
	_, errors, err = LoadAndValidate("../testdata/invalid-catalog.yaml")
	if err != nil {
		t.Fatalf("LoadAndValidate failed: %v", err)
	}
	if !errors.HasErrors() {
		t.Error("Expected validation errors for invalid catalog")
	}
}
