package conn

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"os"
	"reflect"
)

type eventHandler struct {
	fn   reflect.Value
	args []reflect.Type
}

type Server struct {
	clients map[int]*Client
	events  map[string]*eventHandler
}

func NewServer() *Server {
	clients := make(map[int]*Client)
	events := make(map[string]*eventHandler)
	return &Server{
		clients,
		events,
	}
}

func genEventHandler(fn interface{}) (handler *eventHandler, err error) {
	fnValue := reflect.ValueOf(fn)
	handler = new(eventHandler)
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = fmt.Errorf("%v is not a function", fn)
		return
	}
	handler.fn = fnValue
	fnType := fnValue.Type()
	nArgs := fnValue.Type().NumIn()
	handler.args = make([]reflect.Type, nArgs)
	if nArgs == 0 {
		err = errors.New("no arg exists")
		return
	}
	if t := fnType.In(0); t.Kind() != reflect.Ptr || t.Elem().Name() != "Client" {
		err = errors.New("first argument should be of type *Client")
		return
	} else {
		handler.args[0] = t
	}
	for i := 1; i < nArgs; i++ {
		handler.args[i] = fnType.In(i)
	}
	return
}

func (self *Server) Handle(event_name string, fn interface{}) error {
	handler, err := genEventHandler(fn)
	if err != nil {
		return err
	}
	self.events[event_name] = handler
	return nil
}

func (self *Server) Call(message *Message, client *Client) *Message {
	handler := self.events[message.EventName]
	if handler != nil {
		callArgs := make([]reflect.Value, len(message.EventData)+1)
		callArgs[0] = reflect.ValueOf(client)
		for i, arg := range message.EventData {
			callArgs[i+1] = reflect.ValueOf(arg)
		}
		return safeCall(message.EventName, handler.fn, callArgs)
	} else {
		return &Message{
			EventName: message.EventName + " response",
			EventData: []interface{}{},
		}
	}
}

func safeCall(event_name string, fn reflect.Value, args []reflect.Value) *Message {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	ret := fn.Call(args)

	response := &Message{
		EventName: event_name + " response",
		EventData: []interface{}{},
	}
	if len(ret) > 0 {
		retArgs := make([]interface{}, len(ret))
		for i, arg := range ret {
			retArgs[i] = arg.Interface()
		}
		response.EventData = retArgs
	}

	return response
}

func (self *Server) Broadcast(event_name string, args ...interface{}) {
	for _, client := range self.clients {
		client.Emit(event_name, args...)
	}
}

func (self *Server) KillClient(client *Client) {
	self.Broadcast("kill", client.Id)
	self.Call(&Message{"disconnected", []interface{}{}}, client)
	delete(self.clients, client.Id)
}

func (self *Server) Listen(options map[string]string) {
	if pattern := options["pattern"]; pattern != "" {
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
			if _, ok := err.(websocket.HandshakeError); ok {
				http.Error(w, "Not a websocket handshake", 400)
				return
			}
			client := newWebSocketClient(self, ws)
			self.clients[client.Id] = client
			self.Call(&Message{"connected", []interface{}{}}, client)
			client.Emit("connected", client)
			go client.ListenWebSocket()
		})
	}

	if host := options["host"]; host != "" {
		l, err := net.Listen("tcp", host)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}
		// Close the listener when the application closes.
		defer l.Close()
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
			}
			// Handle connections in a new goroutine.
			client := newTCPClient(self, conn)
			self.clients[client.Id] = client
			go client.ListenTCP()
		}
	}
}
