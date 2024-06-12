// Package client provides a client for listening to a Streamr node.
// It does not yet support sending messages to the node.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"github.com/kwilteam/kwil-db/core/log"
)

// Client is a Streamr client. It is meant to connect to a Streamr websocket server.
// One client should be used for each stream subscription.
// It is a thin wrapper around the gorilla/websocket.Conn type that handles connection
// to the Streamr node.
type Client struct {
	conn   *websocket.Conn
	mu     sync.Mutex // mu protects all methods.
	config *ClientConfig
	url    string
}

// NewClient creates a new Streamr client.
// streamrWebsocketUrl is the URL of the Streamr node's websocket server.
// streamrWebsocketUrl should be in the form of "ws://<host>:<port>"/"wss://<host>:<port>".
// streamID is the ID of the stream to subscribe to.
// Opts can be nil, in which case the client will use the default configuration.
func NewClient(ctx context.Context, streamrWebsocketUrl, streamID string, opts *ClientConfig) (*Client, error) {
	conf := DefaultConfig()
	if opts != nil {
		conf.Apply(opts)
	}

	path := "/streams/" + url.PathEscape(streamID) + "/subscribe"
	if opts.ApiKey != nil {
		path += "?apiKey=" + *opts.ApiKey
	}
	fullUrl := streamrWebsocketUrl + path

	conn, req, err := websocket.DefaultDialer.DialContext(ctx, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	return &Client{
		conn:   conn,
		config: opts,
		url:    fullUrl,
	}, nil
}

// Close closes the client's connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

// ReadMessage reads a message from the client's connection.
func (c *Client) ReadMessage() (ev *StreamrEvent, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.readMessage()
}

// readMessage is an internal method that reads a message from the client's connection.
// It does not handle locking.
func (c *Client) readMessage() (ev *StreamrEvent, err error) {
	_, p, err := c.conn.ReadMessage()
	// if error, we need to close and retry
	if err != nil {
		c.config.Logger.Info("failed to read message from Streamr node, attempting to reconnect")
		b := &backoff.Backoff{
			Min:    *c.config.MinRetryDelay,
			Max:    *c.config.MaxRetryDelay,
			Factor: 2,
			Jitter: true,
		}

		for i := 0; i < *c.config.MaxRetrys; i++ {
			time.Sleep(b.Duration())
			c.conn.Close()
			c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
			if err == nil {
				c.config.Logger.Info("reconnected to Streamr node, retrying readMessage")
				// if reconnection is successful, try to read again
				return c.readMessage()
			}
			c.config.Logger.Info("failed to reconnect to Streamr node, attempt %d", i)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to reconnect to Streamr node after %d attempts: %v", *c.config.MaxRetrys, err)
		}
	}

	ev = &StreamrEvent{}
	err = json.Unmarshal(p, ev)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %v", err)
	}

	return ev, nil
}

// ClientConfig is a configuration struct for the Client.
// Pointers are used to differentiate between zero values and non-zero values.
type ClientConfig struct {
	// ApiKey is the API key to use for the connection.
	ApiKey *string
	// MaxRetrys is the maximum number of times to retry the connection on failure.
	// Default is 3.
	MaxRetrys *int
	// MinRetryDelay is the minimum delay between retries.
	// Default is 1 second.
	MinRetryDelay *time.Duration
	// MaxRetryDelay is the maximum delay between retries.
	// Default is 10 seconds.
	MaxRetryDelay *time.Duration
	// Logger is the logger to use for the client.
	Logger *log.SugaredLogger
}

// Apply applies non-default values from the given configuration to the config.
func (c *ClientConfig) Apply(config *ClientConfig) {
	if config.ApiKey != nil {
		c.ApiKey = config.ApiKey
	}
	if config.MaxRetrys != nil {
		c.MaxRetrys = config.MaxRetrys
	}
	if config.MinRetryDelay != nil {
		c.MinRetryDelay = config.MinRetryDelay
	}
	if config.MaxRetryDelay != nil {
		c.MaxRetryDelay = config.MaxRetryDelay
	}
	if config.Logger != nil {
		c.Logger = config.Logger
	}

}

// DefaultConfig returns the default configuration for the client.
func DefaultConfig() *ClientConfig {
	r := 3
	min := time.Second
	max := 10 * time.Second
	l := log.NewNoOp().Sugar()
	return &ClientConfig{
		MaxRetrys:     &r,
		MinRetryDelay: &min,
		MaxRetryDelay: &max,
		Logger:        &l,
	}
}

// StreamrEvent is an event read from a Streamr node.
type StreamrEvent struct {
	// Content is the user-determined content of the event.
	// This is arbitrary and can be anything.
	Content any `json:"content"`
	// Metadata is the metadata of the event, provided
	// by the Streamr network.
	Metadata struct {
		Timestamp      int64  `json:"timestamp"`
		SequenceNumber int64  `json:"sequenceNumber"`
		PublisherID    string `json:"publisherId"`
		MsgChainID     string `json:"msgChainId"`
	} `json:"metadata"`
}
