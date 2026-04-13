package ws

import (
	"encoding/json"
	"mini-exchange/app/entity"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

func NewWSHandler(h Hub) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		client := &Client{
			Conn: c,
			Send: make(chan []byte, 256),
		}

		h.Register(client)

		// Writer
		go func() {
			defer func() {
				h.Unregister(client)
				c.Close()
			}()
			for {
				msg, ok := <-client.Send
				if !ok {
					return
				}
				if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Error().Err(err).Msg("Websocket write error")
					return
				}
			}
		}()

		// Reader
		defer func() {
			h.Unregister(client)
			c.Close()
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error().Err(err).Msg("Websocket read error")
				}
				break
			}

			var req entity.WSMessage
			if err := json.Unmarshal(message, &req); err != nil {
				log.Warn().Err(err).Msg("Invalid WS message format")
				continue
			}

			switch req.Type {
			case "subscribe":
				h.Subscribe(client, req.Channel)
			case "unsubscribe":
				h.Unsubscribe(client, req.Channel)
			default:
				log.Warn().Str("type", req.Type).Msg("Unknown WS message type")
			}
		}
	}
}
