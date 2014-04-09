package conn

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net"
)

var current_id = 0

type Client struct {
	Id     int             `json:"id"`
	server *Server         `json:"-"`
	tcp    net.Conn        `json:"-"`
	socket *websocket.Conn `json:"-"`
}

func newTCPClient(server *Server, c net.Conn) *Client {
	current_id++
	return &Client{current_id, server, c, nil}
}

func newWebSocketClient(server *Server, c *websocket.Conn) *Client {
	current_id++
	return &Client{current_id, server, nil, c}
}

func (self *Client) Close() {
  if self.tcp != nil {
    self.tcp.Close()
  }
  if self.socket != nil {
    self.socket.Close()
  }
}

func (self *Client) ListenTCP() {
	// Make a buffer to hold incoming data.
	reader := bufio.NewReader(self.tcp)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
      self.tcp.Close()
      self.server.KillClient(self)
      return
		}

		message := &Message{}
		json.Unmarshal([]byte(line), message)
		resp := self.server.Call(message, self)
		message_json, _ := json.Marshal(resp)

		self.tcp.Write(append(message_json, "\n"...))
	}
}

func (self *Client) ListenWebSocket() {
	for {
		message := &Message{}
		if err := self.socket.ReadJSON(message); err != nil {
      self.socket.Close()
			self.server.KillClient(self)
			return
		}
		self.socket.WriteJSON(self.server.Call(message, self))
	}
}

func (self *Client) Emit(event_name string, args ...interface{}) {
	message := &Message{
		EventName: event_name,
		EventData: args,
	}
	if self.socket != nil {
		self.socket.WriteJSON(message)
	}
}
