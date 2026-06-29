// Command standards-catalog provides CLI tools for managing standards catalogs.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/grokify/standards-catalog-framework/catalog"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "standards-catalog",
	Short: "CLI for managing standards catalogs",
	Long: `standards-catalog is a CLI tool for creating, validating, and querying
standards catalogs. It supports JSON and YAML formats.`,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(generateWebCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate <catalog-file>",
	Short: "Validate a catalog file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, errors, err := catalog.LoadAndValidate(args[0])
		if err != nil {
			return fmt.Errorf("load catalog: %w", err)
		}

		if errors.HasErrors() {
			fmt.Fprintf(os.Stderr, "Validation failed with %d error(s):\n", len(errors))
			for _, e := range errors {
				fmt.Fprintf(os.Stderr, "  - %s\n", e.Error())
			}
			os.Exit(1)
		}

		fmt.Printf("Catalog '%s' is valid (%d standards)\n", c.Metadata.Name, len(c.Standards))
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list <catalog-file>",
	Short: "List all standards in a catalog",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tVERSION\tSTATUS\tCATEGORY\tLAYER")
		_, _ = fmt.Fprintln(w, "--\t----\t-------\t------\t--------\t-----")

		for _, s := range c.Standards {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				s.ID, s.Name, s.Version, s.Status, s.Category, s.Layer)
		}
		return w.Flush()
	},
}

var showCmd = &cobra.Command{
	Use:   "show <catalog-file> <standard-id>",
	Short: "Show details of a specific standard",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		s := c.FindByID(args[1])
		if s == nil {
			return fmt.Errorf("standard not found: %s", args[1])
		}

		// Pretty print as JSON
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var (
	queryCategory string
	queryLayer    string
	queryStatus   string
	queryOrg      string
	queryTag      string
)

func init() {
	queryCmd.Flags().StringVar(&queryCategory, "category", "", "Filter by category")
	queryCmd.Flags().StringVar(&queryLayer, "layer", "", "Filter by identity layer")
	queryCmd.Flags().StringVar(&queryStatus, "status", "", "Filter by status")
	queryCmd.Flags().StringVar(&queryOrg, "org", "", "Filter by organization")
	queryCmd.Flags().StringVar(&queryTag, "tag", "", "Filter by tag")
}

var queryCmd = &cobra.Command{
	Use:   "query <catalog-file>",
	Short: "Query standards with filters",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		results := c.Standards

		// Apply filters
		if queryCategory != "" {
			results = filterByCategory(results, catalog.Category(queryCategory))
		}
		if queryLayer != "" {
			results = filterByLayer(results, catalog.IdentityLayer(queryLayer))
		}
		if queryStatus != "" {
			results = filterByStatus(results, catalog.StandardStatus(queryStatus))
		}
		if queryOrg != "" {
			results = filterByOrg(results, queryOrg)
		}
		if queryTag != "" {
			results = filterByTag(results, queryTag)
		}

		if len(results) == 0 {
			fmt.Println("No standards found matching criteria")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tVERSION\tSTATUS\tCATEGORY\tLAYER")
		_, _ = fmt.Fprintln(w, "--\t----\t-------\t------\t--------\t-----")

		for _, s := range results {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				s.ID, s.Name, s.Version, s.Status, s.Category, s.Layer)
		}
		return w.Flush()
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats <catalog-file>",
	Short: "Show catalog statistics",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Catalog: %s (v%s)\n", c.Metadata.Name, c.Metadata.Version)
		fmt.Printf("Total standards: %d\n\n", len(c.Standards))

		// Count by status
		statusCounts := make(map[catalog.StandardStatus]int)
		for _, s := range c.Standards {
			statusCounts[s.Status]++
		}
		fmt.Println("By Status:")
		for status, count := range statusCounts {
			fmt.Printf("  %s: %d\n", status, count)
		}

		// Count by category
		categoryCounts := make(map[catalog.Category]int)
		for _, s := range c.Standards {
			categoryCounts[s.Category]++
		}
		fmt.Println("\nBy Category:")
		for cat, count := range categoryCounts {
			fmt.Printf("  %s: %d\n", cat, count)
		}

		// Count by layer
		layerCounts := make(map[catalog.IdentityLayer]int)
		for _, s := range c.Standards {
			layerCounts[s.Layer]++
		}
		fmt.Println("\nBy Layer:")
		for layer, count := range layerCounts {
			fmt.Printf("  %s: %d\n", layer, count)
		}

		// Organizations
		fmt.Printf("\nOrganizations: %s\n", strings.Join(c.GetOrganizations(), ", "))

		// Tags
		tags := c.GetTags()
		if len(tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(tags, ", "))
		}

		return nil
	},
}

