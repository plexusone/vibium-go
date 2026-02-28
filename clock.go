package vibium

import (
	"context"
	"time"
)

// Clock provides control over time in the browser.
type Clock struct {
	client  *BiDiClient
	context string
}

// NewClock creates a new Clock controller.
func NewClock(client *BiDiClient, browsingContext string) *Clock {
	return &Clock{
		client:  client,
		context: browsingContext,
	}
}

// ClockInstallOptions configures clock installation.
type ClockInstallOptions struct {
	// Time to set as the current time (Unix timestamp in milliseconds or time.Time)
	Time interface{}
}

// Install installs fake timers in the browser.
// This replaces native time-related functions like Date, setTimeout, etc.
func (c *Clock) Install(ctx context.Context, opts *ClockInstallOptions) error {
	params := map[string]interface{}{
		"context": c.context,
	}

	if opts != nil && opts.Time != nil {
		switch t := opts.Time.(type) {
		case time.Time:
			params["time"] = t.UnixMilli()
		case int64:
			params["time"] = t
		}
	}

	_, err := c.client.Send(ctx, "vibium:clock.install", params)
	return err
}

// FastForward advances time by the specified number of milliseconds.
// Timers are not fired.
func (c *Clock) FastForward(ctx context.Context, ticks int64) error {
	params := map[string]interface{}{
		"context": c.context,
		"ticks":   ticks,
	}

	_, err := c.client.Send(ctx, "vibium:clock.fastForward", params)
	return err
}

// RunFor advances time by the specified number of milliseconds,
// firing all pending timers.
func (c *Clock) RunFor(ctx context.Context, ticks int64) error {
	params := map[string]interface{}{
		"context": c.context,
		"ticks":   ticks,
	}

	_, err := c.client.Send(ctx, "vibium:clock.runFor", params)
	return err
}

// PauseAt pauses time at the specified timestamp.
func (c *Clock) PauseAt(ctx context.Context, t interface{}) error {
	params := map[string]interface{}{
		"context": c.context,
	}

	switch tv := t.(type) {
	case time.Time:
		params["time"] = tv.UnixMilli()
	case int64:
		params["time"] = tv
	}

	_, err := c.client.Send(ctx, "vibium:clock.pauseAt", params)
	return err
}

// Resume resumes time from a paused state.
func (c *Clock) Resume(ctx context.Context) error {
	params := map[string]interface{}{
		"context": c.context,
	}

	_, err := c.client.Send(ctx, "vibium:clock.resume", params)
	return err
}

// SetFixedTime sets a fixed time that will be returned by Date.now() and new Date().
func (c *Clock) SetFixedTime(ctx context.Context, t interface{}) error {
	params := map[string]interface{}{
		"context": c.context,
	}

	switch tv := t.(type) {
	case time.Time:
		params["time"] = tv.UnixMilli()
	case int64:
		params["time"] = tv
	}

	_, err := c.client.Send(ctx, "vibium:clock.setFixedTime", params)
	return err
}

// SetSystemTime sets the system time.
func (c *Clock) SetSystemTime(ctx context.Context, t interface{}) error {
	params := map[string]interface{}{
		"context": c.context,
	}

	switch tv := t.(type) {
	case time.Time:
		params["time"] = tv.UnixMilli()
	case int64:
		params["time"] = tv
	}

	_, err := c.client.Send(ctx, "vibium:clock.setSystemTime", params)
	return err
}

// SetTimezone sets the timezone.
func (c *Clock) SetTimezone(ctx context.Context, tz string) error {
	params := map[string]interface{}{
		"context":  c.context,
		"timezone": tz,
	}

	_, err := c.client.Send(ctx, "vibium:clock.setTimezone", params)
	return err
}
