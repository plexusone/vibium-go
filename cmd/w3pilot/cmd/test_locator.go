package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var (
	locatorTimeout  time.Duration
	locatorStrategy string
)

// generateLocatorCmd generates a locator for an element
var generateLocatorCmd = &cobra.Command{
	Use:   "generate-locator <selector>",
	Short: "Generate robust locator for element",
	Long: `Generate a robust locator string for an element.

Strategies:
  css    - Generate unique CSS selector (default)
  xpath  - Generate XPath expression
  testid - Use data-testid attribute
  role   - Use ARIA role and name
  text   - Use text content

Examples:
  w3pilot test generate-locator "#submit"
  w3pilot test generate-locator ".button" --strategy xpath
  w3pilot test generate-locator "button" --strategy role
  w3pilot test generate-locator --format json "#login"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), locatorTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		info, err := pilot.GenerateLocator(ctx, selector, &w3pilot.GenerateLocatorOptions{
			Strategy: locatorStrategy,
			Timeout:  locatorTimeout,
		})
		if err != nil {
			return fmt.Errorf("generate locator failed: %w", err)
		}

		Output(info, func(data interface{}) string {
			i := data.(*w3pilot.LocatorInfo)
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Locator: %s\n", i.Locator))
			sb.WriteString(fmt.Sprintf("Strategy: %s\n", i.Strategy))
			if len(i.Metadata) > 0 {
				sb.WriteString("Metadata:\n")
				for k, v := range i.Metadata {
					sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
				}
			}
			return sb.String()
		})

		return nil
	},
}

func init() {
	testCmd.AddCommand(generateLocatorCmd)

	generateLocatorCmd.Flags().DurationVar(&locatorTimeout, "timeout", 5*time.Second, "Timeout")
	generateLocatorCmd.Flags().StringVar(&locatorStrategy, "strategy", "css", "Strategy: css, xpath, testid, role, text")
}
