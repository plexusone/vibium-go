package vibium

import (
	"context"
	"encoding/json"
)

// Download represents a file download.
type Download struct {
	client  *BiDiClient
	context string
	id      string
	URL     string `json:"url"`
	Name    string `json:"suggestedFilename"`
}

// Path returns the path to the downloaded file after it completes.
func (d *Download) Path(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
	}

	result, err := d.client.Send(ctx, "vibium:download.path", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Path, nil
}

// SaveAs saves the download to the specified path.
func (d *Download) SaveAs(ctx context.Context, path string) error {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
		"path":    path,
	}

	_, err := d.client.Send(ctx, "vibium:download.saveAs", params)
	return err
}

// Cancel cancels the download.
func (d *Download) Cancel(ctx context.Context) error {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
	}

	_, err := d.client.Send(ctx, "vibium:download.cancel", params)
	return err
}

// Failure returns the download failure reason, if any.
func (d *Download) Failure(ctx context.Context) (string, error) {
	params := map[string]interface{}{
		"context": d.context,
		"id":      d.id,
	}

	result, err := d.client.Send(ctx, "vibium:download.failure", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Failure string `json:"failure"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", err
	}

	return resp.Failure, nil
}
