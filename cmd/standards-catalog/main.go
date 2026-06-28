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
