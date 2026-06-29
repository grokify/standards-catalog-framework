# Standards Catalog Visualization

Interactive JavaScript visualization for standards catalogs that can be embedded in MkDocs or any HTML page.

## Files

- `standards-graph.html` - Full-page interactive visualization with graph and table views
- `catalog.json` - Generated JSON data file (create using CLI)

## Features

- **Graph View**: Force-directed network showing standards and their relationships
  - Color-coded by identity layer (Human, Agent, Workload, Service)
  - Shows compatibility and supersedes relationships
  - Drag nodes to rearrange
  - Zoom and pan support
  - Click nodes to open specification

- **Table View**: Sortable table with all standard details

- **Filtering**: Filter by layer, organization, category, and status

## Usage

### Option 1: Standalone HTML Page

1. Generate JSON from your catalog:

```bash
standards-catalog show your-catalog.yaml --all > web/catalog.json
```

2. Open `standards-graph.html` in a browser

3. To load external data, modify the script to fetch your JSON:

```javascript
loadCatalogData('catalog.json');
```

### Option 2: Embed in MkDocs

#### Method A: Copy HTML to docs folder

1. Copy `standards-graph.html` to your MkDocs `docs/` folder
2. Link to it from your navigation or markdown:

```markdown
[View Standards Graph](standards-graph.html)
```

MkDocs will copy HTML files directly to the output.

#### Method B: Use iframe in Markdown

```markdown
<iframe
    src="standards-graph.html"
    width="100%"
    height="800"
    style="border: none; border-radius: 8px;"
></iframe>
```

#### Method C: Include in custom theme

Add to your MkDocs theme's extra files and reference from templates.

## Customizing Data

The visualization uses embedded sample data by default. To use your own catalog:

### Option 1: Edit Embedded Data

Replace the `catalogData` object in the HTML file with your catalog JSON.

### Option 2: Load External JSON

Modify the initialization to fetch your catalog:

```javascript
document.addEventListener('DOMContentLoaded', () => {
    loadCatalogData('path/to/your/catalog.json');
});
```

### Option 3: Generate with CLI

```bash
# Convert YAML to JSON
standards-catalog show catalog.yaml --format json > catalog.json
```

## Styling

The visualization uses a dark theme matching the AIStandards.io design system:

- **Human Layer**: Cyan (#00d4ff)
- **Agent Layer**: Orange (#ff6b00)
- **Workload Layer**: Purple (#a855f7)
- **Service Layer**: Green (#22c55e)

To customize colors, edit the CSS variables in the `<style>` section:

```css
:root {
    --accent-human: #00d4ff;
    --accent-agent: #ff6b00;
    --accent-workload: #a855f7;
    --accent-service: #22c55e;
}
```

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

Requires JavaScript enabled. Uses D3.js v7 loaded from CDN.

## Data Format

The visualization expects catalog data in this JSON format:

```json
{
    "apiVersion": "standards-catalog/v1",
    "kind": "Catalog",
    "metadata": {
        "name": "Catalog Name",
        "version": "1.0.0"
    },
    "standards": [
        {
            "id": "standard-id",
            "name": "Standard Name",
            "version": "1.0",
            "status": "adopted",
            "organization": "Organization Name",
            "specUrl": "https://...",
            "category": "authorization",
            "layer": "agent",
            "description": "Description text",
            "compatibleWith": ["other-id"],
            "tags": ["tag1", "tag2"]
        }
    ]
}
```
