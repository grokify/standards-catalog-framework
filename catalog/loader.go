package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load loads a catalog from a file.
// Supports JSON and YAML formats based on file extension.
func Load(path string) (*Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return Parse(data, filepath.Ext(path))
}

// Parse parses catalog data from bytes.
// Extension should include the dot (e.g., ".json", ".yaml").
func Parse(data []byte, ext string) (*Catalog, error) {
	var c Catalog

	switch strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &c); err != nil {
			return nil, fmt.Errorf("parse YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}

	return &c, nil
}

// Save saves a catalog to a file.
// Format is determined by file extension.
func Save(c *Catalog, path string) error {
	var data []byte
	var err error

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		data, err = json.MarshalIndent(c, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(c)
	default:
		return fmt.Errorf("unsupported format: %s", filepath.Ext(path))
	}

	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return os.WriteFile(path, data, 0o600)
}

// LoadAndValidate loads a catalog and validates it.
func LoadAndValidate(path string) (*Catalog, ValidationErrors, error) {
	c, err := Load(path)
	if err != nil {
		return nil, nil, err
	}

	v := NewValidator()
	errors := v.Validate(c)

	return c, errors, nil
}

// FindByID finds a standard by its ID.
func (c *Catalog) FindByID(id string) *Standard {
	for i := range c.Standards {
		if c.Standards[i].ID == id {
			return &c.Standards[i]
		}
	}
	return nil
}

// FindByCategory returns all standards in a category.
func (c *Catalog) FindByCategory(cat Category) []Standard {
	var results []Standard
	for _, s := range c.Standards {
		if s.Category == cat {
			results = append(results, s)
		}
	}
	return results
}

// FindByLayer returns all standards at a given layer.
func (c *Catalog) FindByLayer(layer IdentityLayer) []Standard {
	var results []Standard
	for _, s := range c.Standards {
		if s.Layer == layer {
			results = append(results, s)
		}
	}
	return results
}

// FindByStatus returns all standards with a given status.
func (c *Catalog) FindByStatus(status StandardStatus) []Standard {
	var results []Standard
	for _, s := range c.Standards {
		if s.Status == status {
			results = append(results, s)
		}
	}
	return results
}

// FindByTag returns all standards with a given tag.
func (c *Catalog) FindByTag(tag string) []Standard {
	var results []Standard
	for _, s := range c.Standards {
		for _, t := range s.Tags {
			if t == tag {
				results = append(results, s)
				break
			}
		}
	}
	return results
}

// FindByOrganization returns all standards from a given organization.
func (c *Catalog) FindByOrganization(org string) []Standard {
	var results []Standard
	for _, s := range c.Standards {
		if s.Organization == org {
			results = append(results, s)
		}
	}
	return results
}

// GetOrganizations returns all unique organizations in the catalog.
func (c *Catalog) GetOrganizations() []string {
	seen := make(map[string]bool)
	var orgs []string
	for _, s := range c.Standards {
		if !seen[s.Organization] {
			seen[s.Organization] = true
			orgs = append(orgs, s.Organization)
		}
	}
	return orgs
}

// GetTags returns all unique tags in the catalog.
func (c *Catalog) GetTags() []string {
	seen := make(map[string]bool)
	var tags []string
	for _, s := range c.Standards {
		for _, t := range s.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	return tags
}
