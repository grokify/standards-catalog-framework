# Visualization

Interactive web visualization for exploring standards catalogs.

## Live Demo

<iframe src="standards-graph.html" width="100%" height="700" style="border: 1px solid #2a2a3a; border-radius: 8px;"></iframe>

!!! note "Interactive Features"
    - **Drag** nodes to rearrange the graph
    - **Scroll** to zoom in/out
    - **Click** nodes to open specification URLs
    - **Toggle** between Graph and Table views
    - **Filter** by layer, organization, category, or status

## Generating Visualizations

### From CLI

```bash
# Generate HTML with embedded catalog data
standards-catalog generate-web catalog.yaml -o standards-graph.html
```

### Custom Output Location

```bash
# Generate to docs folder for MkDocs
standards-catalog generate-web catalog.yaml -o docs/standards-graph.html
```

## Features

### Graph View

The force-directed graph shows:

- **Nodes**: Each standard as a colored circle
- **Links**: Relationships between standards
  - Dashed lines: `compatibleWith` relationships
  - Solid red lines: `supersedes` relationships

### Color Coding

Standards are colored by identity layer:

| Layer | Color |
|-------|-------|
| Human | :material-circle:{ style="color: #00d4ff" } Cyan |
| Agent | :material-circle:{ style="color: #ff6b00" } Orange |
| Workload | :material-circle:{ style="color: #a855f7" } Purple |
| Service | :material-circle:{ style="color: #22c55e" } Green |

### Table View

Toggle to Table view for a sortable list with columns:

- Standard name (clickable link to spec)
- Organization
- Layer
- Category
- Status
- Version

### Filtering

Use the sidebar filters to narrow down displayed standards:

- **Identity Layers**: Human, Agent, Workload, Service
- **Organizations**: IETF, CNCF, OpenID Foundation, etc.
- **Categories**: Authorization, Authentication, Identity, etc.
- **Status**: Draft, Proposed, Adopted, Deprecated, Retired

## Embedding in MkDocs

### Method 1: Direct HTML File

Copy the generated HTML to your `docs/` folder:

```bash
standards-catalog generate-web catalog.yaml -o docs/standards-graph.html
```

Link to it in navigation or markdown:

```markdown
[View Standards Graph](standards-graph.html)
```

### Method 2: Iframe Embed

Embed in any markdown page:

```markdown
<iframe
    src="standards-graph.html"
    width="100%"
    height="700"
    style="border: 1px solid #2a2a3a; border-radius: 8px;">
</iframe>
```

### Method 3: Full Page

Create a page that links directly to the visualization:

```yaml
# mkdocs.yml
nav:
  - Visualization: standards-graph.html
```

## Customization

### Editing Colors

The visualization uses CSS variables. Edit the `<style>` section:

```css
:root {
    --accent-human: #00d4ff;    /* Cyan for human layer */
    --accent-agent: #ff6b00;    /* Orange for agent layer */
    --accent-workload: #a855f7; /* Purple for workload layer */
    --accent-service: #22c55e;  /* Green for service layer */
}
```

### Modifying Data

The generated HTML contains embedded catalog data. To update:

1. Regenerate with updated catalog: `standards-catalog generate-web catalog.yaml`
2. Or edit the `catalogData` object in the HTML directly

### Loading External Data

Modify the HTML to load data from an external JSON file:

```javascript
document.addEventListener('DOMContentLoaded', () => {
    fetch('catalog.json')
        .then(response => response.json())
        .then(data => {
            Object.assign(catalogData, data);
            init();
        });
});
```

## Browser Support

The visualization requires a modern browser with JavaScript enabled:

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

Uses D3.js v7 loaded from CDN.
