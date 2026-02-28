package cmd

import (
	"fmt"
	"sort"

	"github.com/agentplexus/vibium-go/rpa/activity"
	"github.com/spf13/cobra"
)

var listCategory string

var listCmd = &cobra.Command{
	Use:   "list [resource]",
	Short: "List available resources",
	Long: `List available activities and other resources.

Examples:
  # List all activities
  vibium-rpa list activities

  # List activities in a category
  vibium-rpa list activities --category browser
`,
}

var listActivitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List available activities",
	Long: `List all available activities that can be used in workflows.

Activities are grouped by category (e.g., browser, element, file, http, util).

Examples:
  # List all activities
  vibium-rpa list activities

  # List only browser activities
  vibium-rpa list activities --category browser
`,
	RunE: listActivities,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listActivitiesCmd)

	listActivitiesCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter by category")
}

func listActivities(cmd *cobra.Command, args []string) error {
	byCategory := activity.DefaultRegistry.ListByCategory()

	// Get sorted categories
	categories := make([]string, 0, len(byCategory))
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	if listCategory != "" {
		// Show single category
		activities, ok := byCategory[listCategory]
		if !ok {
			return fmt.Errorf("unknown category: %s", listCategory)
		}

		fmt.Printf("Category: %s (%d activities)\n\n", listCategory, len(activities))
		for _, a := range activities {
			fmt.Printf("  %s\n", a)
		}
		return nil
	}

	// Show all categories
	total := 0
	for _, activities := range byCategory {
		total += len(activities)
	}
	fmt.Printf("Available Activities (%d total)\n\n", total)

	for _, category := range categories {
		activities := byCategory[category]
		fmt.Printf("%s (%d):\n", category, len(activities))
		for _, a := range activities {
			fmt.Printf("  - %s\n", a)
		}
		fmt.Println()
	}

	fmt.Println("Use 'vibium-rpa list activities --category <name>' to filter by category.")

	return nil
}
