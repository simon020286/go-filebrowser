package ws

import (
	"log"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

// Hub struct definition.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

// NewHub constructor.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// ServeWs connect new client.
func (h *Hub) ServeWs(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := &Client{hub: h, conn: conn, send: make(chan Message, 256)}
		client.hub.register <- client

		go client.writePump()
		client.readPump()
	})

	if err != nil {
		log.Println(err)
	}
}

// Run hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Broadcast message to all clients.
func (h *Hub) Broadcast(event string, payload interface{}) {
	h.broadcast <- NewMessage(event, payload)
}
