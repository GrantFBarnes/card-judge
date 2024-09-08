package websocket

import (
	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/database"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Lobby ID
	lobbyId uuid.UUID

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(lobbyId uuid.UUID) *Hub {
	return &Hub{
		lobbyId:    lobbyId,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
			if len(h.clients) == 0 {
				database.DeleteLobby(h.lobbyId)
				delete(lobbyHubs, h.lobbyId)
				return
			}
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
	database.AddPlayerToLobby(h.lobbyId, client.player.Id)
	h.broadcastMessage([]byte("Player " + client.player.Name + " Joined"))
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		database.RemovePlayerFromLobby(h.lobbyId, client.player.Id)
	}
	h.broadcastMessage([]byte("Player " + client.player.Name + " Left"))
}

func (h *Hub) broadcastMessage(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}