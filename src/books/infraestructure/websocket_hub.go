package infraestructure

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage estructura para mensajes en WebSocket (tareas)
// Type: "create" | "update" | "delete" | "subscribe"
type WebSocketMessage struct {
	Type   string      `json:"type"`
	TaskID int32       `json:"task_id"`
	Data   interface{} `json:"data"`
}

// Hub gestiona todas las conexiones WebSocket del módulo tareas
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WebSocketMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client representa una conexión de cliente
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan WebSocketMessage
	taskID int32
	userID int32
}

var hub *Hub

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan WebSocketMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.run()
}

func GetHub() *Hub {
	return hub
}

// run ejecuta el bucle principal del hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[Tasks WS] Cliente registrado. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("[Tasks WS] Cliente desregistrado. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// readPump lee mensajes del cliente
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(1440 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(1440 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[Tasks WS] Error: %v", err)
			}
			return
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[Tasks WS] Error al parsear mensaje: %v", err)
			continue
		}

		c.hub.broadcast <- msg
	}
}

// writePump escribe mensajes al cliente
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if err := json.NewEncoder(w).Encode(message); err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// BroadcastEvent envía un evento a todos los clientes
func (h *Hub) BroadcastEvent(eventType string, taskID int32, data interface{}) {
	message := WebSocketMessage{
		Type:   eventType,
		TaskID: taskID,
		Data:   data,
	}
	select {
	case h.broadcast <- message:
	default:
		log.Println("[Tasks WS] Broadcast canal lleno")
	}
}
