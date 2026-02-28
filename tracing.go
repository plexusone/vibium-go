package vibium

import (
	"context"
	"encoding/base64"
	"encoding/json"
)

// Tracing provides control over trace recording.
type Tracing struct {
	client      *BiDiClient
	userContext string
}

// TracingStartOptions configures trace recording.
type TracingStartOptions struct {
	// Name is the trace file name (optional).
	Name string

	// Screenshots includes screenshots in the trace.
	Screenshots bool

	// Snapshots includes DOM snapshots in the trace.
	Snapshots bool

	// Sources includes source files in the trace.
	Sources bool

	// Title is the trace title (shown in trace viewer).
	Title string

	// Categories specifies which trace categories to include.
	Categories []string
}

// TracingStopOptions configures how to stop trace recording.
type TracingStopOptions struct {
	// Path to save the trace file.
	Path string
}

// TracingChunkOptions configures trace chunk recording.
type TracingChunkOptions struct {
	// Name for the chunk.
	Name string

	// Title for the chunk.
	Title string
}

// TracingGroupOptions configures trace groups.
type TracingGroupOptions struct {
	// Location to associate with the group.
	Location string
}

// Start starts trace recording.
func (t *Tracing) Start(ctx context.Context, opts *TracingStartOptions) error {
	params := map[string]interface{}{
		"userContext": t.userContext,
	}

	if opts != nil {
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Screenshots {
			params["screenshots"] = opts.Screenshots
		}
		if opts.Snapshots {
			params["snapshots"] = opts.Snapshots
		}
		if opts.Sources {
			params["sources"] = opts.Sources
		}
		if opts.Title != "" {
			params["title"] = opts.Title
		}
		if len(opts.Categories) > 0 {
			params["categories"] = opts.Categories
		}
	}

	_, err := t.client.Send(ctx, "vibium:tracing.start", params)
	return err
}

// Stop stops trace recording and returns the trace data.
func (t *Tracing) Stop(ctx context.Context, opts *TracingStopOptions) ([]byte, error) {
	params := map[string]interface{}{
		"userContext": t.userContext,
	}

	if opts != nil && opts.Path != "" {
		params["path"] = opts.Path
	}

	result, err := t.client.Send(ctx, "vibium:tracing.stop", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Data)
}

// StartChunk starts a new trace chunk.
func (t *Tracing) StartChunk(ctx context.Context, opts *TracingChunkOptions) error {
	params := map[string]interface{}{
		"userContext": t.userContext,
	}

	if opts != nil {
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Title != "" {
			params["title"] = opts.Title
		}
	}

	_, err := t.client.Send(ctx, "vibium:tracing.startChunk", params)
	return err
}

// StopChunk stops the current trace chunk and returns the data.
func (t *Tracing) StopChunk(ctx context.Context, opts *TracingChunkOptions) ([]byte, error) {
	params := map[string]interface{}{
		"userContext": t.userContext,
	}

	if opts != nil {
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Title != "" {
			params["title"] = opts.Title
		}
	}

	result, err := t.client.Send(ctx, "vibium:tracing.stopChunk", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Data)
}

// StartGroup starts a new trace group.
func (t *Tracing) StartGroup(ctx context.Context, name string, opts *TracingGroupOptions) error {
	params := map[string]interface{}{
		"userContext": t.userContext,
		"name":        name,
	}

	if opts != nil && opts.Location != "" {
		params["location"] = opts.Location
	}

	_, err := t.client.Send(ctx, "vibium:tracing.startGroup", params)
	return err
}

// StopGroup stops the current trace group.
func (t *Tracing) StopGroup(ctx context.Context) error {
	params := map[string]interface{}{
		"userContext": t.userContext,
	}

	_, err := t.client.Send(ctx, "vibium:tracing.stopGroup", params)
	return err
}