func filterByCategory(standards []catalog.Standard, cat catalog.Category) []catalog.Standard {
	var results []catalog.Standard
	for _, s := range standards {
		if s.Category == cat {
			results = append(results, s)
		}
	}
	return results
}

func filterByLayer(standards []catalog.Standard, layer catalog.IdentityLayer) []catalog.Standard {
	var results []catalog.Standard
	for _, s := range standards {
		if s.Layer == layer {
			results = append(results, s)
		}
	}
	return results
}

func filterByStatus(standards []catalog.Standard, status catalog.StandardStatus) []catalog.Standard {
	var results []catalog.Standard
	for _, s := range standards {
		if s.Status == status {
			results = append(results, s)
		}
	}
	return results
}

func filterByOrg(standards []catalog.Standard, org string) []catalog.Standard {
	var results []catalog.Standard
	for _, s := range standards {
		if s.Organization == org {
			results = append(results, s)
		}
	}
	return results
}

func filterByTag(standards []catalog.Standard, tag string) []catalog.Standard {
	var results []catalog.Standard
	for _, s := range standards {
		for _, t := range s.Tags {
			if t == tag {
				results = append(results, s)
				break
			}
		}
	}
	return results
}

var exportCmd = &cobra.Command{
	Use:   "export <catalog-file>",
	Short: "Export catalog as JSON",
	Long:  `Export the catalog as JSON for use with web visualizations or other tools.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var generateWebOutput string

func init() {
	generateWebCmd.Flags().StringVarP(&generateWebOutput, "output", "o", "standards-graph.html", "Output HTML file path")
}

var generateWebCmd = &cobra.Command{
	Use:   "generate-web <catalog-file>",
	Short: "Generate interactive web visualization",
	Long: `Generate a standalone HTML file with an interactive visualization
of the standards catalog. The HTML file can be opened directly in a browser
or embedded in MkDocs.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load(args[0])
		if err != nil {
			return err
		}

		// Convert catalog to JSON
		catalogJSON, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("marshal catalog: %w", err)
		}

		// Generate HTML with embedded data
		html := generateVisualizationHTML(string(catalogJSON))

		// Write to output file (0644 is appropriate for HTML files served by web servers)
		if err := os.WriteFile(generateWebOutput, []byte(html), 0644); err != nil { //nolint:gosec // HTML files need to be readable
			return fmt.Errorf("write output: %w", err)
		}

		fmt.Printf("Generated visualization: %s\n", generateWebOutput)
		return nil
	},
}

