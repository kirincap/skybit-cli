package ws

import (
	"context"
	"encoding/json"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	conn *websocket.Conn
}

func New(url string) (*Client, error) {
	c := &Client{}
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, nil
}

func (c *Client) Send(v any) error {
	return wsjson.Write(context.Background(), c.conn, v)
}

func (c *Client) Read(handler func(json.RawMessage) error) error {
	for {
		var msg json.RawMessage
		if err := wsjson.Read(context.Background(), c.conn, &msg); err != nil {
			return err
		}
		if err := handler(msg); err != nil {
			return err
		}
	}
}
