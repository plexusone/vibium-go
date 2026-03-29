package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	w3pilot "github.com/plexusone/w3pilot"
)

// Manager handles named state snapshots.
type Manager struct {
	dir string
}

// StateInfo contains metadata about a saved state.
type StateInfo struct {
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	Size       int64     `json:"size"`
	Origins    []string  `json:"origins,omitempty"` // List of origin domains
	NumCookies int       `json:"num_cookies"`
}

// NewManager creates a new state manager.
// If dir is empty, uses ~/.w3pilot/states/
func NewManager(dir string) (*Manager, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dir = filepath.Join(home, ".w3pilot", "states")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create states directory: %w", err)
	}

	return &Manager{dir: dir}, nil
}

// Save saves a storage state with the given name.
func (m *Manager) Save(name string, state *w3pilot.StorageState) error {
	if name == "" {
		return fmt.Errorf("state name cannot be empty")
	}

	// Validate name (alphanumeric, dash, underscore)
	for _, c := range name {
		if !isValidNameChar(c) {
			return fmt.Errorf("invalid state name: only alphanumeric, dash, and underscore allowed")
		}
	}

	path := m.statePath(name)

	// Wrap state with metadata
	wrapper := struct {
		Name      string                `json:"name"`
		CreatedAt time.Time             `json:"created_at"`
		State     *w3pilot.StorageState `json:"state"`
	}{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		State:     state,
	}

	data, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// Load loads a storage state by name.
func (m *Manager) Load(name string) (*w3pilot.StorageState, error) {
	path := m.statePath(name)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("state not found: %s", name)
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var wrapper struct {
		State *w3pilot.StorageState `json:"state"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	if wrapper.State == nil {
		return nil, fmt.Errorf("invalid state file: missing state data")
	}

	return wrapper.State, nil
}

// List returns information about all saved states.
func (m *Manager) List() ([]StateInfo, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []StateInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read states directory: %w", err)
	}

	var states []StateInfo

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // Remove .json

		info, err := m.getStateInfo(name)
		if err != nil {
			continue // Skip invalid files
		}

		states = append(states, *info)
	}

	// Sort by created time, newest first
	sort.Slice(states, func(i, j int) bool {
		return states[i].CreatedAt.After(states[j].CreatedAt)
	})

	return states, nil
}

// Delete removes a saved state by name.
func (m *Manager) Delete(name string) error {
	path := m.statePath(name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("state not found: %s", name)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete state: %w", err)
	}

	return nil
}

// Exists checks if a state with the given name exists.
func (m *Manager) Exists(name string) bool {
	path := m.statePath(name)
	_, err := os.Stat(path)
	return err == nil
}

// statePath returns the file path for a state name.
func (m *Manager) statePath(name string) string {
	return filepath.Join(m.dir, name+".json")
}

// getStateInfo retrieves info about a specific state.
func (m *Manager) getStateInfo(name string) (*StateInfo, error) {
	path := m.statePath(name)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Name      string                `json:"name"`
		CreatedAt time.Time             `json:"created_at"`
		State     *w3pilot.StorageState `json:"state"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	info := &StateInfo{
		Name:      name,
		CreatedAt: wrapper.CreatedAt,
		Size:      fileInfo.Size(),
	}

	if wrapper.State != nil {
		info.NumCookies = len(wrapper.State.Cookies)
		for _, origin := range wrapper.State.Origins {
			info.Origins = append(info.Origins, origin.Origin)
		}
	}

	return info, nil
}

func isValidNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}
