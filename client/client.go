package client

import (
	"github.com/tanema/botbattle/conn"
	"bufio"
	"encoding/json"
	"net"
)

type BotClient struct {
	conn net.Conn
  host string
  name string
  ArenaHeight int
  ArenaWidth  int
}

//duplicated here so the client user does not have to import game
type Status struct {
  Id        int     `json:"id"`
  X         int     `json:"x"`
  Y         int     `json:"y"`
  Rotation  int     `json:"rotation"`
  Name      string  `json:"name"`
  Health    int     `json:"health"`
  KillCount int     `json:"kill_count"`
}

func NewBotClient(host, name string) (*BotClient, error) {
  newClient := &BotClient{host: host, name: name}
	return newClient.Reconnect()
}

func (self *BotClient) Reconnect() (*BotClient, error) {
	conn, err := net.Dial("tcp", self.host)
	if err != nil {
    return nil, err
	}
  self.conn = conn
	resp, err := self.request("register", self.name)
  if err != nil {
    return nil, err
  }
	self.ArenaWidth = int(resp.EventData[0].(float64))
	self.ArenaHeight = int(resp.EventData[1].(float64))
  return self, nil
}

func (self *BotClient) Status() (*Status, error) {
	resp, err := self.request("status")
  if err != nil {
    return nil, err
  }
  my_status := &Status{}
	json.Unmarshal([]byte(resp.EventData[0].(string)), my_status)
	return my_status, nil
}

func (self *BotClient) RotLeft() (*Status, error) {
	resp, err := self.request("rotate left")
  if err != nil {
    return nil, err
  }
  my_status := &Status{}
	json.Unmarshal([]byte(resp.EventData[0].(string)), my_status)
	return my_status, nil
}

func (self *BotClient) RotRight() (*Status, error) {
	resp, err := self.request("rotate right")
  if err != nil {
    return nil, err
  }
  my_status := &Status{}
	json.Unmarshal([]byte(resp.EventData[0].(string)), my_status)
	return my_status, nil
}

func (self *BotClient) MoveForward() (*Status, error) {
	resp, err := self.request("move forward")
  if err != nil {
    return nil, err
  }
  my_status := &Status{}
	json.Unmarshal([]byte(resp.EventData[0].(string)), my_status)
	return my_status, nil
}

func (self *BotClient) MoveBackward() (*Status, error) {
	resp, err := self.request("move backward")
  if err != nil {
    return nil, err
  }
  my_status := &Status{}
	json.Unmarshal([]byte(resp.EventData[0].(string)), my_status)
	return my_status, nil
}

func (self *BotClient) Scan() ([]*Status, error) {
	resp, err := self.request("scan")
  if err != nil {
    return nil, err
  }
	result := []*Status{}
	if statuses, ok := resp.EventData[0].([]interface{}); ok {
		for _, state_interface := range statuses {
      bot_status := &Status{}
      json.Unmarshal([]byte(state_interface.(string)), bot_status)
			result = append(result, bot_status)
		}
	}
	return result, nil
}

func (self *BotClient) FireGun() (bool, error) {
	resp, err := self.request("fire gun")
  if err != nil {
    return false, err
  }
	return resp.EventData[0].(bool), nil
}

func (self *BotClient) FireCannon() (bool, error) {
	resp, err := self.request("fire cannon")
  if err != nil {
    return false, err
  }
	return resp.EventData[0].(bool), nil
}

func (self *BotClient) Shield() (bool, error) {
	resp, err := self.request("shield")
  if err != nil {
    return false, err
  }
	return resp.EventData[0].(bool), nil
}

func (self *BotClient) request(line string, args ...interface{}) (*conn.Message, error) {
	message := conn.Message{
		EventName: line,
		EventData: args,
	}
	message_json, err := json.Marshal(message)
  if err != nil {
    return nil, err
  }
	_, err = self.conn.Write(append(message_json, "\n"...))
	if err != nil {
		self.conn.Close()
    self.conn = nil
    return nil, err
	}

	reply, reply_err := bufio.NewReader(self.conn).ReadString('\n')
	if reply_err != nil {
		self.conn.Close()
    self.conn = nil
    return nil, reply_err
	}

	response := &conn.Message{}
	json.Unmarshal([]byte(reply), response)

	return response, nil
}