func generateVisualizationHTML(catalogJSON string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Standards Catalog Visualization</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        :root {
            --bg-primary: #0a0a0f;
            --bg-secondary: #12121a;
            --bg-tertiary: #1a1a24;
            --text-primary: #e8e8e8;
            --text-secondary: #a0a0a0;
            --border-color: #2a2a3a;
            --accent-human: #00d4ff;
            --accent-agent: #ff6b00;
            --accent-workload: #a855f7;
            --accent-service: #22c55e;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
        }
        .container { display: flex; height: 100vh; }
        .sidebar {
            width: 280px;
            background: var(--bg-secondary);
            border-right: 1px solid var(--border-color);
            padding: 20px;
            overflow-y: auto;
            flex-shrink: 0;
        }
        .sidebar h2 {
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-secondary);
            margin-bottom: 12px;
            padding-bottom: 8px;
            border-bottom: 1px solid var(--border-color);
        }
        .filter-section { margin-bottom: 24px; }
        .filter-group { margin-bottom: 8px; }
        .filter-group label {
            display: flex;
            align-items: center;
            gap: 8px;
            cursor: pointer;
            padding: 6px 8px;
            border-radius: 4px;
            transition: background 0.2s;
        }
        .filter-group label:hover { background: var(--bg-tertiary); }
        .filter-group input[type="checkbox"] {
            width: 16px;
            height: 16px;
            accent-color: var(--accent-human);
        }
        .color-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            flex-shrink: 0;
        }
        .color-dot.human { background: var(--accent-human); }
        .color-dot.agent { background: var(--accent-agent); }
        .color-dot.workload { background: var(--accent-workload); }
        .color-dot.service { background: var(--accent-service); }
        .main-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
        .header {
            padding: 16px 24px;
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border-color);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header h1 { font-size: 20px; font-weight: 600; }
        .stats { display: flex; gap: 24px; font-size: 13px; color: var(--text-secondary); }
        .stats span { color: var(--text-primary); font-weight: 600; }
        .graph-container { flex: 1; position: relative; overflow: hidden; }
        #graph { width: 100%; height: 100%; }
        .node { cursor: pointer; transition: opacity 0.2s; }
        .node:hover { opacity: 0.8; }
        .node-label {
            font-size: 11px;
            fill: var(--text-primary);
            pointer-events: none;
            text-anchor: middle;
        }
        .link { stroke-opacity: 0.4; transition: stroke-opacity 0.2s; }
        .link.compatible { stroke: #4a5568; stroke-dasharray: 4,2; }
        .link.supersedes { stroke: #ef4444; }
        .tooltip {
            position: absolute;
            background: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 16px;
            max-width: 320px;
            pointer-events: none;
            opacity: 0;
            transition: opacity 0.2s;
            z-index: 1000;
            box-shadow: 0 4px 20px rgba(0,0,0,0.4);
        }
        .tooltip.visible { opacity: 1; }
        .tooltip h3 { font-size: 16px; margin-bottom: 8px; display: flex; align-items: center; gap: 8px; }
        .tooltip .org { font-size: 12px; color: var(--text-secondary); margin-bottom: 12px; }
        .tooltip .description { font-size: 13px; line-height: 1.5; margin-bottom: 12px; color: var(--text-secondary); }
        .tooltip .meta { display: flex; flex-wrap: wrap; gap: 8px; }
        .tooltip .tag { font-size: 11px; padding: 3px 8px; border-radius: 4px; background: var(--bg-secondary); border: 1px solid var(--border-color); }
        .tooltip .status { font-size: 11px; padding: 3px 8px; border-radius: 4px; font-weight: 500; }
        .tooltip .status.adopted { background: rgba(34, 197, 94, 0.2); color: #22c55e; }
        .tooltip .status.draft { background: rgba(234, 179, 8, 0.2); color: #eab308; }
        .tooltip .status.proposed { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
        .tooltip .status.deprecated { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
        .tooltip .status.retired { background: rgba(107, 114, 128, 0.2); color: #6b7280; }
        .tooltip .link-btn {
            display: inline-block;
            margin-top: 12px;
            padding: 6px 12px;
            background: var(--accent-human);
            color: var(--bg-primary);
            text-decoration: none;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 500;
            pointer-events: auto;
        }
        .tooltip .link-btn:hover { opacity: 0.9; }
        .legend {
            position: absolute;
            bottom: 20px;
            right: 20px;
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 12px 16px;
        }
        .legend h4 { font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; color: var(--text-secondary); margin-bottom: 8px; }
        .legend-item { display: flex; align-items: center; gap: 8px; font-size: 12px; margin-bottom: 4px; }
        .legend-line { width: 24px; height: 2px; }
        .legend-line.compatible { background: repeating-linear-gradient(90deg, #4a5568 0px, #4a5568 4px, transparent 4px, transparent 6px); }
        .legend-line.supersedes { background: #ef4444; }
        .view-toggle { display: flex; gap: 4px; background: var(--bg-tertiary); padding: 4px; border-radius: 6px; }
        .view-toggle button {
            padding: 6px 12px;
            border: none;
            background: transparent;
            color: var(--text-secondary);
            cursor: pointer;
            border-radius: 4px;
            font-size: 13px;
            transition: all 0.2s;
        }
        .view-toggle button.active { background: var(--accent-human); color: var(--bg-primary); }
        .view-toggle button:hover:not(.active) { color: var(--text-primary); }
        .matrix-container { display: none; padding: 24px; overflow: auto; height: 100%; }
        .matrix-container.visible { display: block; }
        .matrix-table { border-collapse: collapse; width: 100%; }
        .matrix-table th, .matrix-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid var(--border-color); }
        .matrix-table th {
            background: var(--bg-tertiary);
            font-weight: 500;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-secondary);
            position: sticky;
            top: 0;
        }
        .matrix-table tr:hover td { background: var(--bg-tertiary); }
        .matrix-table .standard-name { font-weight: 500; }
        .matrix-table .standard-name a { color: var(--accent-human); text-decoration: none; }
        .matrix-table .standard-name a:hover { text-decoration: underline; }
        .layer-badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 500; }
        .layer-badge.human { background: rgba(0, 212, 255, 0.2); color: var(--accent-human); }
        .layer-badge.agent { background: rgba(255, 107, 0, 0.2); color: var(--accent-agent); }
        .layer-badge.workload { background: rgba(168, 85, 247, 0.2); color: var(--accent-workload); }
        .layer-badge.service { background: rgba(34, 197, 94, 0.2); color: var(--accent-service); }
    </style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="filter-section">
                <h2>Identity Layers</h2>
                <div class="filter-group" id="layer-filters"></div>
            </div>
            <div class="filter-section">
                <h2>Organizations</h2>
                <div class="filter-group" id="org-filters"></div>
            </div>
            <div class="filter-section">
                <h2>Categories</h2>
                <div class="filter-group" id="category-filters"></div>
            </div>
            <div class="filter-section">
                <h2>Status</h2>
                <div class="filter-group" id="status-filters"></div>
            </div>
        </aside>
        <main class="main-content">
            <header class="header">
                <h1 id="catalog-title">Standards Catalog</h1>
                <div style="display: flex; align-items: center; gap: 24px;">
                    <div class="stats">
                        <div>Standards: <span id="standard-count">0</span></div>
                        <div>Organizations: <span id="org-count">0</span></div>
                    </div>
                    <div class="view-toggle">
                        <button class="active" data-view="graph">Graph</button>
                        <button data-view="matrix">Table</button>
                    </div>
                </div>
            </header>
            <div class="graph-container" id="graph-view">
                <svg id="graph"></svg>
                <div class="tooltip" id="tooltip"></div>
                <div class="legend">
                    <h4>Relationships</h4>
                    <div class="legend-item"><div class="legend-line compatible"></div><span>Compatible With</span></div>
                    <div class="legend-item"><div class="legend-line supersedes"></div><span>Supersedes</span></div>
                </div>
            </div>
            <div class="matrix-container" id="matrix-view">
                <table class="matrix-table">
                    <thead><tr><th>Standard</th><th>Organization</th><th>Layer</th><th>Category</th><th>Status</th><th>Version</th></tr></thead>
                    <tbody id="matrix-body"></tbody>
                </table>
            </div>
        </main>
    </div>
    <script>
        const catalogData = ` + catalogJSON + `;
        const layerColors = { human: '#00d4ff', agent: '#ff6b00', workload: '#a855f7', service: '#22c55e' };
        let filters = { layers: new Set(['human', 'agent', 'workload', 'service']), organizations: new Set(), categories: new Set(), statuses: new Set() };
        let simulation;

        function init() {
            document.getElementById('catalog-title').textContent = catalogData.metadata.name || 'Standards Catalog';
            const organizations = [...new Set(catalogData.standards.map(s => s.organization))].sort();
            const categories = [...new Set(catalogData.standards.map(s => s.category))].sort();
            const statuses = [...new Set(catalogData.standards.map(s => s.status))];
            const layers = ['human', 'agent', 'workload', 'service'];
            filters.organizations = new Set(organizations);
            filters.categories = new Set(categories);
            filters.statuses = new Set(statuses);
            createFilters('layer-filters', layers, 'layers', true);
            createFilters('org-filters', organizations, 'organizations');
            createFilters('category-filters', categories, 'categories');
            createFilters('status-filters', statuses, 'statuses');
            document.getElementById('standard-count').textContent = catalogData.standards.length;
            document.getElementById('org-count').textContent = organizations.length;
            createGraph();
            createMatrix();
            document.querySelectorAll('.view-toggle button').forEach(btn => {
                btn.addEventListener('click', () => {
                    document.querySelectorAll('.view-toggle button').forEach(b => b.classList.remove('active'));
                    btn.classList.add('active');
                    const view = btn.dataset.view;
                    document.getElementById('graph-view').style.display = view === 'graph' ? 'block' : 'none';
                    document.getElementById('matrix-view').classList.toggle('visible', view === 'matrix');
                });
            });
        }

        function createFilters(containerId, values, filterKey, isLayer = false) {
            const container = document.getElementById(containerId);
            values.forEach(value => {
                const label = document.createElement('label');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.checked = true;
                checkbox.addEventListener('change', () => {
                    checkbox.checked ? filters[filterKey].add(value) : filters[filterKey].delete(value);
                    updateVisualization();
                });
                label.appendChild(checkbox);
                if (isLayer) {
                    const dot = document.createElement('span');
                    dot.className = 'color-dot ' + value;
                    label.appendChild(dot);
                }
                label.appendChild(document.createTextNode(capitalize(value)));
                container.appendChild(label);
            });
        }

        function capitalize(str) { return str.charAt(0).toUpperCase() + str.slice(1); }
        function getFilteredStandards() {
            return catalogData.standards.filter(s =>
                filters.layers.has(s.layer) && filters.organizations.has(s.organization) &&
                filters.categories.has(s.category) && filters.statuses.has(s.status)
            );
        }

        function createGraph() {
            const svg = d3.select('#graph');
            const container = document.getElementById('graph-view');
            const width = container.clientWidth;
            const height = container.clientHeight;
            svg.attr('width', width).attr('height', height);
            svg.append('g').attr('class', 'links');
            svg.append('g').attr('class', 'nodes');
            const zoom = d3.zoom().scaleExtent([0.3, 3]).on('zoom', (e) => svg.selectAll('g').attr('transform', e.transform));
            svg.call(zoom);
            updateGraph();
            window.addEventListener('resize', () => {
                svg.attr('width', container.clientWidth).attr('height', container.clientHeight);
                if (simulation) { simulation.force('center', d3.forceCenter(container.clientWidth / 2, container.clientHeight / 2)); simulation.alpha(0.3).restart(); }
            });
        }

        function updateGraph() {
            const svg = d3.select('#graph');
            const width = parseInt(svg.attr('width'));
            const height = parseInt(svg.attr('height'));
            const filteredStandards = getFilteredStandards();
            const filteredIds = new Set(filteredStandards.map(s => s.id));
            const nodes = filteredStandards.map(s => ({ ...s, radius: 24 }));
            const links = [];
            filteredStandards.forEach(s => {
                (s.compatibleWith || []).forEach(t => { if (filteredIds.has(t)) links.push({ source: s.id, target: t, type: 'compatible' }); });
                (s.supersedes || []).forEach(t => { if (filteredIds.has(t)) links.push({ source: s.id, target: t, type: 'supersedes' }); });
            });
            if (simulation) simulation.stop();
            simulation = d3.forceSimulation(nodes)
                .force('link', d3.forceLink(links).id(d => d.id).distance(120))
                .force('charge', d3.forceManyBody().strength(-400))
                .force('center', d3.forceCenter(width / 2, height / 2))
                .force('collision', d3.forceCollide().radius(40));
            const linkSel = svg.select('.links').selectAll('.link').data(links, d => d.source.id + '-' + d.target.id);
            linkSel.exit().remove();
            const linkEnter = linkSel.enter().append('line').attr('class', d => 'link ' + d.type).attr('stroke-width', 2);
            const linkMerge = linkEnter.merge(linkSel);
            const nodeSel = svg.select('.nodes').selectAll('.node').data(nodes, d => d.id);
            nodeSel.exit().remove();
            const nodeEnter = nodeSel.enter().append('g').attr('class', 'node').call(d3.drag().on('start', dragStarted).on('drag', dragged).on('end', dragEnded));
            nodeEnter.append('circle').attr('r', d => d.radius).attr('fill', d => layerColors[d.layer]).attr('stroke', '#fff').attr('stroke-width', 2);
            nodeEnter.append('text').attr('class', 'node-label').attr('dy', 4).text(d => d.name.length > 12 ? d.name.substring(0, 10) + '...' : d.name);
            const nodeMerge = nodeEnter.merge(nodeSel);
            nodeMerge.on('mouseover', showTooltip).on('mouseout', hideTooltip).on('click', (e, d) => { if (d.specUrl) window.open(d.specUrl, '_blank'); });
            simulation.on('tick', () => {
                linkMerge.attr('x1', d => d.source.x).attr('y1', d => d.source.y).attr('x2', d => d.target.x).attr('y2', d => d.target.y);
                nodeMerge.attr('transform', d => 'translate(' + d.x + ',' + d.y + ')');
            });
        }

        function dragStarted(e, d) { if (!e.active) simulation.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; }
        function dragged(e, d) { d.fx = e.x; d.fy = e.y; }
        function dragEnded(e, d) { if (!e.active) simulation.alphaTarget(0); d.fx = null; d.fy = null; }

        function showTooltip(event, d) {
            const tooltip = document.getElementById('tooltip');
            tooltip.innerHTML = '<h3><span class="color-dot ' + d.layer + '" style="display:inline-block;width:10px;height:10px;border-radius:50%;"></span> ' + d.name + '</h3>' +
                '<div class="org">' + d.organization + ' &bull; v' + d.version + '</div>' +
                '<div class="description">' + (d.description || 'No description.') + '</div>' +
                '<div class="meta"><span class="status ' + d.status + '">' + capitalize(d.status) + '</span><span class="tag">' + capitalize(d.layer) + '</span><span class="tag">' + capitalize(d.category) + '</span></div>' +
                '<a href="' + d.specUrl + '" class="link-btn" target="_blank">View Specification</a>';
            const rect = event.target.getBoundingClientRect();
            const containerRect = document.getElementById('graph-view').getBoundingClientRect();
            let left = rect.right - containerRect.left + 10;
            let top = rect.top - containerRect.top;
            if (left + 320 > containerRect.width) left = rect.left - containerRect.left - 330;
            if (top + 200 > containerRect.height) top = containerRect.height - 210;
            tooltip.style.left = left + 'px';
            tooltip.style.top = top + 'px';
            tooltip.classList.add('visible');
        }

        function hideTooltip() { document.getElementById('tooltip').classList.remove('visible'); }
        function createMatrix() { updateMatrix(); }
        function updateMatrix() {
            const tbody = document.getElementById('matrix-body');
            tbody.innerHTML = getFilteredStandards().map(s =>
                '<tr><td class="standard-name"><a href="' + s.specUrl + '" target="_blank">' + s.name + '</a></td>' +
                '<td>' + s.organization + '</td><td><span class="layer-badge ' + s.layer + '">' + capitalize(s.layer) + '</span></td>' +
                '<td>' + capitalize(s.category) + '</td><td><span class="status ' + s.status + '">' + capitalize(s.status) + '</span></td><td>' + s.version + '</td></tr>'
            ).join('');
        }
        function updateVisualization() {
            document.getElementById('standard-count').textContent = getFilteredStandards().length;
            updateGraph();
            updateMatrix();
        }
        document.addEventListener('DOMContentLoaded', init);
    </script>
</body>
</html>`
}
