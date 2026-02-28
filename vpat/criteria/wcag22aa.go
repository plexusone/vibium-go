// Package criteria defines WCAG success criteria and their axe-core rule mappings.
package criteria

// Criterion defines a WCAG success criterion.
type Criterion struct {
	ID          string   // e.g., "1.1.1"
	Name        string   // e.g., "Non-text Content"
	Level       string   // "A", "AA", or "AAA"
	Description string   // Brief description
	AxeRules    []string // Mapped axe-core rules
	CanAutomate bool     // Whether this can be fully automated
}

// WCAG22AA returns all WCAG 2.2 Level A and AA criteria.
func WCAG22AA() []Criterion {
	return []Criterion{
		// Principle 1: Perceivable
		// Guideline 1.1 Text Alternatives
		{
			ID:          "1.1.1",
			Name:        "Non-text Content",
			Level:       "A",
			Description: "All non-text content has a text alternative",
			AxeRules:    []string{"image-alt", "input-image-alt", "area-alt", "object-alt", "svg-img-alt"},
			CanAutomate: true,
		},

		// Guideline 1.2 Time-based Media
		{
			ID:          "1.2.1",
			Name:        "Audio-only and Video-only (Prerecorded)",
			Level:       "A",
			Description: "Alternatives for prerecorded audio-only and video-only content",
			AxeRules:    []string{"video-caption", "audio-caption"},
			CanAutomate: false,
		},
		{
			ID:          "1.2.2",
			Name:        "Captions (Prerecorded)",
			Level:       "A",
			Description: "Captions are provided for prerecorded audio content",
			AxeRules:    []string{"video-caption"},
			CanAutomate: false,
		},
		{
			ID:          "1.2.3",
			Name:        "Audio Description or Media Alternative (Prerecorded)",
			Level:       "A",
			Description: "Alternative or audio description for prerecorded video",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.2.4",
			Name:        "Captions (Live)",
			Level:       "AA",
			Description: "Captions are provided for live audio content",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.2.5",
			Name:        "Audio Description (Prerecorded)",
			Level:       "AA",
			Description: "Audio description for prerecorded video content",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Guideline 1.3 Adaptable
		{
			ID:          "1.3.1",
			Name:        "Info and Relationships",
			Level:       "A",
			Description: "Information and relationships conveyed through presentation can be programmatically determined",
			AxeRules:    []string{"definition-list", "dlitem", "list", "listitem", "table-fake-caption", "td-headers-attr", "th-has-data-cells", "empty-table-header", "scope-attr-valid", "p-as-heading"},
			CanAutomate: true,
		},
		{
			ID:          "1.3.2",
			Name:        "Meaningful Sequence",
			Level:       "A",
			Description: "Correct reading sequence can be programmatically determined",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.3.3",
			Name:        "Sensory Characteristics",
			Level:       "A",
			Description: "Instructions don't rely solely on sensory characteristics",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.3.4",
			Name:        "Orientation",
			Level:       "AA",
			Description: "Content does not restrict its view to a single orientation",
			AxeRules:    []string{"css-orientation-lock"},
			CanAutomate: true,
		},
		{
			ID:          "1.3.5",
			Name:        "Identify Input Purpose",
			Level:       "AA",
			Description: "Input field purpose can be programmatically determined",
			AxeRules:    []string{"autocomplete-valid"},
			CanAutomate: true,
		},

		// Guideline 1.4 Distinguishable
		{
			ID:          "1.4.1",
			Name:        "Use of Color",
			Level:       "A",
			Description: "Color is not the only visual means of conveying information",
			AxeRules:    []string{"link-in-text-block"},
			CanAutomate: false,
		},
		{
			ID:          "1.4.2",
			Name:        "Audio Control",
			Level:       "A",
			Description: "Mechanism to pause or stop audio that plays automatically",
			AxeRules:    []string{"no-autoplay-audio"},
			CanAutomate: true,
		},
		{
			ID:          "1.4.3",
			Name:        "Contrast (Minimum)",
			Level:       "AA",
			Description: "Text has a contrast ratio of at least 4.5:1",
			AxeRules:    []string{"color-contrast"},
			CanAutomate: true,
		},
		{
			ID:          "1.4.4",
			Name:        "Resize Text",
			Level:       "AA",
			Description: "Text can be resized up to 200% without loss of functionality",
			AxeRules:    []string{"meta-viewport"},
			CanAutomate: false,
		},
		{
			ID:          "1.4.5",
			Name:        "Images of Text",
			Level:       "AA",
			Description: "Text is used to convey information rather than images of text",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.4.10",
			Name:        "Reflow",
			Level:       "AA",
			Description: "Content can reflow without horizontal scrolling at 320 CSS pixels",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.4.11",
			Name:        "Non-text Contrast",
			Level:       "AA",
			Description: "UI components and graphics have a contrast ratio of at least 3:1",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.4.12",
			Name:        "Text Spacing",
			Level:       "AA",
			Description: "No loss of content when text spacing is adjusted",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "1.4.13",
			Name:        "Content on Hover or Focus",
			Level:       "AA",
			Description: "Additional content triggered by hover/focus is dismissible, hoverable, and persistent",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Principle 2: Operable
		// Guideline 2.1 Keyboard Accessible
		{
			ID:          "2.1.1",
			Name:        "Keyboard",
			Level:       "A",
			Description: "All functionality is operable via keyboard",
			AxeRules:    []string{"scrollable-region-focusable"},
			CanAutomate: false,
		},
		{
			ID:          "2.1.2",
			Name:        "No Keyboard Trap",
			Level:       "A",
			Description: "Keyboard focus can be moved away from any component",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.1.4",
			Name:        "Character Key Shortcuts",
			Level:       "A",
			Description: "Single character key shortcuts can be turned off or remapped",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Guideline 2.2 Enough Time
		{
			ID:          "2.2.1",
			Name:        "Timing Adjustable",
			Level:       "A",
			Description: "Time limits can be turned off, adjusted, or extended",
			AxeRules:    []string{"meta-refresh"},
			CanAutomate: false,
		},
		{
			ID:          "2.2.2",
			Name:        "Pause, Stop, Hide",
			Level:       "A",
			Description: "Moving, blinking, scrolling content can be paused, stopped, or hidden",
			AxeRules:    []string{"blink", "marquee"},
			CanAutomate: true,
		},

		// Guideline 2.3 Seizures and Physical Reactions
		{
			ID:          "2.3.1",
			Name:        "Three Flashes or Below Threshold",
			Level:       "A",
			Description: "No content flashes more than three times per second",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Guideline 2.4 Navigable
		{
			ID:          "2.4.1",
			Name:        "Bypass Blocks",
			Level:       "A",
			Description: "Mechanism to bypass repeated blocks of content",
			AxeRules:    []string{"bypass", "region"},
			CanAutomate: true,
		},
		{
			ID:          "2.4.2",
			Name:        "Page Titled",
			Level:       "A",
			Description: "Pages have titles that describe topic or purpose",
			AxeRules:    []string{"document-title"},
			CanAutomate: true,
		},
		{
			ID:          "2.4.3",
			Name:        "Focus Order",
			Level:       "A",
			Description: "Focus order preserves meaning and operability",
			AxeRules:    []string{"tabindex"},
			CanAutomate: false,
		},
		{
			ID:          "2.4.4",
			Name:        "Link Purpose (In Context)",
			Level:       "A",
			Description: "Link purpose can be determined from link text or context",
			AxeRules:    []string{"link-name"},
			CanAutomate: true,
		},
		{
			ID:          "2.4.5",
			Name:        "Multiple Ways",
			Level:       "AA",
			Description: "More than one way to locate a page within a set",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.4.6",
			Name:        "Headings and Labels",
			Level:       "AA",
			Description: "Headings and labels describe topic or purpose",
			AxeRules:    []string{"empty-heading"},
			CanAutomate: false,
		},
		{
			ID:          "2.4.7",
			Name:        "Focus Visible",
			Level:       "AA",
			Description: "Keyboard focus indicator is visible",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.4.11",
			Name:        "Focus Not Obscured (Minimum)",
			Level:       "AA",
			Description: "Focused element is not entirely hidden by other content",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Guideline 2.5 Input Modalities
		{
			ID:          "2.5.1",
			Name:        "Pointer Gestures",
			Level:       "A",
			Description: "Multipoint or path-based gestures have single-pointer alternatives",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.5.2",
			Name:        "Pointer Cancellation",
			Level:       "A",
			Description: "Single-pointer functionality can be cancelled",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.5.3",
			Name:        "Label in Name",
			Level:       "A",
			Description: "Visible label is part of accessible name",
			AxeRules:    []string{"label-content-name-mismatch"},
			CanAutomate: true,
		},
		{
			ID:          "2.5.4",
			Name:        "Motion Actuation",
			Level:       "A",
			Description: "Motion-triggered functionality can be disabled and has alternatives",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.5.7",
			Name:        "Dragging Movements",
			Level:       "AA",
			Description: "Dragging functionality has single-pointer alternatives",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "2.5.8",
			Name:        "Target Size (Minimum)",
			Level:       "AA",
			Description: "Touch targets are at least 24x24 CSS pixels",
			AxeRules:    []string{"target-size"},
			CanAutomate: true,
		},

		// Principle 3: Understandable
		// Guideline 3.1 Readable
		{
			ID:          "3.1.1",
			Name:        "Language of Page",
			Level:       "A",
			Description: "Default human language can be programmatically determined",
			AxeRules:    []string{"html-has-lang", "html-lang-valid"},
			CanAutomate: true,
		},
		{
			ID:          "3.1.2",
			Name:        "Language of Parts",
			Level:       "AA",
			Description: "Language of parts can be programmatically determined",
			AxeRules:    []string{"valid-lang"},
			CanAutomate: true,
		},

		// Guideline 3.2 Predictable
		{
			ID:          "3.2.1",
			Name:        "On Focus",
			Level:       "A",
			Description: "Focus does not trigger unexpected context changes",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.2.2",
			Name:        "On Input",
			Level:       "A",
			Description: "Input does not trigger unexpected context changes",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.2.3",
			Name:        "Consistent Navigation",
			Level:       "AA",
			Description: "Navigation mechanisms are consistent across pages",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.2.4",
			Name:        "Consistent Identification",
			Level:       "AA",
			Description: "Components with same functionality are identified consistently",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.2.6",
			Name:        "Consistent Help",
			Level:       "A",
			Description: "Help mechanisms are in consistent locations",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Guideline 3.3 Input Assistance
		{
			ID:          "3.3.1",
			Name:        "Error Identification",
			Level:       "A",
			Description: "Input errors are identified and described in text",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.3.2",
			Name:        "Labels or Instructions",
			Level:       "A",
			Description: "Labels or instructions are provided for user input",
			AxeRules:    []string{"label", "select-name", "input-button-name"},
			CanAutomate: true,
		},
		{
			ID:          "3.3.3",
			Name:        "Error Suggestion",
			Level:       "AA",
			Description: "Suggestions are provided when input errors are detected",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.3.4",
			Name:        "Error Prevention (Legal, Financial, Data)",
			Level:       "AA",
			Description: "Submissions are reversible, verifiable, or confirmable",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.3.7",
			Name:        "Redundant Entry",
			Level:       "A",
			Description: "Previously entered information is auto-populated or available for selection",
			AxeRules:    []string{},
			CanAutomate: false,
		},
		{
			ID:          "3.3.8",
			Name:        "Accessible Authentication (Minimum)",
			Level:       "AA",
			Description: "Authentication does not require cognitive function test",
			AxeRules:    []string{},
			CanAutomate: false,
		},

		// Principle 4: Robust
		// Guideline 4.1 Compatible
		{
			ID:          "4.1.1",
			Name:        "Parsing",
			Level:       "A",
			Description: "No major parsing errors (obsolete in WCAG 2.2)",
			AxeRules:    []string{"duplicate-id", "duplicate-id-active", "duplicate-id-aria"},
			CanAutomate: true,
		},
		{
			ID:          "4.1.2",
			Name:        "Name, Role, Value",
			Level:       "A",
			Description: "Name, role, and value can be programmatically determined",
			AxeRules: []string{
				"aria-allowed-attr", "aria-allowed-role", "aria-command-name",
				"aria-dialog-name", "aria-hidden-body", "aria-hidden-focus",
				"aria-input-field-name", "aria-meter-name", "aria-progressbar-name",
				"aria-required-attr", "aria-required-children", "aria-required-parent",
				"aria-roledescription", "aria-roles", "aria-toggle-field-name",
				"aria-tooltip-name", "aria-valid-attr", "aria-valid-attr-value",
				"button-name", "form-field-multiple-labels", "frame-title",
				"input-button-name", "role-img-alt",
			},
			CanAutomate: true,
		},
		{
			ID:          "4.1.3",
			Name:        "Status Messages",
			Level:       "AA",
			Description: "Status messages can be programmatically determined",
			AxeRules:    []string{"aria-live-region-attr"},
			CanAutomate: false,
		},
	}
}

// GetByLevel filters criteria by WCAG level.
func GetByLevel(criteria []Criterion, levels ...string) []Criterion {
	levelSet := make(map[string]bool)
	for _, l := range levels {
		levelSet[l] = true
	}

	var result []Criterion
	for _, c := range criteria {
		if levelSet[c.Level] {
			result = append(result, c)
		}
	}
	return result
}

// GetAutomatable returns criteria that can be fully automated.
func GetAutomatable(criteria []Criterion) []Criterion {
	var result []Criterion
	for _, c := range criteria {
		if c.CanAutomate {
			result = append(result, c)
		}
	}
	return result
}
