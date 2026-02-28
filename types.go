// Package vibium provides a Go client for browser automation via the Vibium platform.
// It communicates with the Vibium clicker binary over WebSocket using the WebDriver BiDi protocol.
package vibium

import "time"

// BoundingBox represents the position and size of an element.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ElementInfo contains metadata about a DOM element.
type ElementInfo struct {
	Tag  string      `json:"tag"`
	Text string      `json:"text"`
	Box  BoundingBox `json:"box"`
}

// LaunchOptions configures browser launch behavior.
type LaunchOptions struct {
	// Headless runs the browser without a visible window.
	Headless bool

	// Port specifies the WebSocket port. If 0, an available port is auto-selected.
	Port int

	// ExecutablePath specifies a custom path to the clicker binary.
	ExecutablePath string
}

// FindOptions configures element finding behavior.
type FindOptions struct {
	// Timeout specifies how long to wait for the element to appear.
	// Default is 30 seconds.
	Timeout time.Duration

	// Semantic selectors for finding elements by accessibility properties.

	// Role matches elements by ARIA role (e.g., "button", "textbox").
	Role string

	// Text matches elements containing the specified text.
	Text string

	// Label matches elements by their associated label text.
	Label string

	// Placeholder matches input elements by placeholder attribute.
	Placeholder string

	// TestID matches elements by data-testid attribute.
	TestID string

	// Alt matches image elements by alt attribute.
	Alt string

	// Title matches elements by title attribute.
	Title string

	// XPath matches elements using an XPath expression.
	XPath string

	// Near finds elements near another element specified by selector.
	Near string
}

// SelectOptionValues specifies which options to select in a <select> element.
type SelectOptionValues struct {
	// Values selects options by their value attribute.
	Values []string

	// Labels selects options by their visible text.
	Labels []string

	// Indexes selects options by their zero-based index.
	Indexes []int
}

// Viewport represents the browser viewport dimensions.
type Viewport struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WindowState represents the browser window state.
type WindowState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	State     string `json:"state"` // "normal", "minimized", "maximized", "fullscreen"
	IsVisible bool   `json:"isVisible"`
}

// SetWindowOptions configures window state.
type SetWindowOptions struct {
	X      *int
	Y      *int
	Width  *int
	Height *int
	State  string // "normal", "minimized", "maximized", "fullscreen"
}

// PDFOptions configures PDF generation.
type PDFOptions struct {
	Path            string
	Scale           float64
	DisplayHeader   bool
	DisplayFooter   bool
	PrintBackground bool
	Landscape       bool
	PageRanges      string
	Format          string // "Letter", "Legal", "Tabloid", "A0"-"A6"
	Width           string
	Height          string
	Margin          *PDFMargin
}

// PDFMargin configures PDF page margins.
type PDFMargin struct {
	Top    string
	Right  string
	Bottom string
	Left   string
}

// FrameInfo contains metadata about a frame.
type FrameInfo struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// EmulateMediaOptions configures media emulation.
type EmulateMediaOptions struct {
	Media         string // "screen", "print", or ""
	ColorScheme   string // "light", "dark", "no-preference", or ""
	ReducedMotion string // "reduce", "no-preference", or ""
	ForcedColors  string // "active", "none", or ""
}

// Geolocation represents geographic coordinates.
type Geolocation struct {
	Latitude  float64
	Longitude float64
	Accuracy  float64
}

// Cookie represents a browser cookie.
type Cookie struct {
	Name         string  `json:"name"`
	Value        string  `json:"value"`
	Domain       string  `json:"domain"`
	Path         string  `json:"path"`
	Expires      float64 `json:"expires"`
	HTTPOnly     bool    `json:"httpOnly"`
	Secure       bool    `json:"secure"`
	SameSite     string  `json:"sameSite"`
	PartitionKey string  `json:"partitionKey,omitempty"`
}

// SetCookieParam represents parameters for setting a cookie.
type SetCookieParam struct {
	Name         string  `json:"name"`
	Value        string  `json:"value"`
	URL          string  `json:"url,omitempty"`
	Domain       string  `json:"domain,omitempty"`
	Path         string  `json:"path,omitempty"`
	Expires      float64 `json:"expires,omitempty"`
	HTTPOnly     bool    `json:"httpOnly,omitempty"`
	Secure       bool    `json:"secure,omitempty"`
	SameSite     string  `json:"sameSite,omitempty"`
	PartitionKey string  `json:"partitionKey,omitempty"`
}

// StorageState represents browser storage state including cookies and localStorage.
type StorageState struct {
	Cookies []Cookie             `json:"cookies"`
	Origins []StorageStateOrigin `json:"origins"`
}

// StorageStateOrigin represents localStorage for an origin.
type StorageStateOrigin struct {
	Origin       string            `json:"origin"`
	LocalStorage map[string]string `json:"localStorage"`
}

// ActionOptions configures action behavior (click, type).
type ActionOptions struct {
	// Timeout specifies how long to wait for actionability.
	// Default is 30 seconds.
	Timeout time.Duration
}

// DefaultTimeout is the default timeout for finding elements and waiting for actionability.
const DefaultTimeout = 30 * time.Second
