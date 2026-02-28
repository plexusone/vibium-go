package activity

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Registry holds all available activities.
type Registry struct {
	mu         sync.RWMutex
	activities map[string]Activity
}

// NewRegistry creates a new activity registry.
func NewRegistry() *Registry {
	return &Registry{
		activities: make(map[string]Activity),
	}
}

// Register adds an activity to the registry.
func (r *Registry) Register(a Activity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.activities[a.Name()] = a
}

// Get retrieves an activity by name.
func (r *Registry) Get(name string) (Activity, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.activities[name]
	return a, ok
}

// MustGet retrieves an activity by name, panicking if not found.
func (r *Registry) MustGet(name string) Activity {
	a, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("activity not found: %s", name))
	}
	return a
}

// List returns all registered activity names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.activities))
	for name := range r.activities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListByCategory returns activities grouped by category (prefix before ".").
func (r *Registry) ListByCategory() map[string][]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make(map[string][]string)
	for name := range r.activities {
		category := "other"
		if idx := strings.Index(name, "."); idx > 0 {
			category = name[:idx]
		}
		categories[category] = append(categories[category], name)
	}

	// Sort activities within each category
	for category := range categories {
		sort.Strings(categories[category])
	}

	return categories
}

// Count returns the number of registered activities.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.activities)
}

// Categories returns all activity categories.
func (r *Registry) Categories() []string {
	byCategory := r.ListByCategory()
	categories := make([]string, 0, len(byCategory))
	for category := range byCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// ListCategory returns all activities in a specific category.
func (r *Registry) ListCategory(category string) []string {
	byCategory := r.ListByCategory()
	if activities, ok := byCategory[category]; ok {
		return activities
	}
	return nil
}

// DefaultRegistry is the global default registry.
var DefaultRegistry = NewRegistry()

// Register adds an activity to the default registry.
func Register(a Activity) {
	DefaultRegistry.Register(a)
}

// Get retrieves an activity from the default registry.
func Get(name string) (Activity, bool) {
	return DefaultRegistry.Get(name)
}

// List returns all activity names from the default registry.
func List() []string {
	return DefaultRegistry.List()
}

// init registers all built-in activities.
func init() {
	// Browser activities
	Register(&NavigateActivity{})
	Register(&ClickActivity{})
	Register(&FillActivity{})
	Register(&TypeActivity{})
	Register(&SelectOptionActivity{})
	Register(&CheckActivity{})
	Register(&UncheckActivity{})
	Register(&ScrollActivity{})
	Register(&ScreenshotActivity{})
	Register(&PDFActivity{})

	// Element activities
	Register(&FindActivity{})
	Register(&FindAllActivity{})
	Register(&GetTextActivity{})
	Register(&GetValueActivity{})
	Register(&GetAttributeActivity{})
	Register(&WaitForActivity{})
	Register(&IsVisibleActivity{})

	// Data activities
	Register(&ScrapeTableActivity{})

	// File activities
	Register(&FileReadActivity{})
	Register(&FileWriteActivity{})
	Register(&FileExistsActivity{})
	Register(&FileDeleteActivity{})

	// HTTP activities
	Register(&HTTPGetActivity{})
	Register(&HTTPPostActivity{})
	Register(&HTTPDownloadActivity{})

	// Utility activities
	Register(&LogActivity{})
	Register(&WaitActivity{})
	Register(&AssertActivity{})
	Register(&SetVariableActivity{})
}
