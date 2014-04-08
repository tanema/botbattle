package client

import (
	"botbattle/conn"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type BotClient struct {
	conn net.Conn
  ArenaHeight int
  ArenaWidth  int
}

type Status struct {
  Id           int
	X            int
	Y            int
	Rotation     int
	Health       int
}

func NewBotClient(host, name string) *BotClient {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		println("could not connect to that host")
		fmt.Println(err)
		os.Exit(0)
	}
  newClient := new(BotClient)
  newClient.conn = conn
	resp := newClient.request("register", name)
	newClient.ArenaWidth = int(resp.EventData[0].(float64))
	newClient.ArenaHeight = int(resp.EventData[1].(float64))
	return newClient
}

func (self *BotClient) Status() *Status {
	resp := self.request("status")
	id := int(resp.EventData[0].(float64))
	x := int(resp.EventData[1].(float64))
	y := int(resp.EventData[2].(float64))
	rotation := int(resp.EventData[3].(float64))
	health := int(resp.EventData[4].(float64))
	return &Status{id, x, y, rotation, health}
}

func (self *BotClient) RotLeft() (rot int) {
	resp := self.request("rotate left")
	rot = int(resp.EventData[0].(float64))
	return
}

func (self *BotClient) RotRight() (rot int) {
	resp := self.request("rotate right")
	rot = int(resp.EventData[0].(float64))
	return
}

func (self *BotClient) MoveForward() (x int, y int) {
	resp := self.request("move forward")
	x = int(resp.EventData[0].(float64))
	y = int(resp.EventData[1].(float64))
	return
}

func (self *BotClient) MoveBackward() (x int, y int) {
	resp := self.request("move backward")
	x = int(resp.EventData[0].(float64))
	y = int(resp.EventData[1].(float64))
	return
}

func (self *BotClient) Scan() []*Status {
	resp := self.request("scan")
	result := []*Status{}
	if statuses, ok := resp.EventData[0].([]interface{}); ok {
		for _, state_interface := range statuses {
			state_array := state_interface.([]interface{})
			new_status := new(Status)
			new_status.Id = int(state_array[0].(float64))
			new_status.X = int(state_array[1].(float64))
			new_status.Y = int(state_array[2].(float64))
			new_status.Rotation = int(state_array[3].(float64))
			new_status.Health = int(state_array[4].(float64))
			result = append(result, new_status)
		}
	}
	return result
}

func (self *BotClient) FireGun() bool {
	resp := self.request("fire gun")
	return resp.EventData[0].(bool)
}

func (self *BotClient) FireCannon() bool {
	resp := self.request("fire cannon")
	return resp.EventData[0].(bool)
}

func (self *BotClient) Shield() bool {
	resp := self.request("shield")
	return resp.EventData[0].(bool)
}

func (self *BotClient) request(line string, args ...interface{}) *conn.Message {
	message := conn.Message{
		EventName: line,
		EventData: args,
	}
	message_json, _ := json.Marshal(message)
	_, err := self.conn.Write(append(message_json, "\n"...))
	if err != nil {
		println("looks like you died")
		os.Exit(1)
	}

	reply, err := bufio.NewReader(self.conn).ReadString('\n')
	if err != nil {
		println("looks like you died")
		os.Exit(1)
	}

	response := &conn.Message{}
	json.Unmarshal([]byte(reply), response)

	return response
}
