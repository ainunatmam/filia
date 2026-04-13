package ws

import (
	"encoding/json"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
	"strings"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub interface {
	Run()
	Register(client *Client)
	Unregister(client *Client)
	Subscribe(client *Client, channel string)
	Unsubscribe(client *Client, channel string)
	Broadcast(channel string, data interface{})
}

type hub struct {
	clients    map[*Client]bool
	channels   map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMessage
	mu         sync.RWMutex
}

type broadcastMessage struct {
	channel string
	data    interface{}
}

func NewHub() Hub {
	return &hub{
		clients:    make(map[*Client]bool),
		channels:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan broadcastMessage, 1000),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			libraries.Logger.Info().Msg("WebSocket: Client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				for channel := range h.channels {
					delete(h.channels[channel], client)
				}
				close(client.Send)
			}
			h.mu.Unlock()
			libraries.Logger.Info().Msg("WebSocket: Client unregistered")

		case msg := <-h.broadcast:
			h.mu.RLock()
			// Normalize channel name to UpperCase and trim space
			channel := strings.ToUpper(strings.TrimSpace(msg.channel))
			subscribers, exists := h.channels[channel]
			
			if !exists || len(subscribers) == 0 {
				h.mu.RUnlock()
				continue
			}

			payload := entity.WSMessage{
				Type:    "broadcast",
				Channel: channel,
				Data:    msg.data,
			}
			b, err := json.Marshal(payload)
			if err != nil {
				libraries.Logger.Error().Err(err).Msg("WebSocket: Marshal error")
				h.mu.RUnlock()
				continue
			}

			libraries.Logger.Debug().Str("channel", channel).Int("subscribers", len(subscribers)).Msg("WebSocket: Broadcasting message")
			
			for client := range subscribers {
				select {
				case client.Send <- b:
				default:
					libraries.Logger.Warn().Msg("WebSocket: Client send buffer full, dropping message")
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *hub) Register(client *Client) {
	h.register <- client
}

func (h *hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *hub) Subscribe(client *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Normalize channel name to UpperCase and trim space
	channel = strings.ToUpper(strings.TrimSpace(channel))
	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*Client]bool)
	}
	h.channels[channel][client] = true
	libraries.Logger.Info().Str("channel", channel).Msg("WebSocket: Client subscribed")
}

func (h *hub) Unsubscribe(client *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Normalize channel name to UpperCase and trim space
	channel = strings.ToUpper(strings.TrimSpace(channel))
	if _, ok := h.channels[channel]; ok {
		delete(h.channels[channel], client)
	}
	libraries.Logger.Info().Str("channel", channel).Msg("WebSocket: Client unsubscribed")
}

func (h *hub) Broadcast(channel string, data interface{}) {
	h.broadcast <- broadcastMessage{
		channel: channel,
		data:    data,
	}
}
