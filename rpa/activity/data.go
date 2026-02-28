package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	vibium "github.com/plexusone/vibium-go"
)

// ScrapeTableActivity extracts data from an HTML table.
type ScrapeTableActivity struct{}

func (a *ScrapeTableActivity) Name() string { return "data.scrapeTable" }

func (a *ScrapeTableActivity) Execute(ctx context.Context, params map[string]any, env *Environment) (any, error) {
	selector := GetString(params, "selector")
	if selector == "" {
		return nil, fmt.Errorf("selector parameter is required")
	}

	timeout := time.Duration(GetIntDefault(params, "timeout", 30000)) * time.Millisecond
	opts := &vibium.FindOptions{Timeout: timeout}

	// Find the table element
	el, err := env.Vibe.Find(ctx, selector, opts)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
	}
	_ = el // We use the selector directly in JavaScript

	// JavaScript to extract table data
	script := `
		(selector) => {
			const table = document.querySelector(selector);
			if (!table) return JSON.stringify({ error: 'Table not found' });

			const headers = [];
			const headerRow = table.querySelector('thead tr, tr:first-child');
			if (headerRow) {
				const headerCells = headerRow.querySelectorAll('th, td');
				headerCells.forEach(cell => {
					headers.push(cell.textContent.trim());
				});
			}

			const rows = [];
			const bodyRows = table.querySelectorAll('tbody tr, tr');
			const startIndex = headerRow ? 1 : 0;

			for (let i = startIndex; i < bodyRows.length; i++) {
				const row = bodyRows[i];
				const cells = row.querySelectorAll('td, th');
				const rowData = {};

				cells.forEach((cell, j) => {
					const key = headers[j] || 'col' + j;
					rowData[key] = cell.textContent.trim();
				});

				if (Object.keys(rowData).length > 0) {
					rows.push(rowData);
				}
			}

			return JSON.stringify({ headers: headers, rows: rows });
		}
	`

	result, err := env.Vibe.Evaluate(ctx, fmt.Sprintf("return (%s)('%s')", script, selector))
	if err != nil {
		return nil, fmt.Errorf("table extraction failed: %w", err)
	}

	// Parse the JSON result
	resultStr, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	var tableData struct {
		Error   string              `json:"error,omitempty"`
		Headers []string            `json:"headers"`
		Rows    []map[string]string `json:"rows"`
	}
	if err := json.Unmarshal([]byte(resultStr), &tableData); err != nil {
		return nil, fmt.Errorf("failed to parse table data: %w", err)
	}

	if tableData.Error != "" {
		return nil, fmt.Errorf("table extraction error: %s", tableData.Error)
	}

	// Return appropriate format based on params
	if GetBool(params, "rowsOnly") {
		return tableData.Rows, nil
	}

	return map[string]any{
		"headers": tableData.Headers,
		"rows":    tableData.Rows,
		"count":   len(tableData.Rows),
	}, nil
}
