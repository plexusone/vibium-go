package vibium

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// BiDiCommand represents a WebDriver BiDi command.
type BiDiCommand struct {
	ID     int64       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BiDiResponse represents a WebDriver BiDi response.
type BiDiResponse struct {
	ID     int64           `json:"id"`
	Type   string          `json:"type"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
	Message string         `json:"message,omitempty"`
}

// BiDiClient manages WebSocket communication with the clicker server.
type BiDiClient struct {
	conn      *websocket.Conn
	url       string
	nextID    atomic.Int64
	pending   map[int64]chan *BiDiResponse
	pendingMu sync.Mutex
	closed    atomic.Bool
	closeCh   chan struct{}
}

// NewBiDiClient creates a new BiDi client.
func NewBiDiClient() *BiDiClient {
	return &BiDiClient{
		pending: make(map[int64]chan *BiDiResponse),
		closeCh: make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection to the clicker server.
func (c *BiDiClient) Connect(ctx context.Context, url string) error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return &ConnectionError{URL: url, Cause: err}
	}

	c.conn = conn
	c.url = url

	// Start message receiver
	go c.receiveLoop()

	return nil
}

// Close closes the WebSocket connection.
func (c *BiDiClient) Close() error {
	if c.closed.Swap(true) {
		return nil // Already closed
	}

	close(c.closeCh)

	// Reject all pending requests
	c.pendingMu.Lock()
	for _, ch := range c.pending {
		close(ch)
	}
	c.pending = make(map[int64]chan *BiDiResponse)
	c.pendingMu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send sends a command and waits for the response.
func (c *BiDiClient) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	if c.closed.Load() {
		return nil, ErrConnectionClosed
	}

	id := c.nextID.Add(1)
	cmd := BiDiCommand{
		ID:     id,
		Method: method,
		Params: params,
	}

	// Create response channel
	respCh := make(chan *BiDiResponse, 1)
	c.pendingMu.Lock()
	c.pending[id] = respCh
	c.pendingMu.Unlock()

	// Clean up on exit
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	// Send command
	if err := c.conn.WriteJSON(cmd); err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.closeCh:
		return nil, ErrConnectionClosed
	case resp, ok := <-respCh:
		if !ok {
			return nil, ErrConnectionClosed
		}
		if resp.Type == "error" || resp.Error != "" {
			return nil, &BiDiError{
				ErrorType: resp.Error,
				Message:   resp.Message,
			}
		}
		return resp.Result, nil
	}
}

func (c *BiDiClient) receiveLoop() {
	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		var resp BiDiResponse
		if err := c.conn.ReadJSON(&resp); err != nil {
			if c.closed.Load() {
				return
			}
			// Connection error - close everything
			_ = c.Close()
			return
		}

		// Route response to waiting request
		c.pendingMu.Lock()
		ch, ok := c.pending[resp.ID]
		c.pendingMu.Unlock()

		if ok {
			select {
			case ch <- &resp:
			default:
			}
		}
	}
}
