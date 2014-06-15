package game

import (
	"github.com/tanema/botbattle/conn"
	"encoding/json"
	"time"
)

const (
	ARENA_HEIGHT = 11
	ARENA_WIDTH  = 23
)

type Scene struct {
	serv *conn.Server
	bots map[int]*Bot
}

func NewScene() *Scene {
	server := conn.NewServer()
	newScene := &Scene{
		server,
		make(map[int]*Bot),
	}
	newScene.bindActions()
	return newScene
}

func (self *Scene) bindActions() {
	self.serv.Handle("connected", self.onWebSocketConnected)
	self.serv.Handle("register", self.onRegister)
	self.serv.Handle("status", self.onStatus)
	self.serv.Handle("disconnected", self.onBotDisconnect)
	self.serv.Handle("rotate left", self.onBotRotLeft)
	self.serv.Handle("rotate right", self.onBotRotRight)
	self.serv.Handle("move forward", self.onBotMoveForward)
	self.serv.Handle("move backward", self.onBotMoveBackward)
	self.serv.Handle("fire gun", self.onFireGun)
	self.serv.Handle("fire cannon", self.onFireCannon)
	self.serv.Handle("scan", self.onScan)
	self.serv.Handle("shield", self.onShield)
}


func (self *Scene) onWebSocketConnected(client *conn.Client) {
  result := []Status{}
  for _, bot := range self.bots {
    result = append(result, Status{bot.client.Id, bot.x, bot.y, bot.rotation, bot.name, bot.health, bot.killcount})
  }
  go func(){
	  time.Sleep(1000 * time.Millisecond)
    client.Emit("connected", result, ARENA_WIDTH, ARENA_HEIGHT)
  }()
  return
}

func (self *Scene) onRegister(client *conn.Client, name string) (int, int) {
	newBot := NewBot(self, client, name)
	self.bots[client.Id] = newBot
	self.serv.Broadcast("register", newBot.client.Id, newBot.x, newBot.y, newBot.rotation, name)
	return ARENA_WIDTH, ARENA_HEIGHT
}

func (self *Scene) onStatus(client *conn.Client) string {
	if bot := self.bots[client.Id]; bot != nil {
	  json_resp, _ := json.Marshal(bot.Status())
	  return string(json_resp)
	}

	json_resp, _ := json.Marshal(Status{})
	return string(json_resp)
}

func (self *Scene) onBotDisconnect(client *conn.Client) {
	if bot := self.bots[client.Id]; bot != nil {
		self.KillBot(bot)
	}
}

func (self *Scene) KillBot(bot *Bot) {
  bot.client.Close()
	delete(self.bots, bot.client.Id)
}

func (self *Scene) onBotRotLeft(client *conn.Client) string {
	if bot := self.bots[client.Id]; bot != nil {
    status := bot.RotLeft()
		self.serv.Broadcast("rotate", bot.client.Id, status.Rotation)
	  json_resp, _ := json.Marshal(status)
	  return string(json_resp)
	}
  return ""
}

func (self *Scene) onBotRotRight(client *conn.Client) string {
	if bot := self.bots[client.Id]; bot != nil {
    status := bot.RotRight()
		self.serv.Broadcast("rotate", bot.client.Id, status.Rotation)
	  json_resp, _ := json.Marshal(status)
	  return string(json_resp)
	}
  return ""
}

func (self *Scene) onBotMoveForward(client *conn.Client) string {
	if bot := self.bots[client.Id]; bot != nil {
    status := bot.MoveForward()
		self.serv.Broadcast("move", bot.client.Id, status.X, status.Y)
    json_resp, _ := json.Marshal(status)
    return string(json_resp)
	}
  return ""
}

func (self *Scene) onBotMoveBackward(client *conn.Client) string {
	if bot := self.bots[client.Id]; bot != nil {
    status := bot.MoveBackward()
		self.serv.Broadcast("move", bot.client.Id, status.X, status.Y)
    json_resp, _ := json.Marshal(status)
    return string(json_resp)
	}
  return ""
}

func (self *Scene) onFireGun(client *conn.Client) bool {
	if bot := self.bots[client.Id]; bot != nil {
		self.serv.Broadcast("fire gun", bot.client.Id)
		return bot.FireGun()
	}
	return false
}

func (self *Scene) onFireCannon(client *conn.Client) bool {
	if bot := self.bots[client.Id]; bot != nil {
		self.serv.Broadcast("fire cannon", bot.client.Id)
		return bot.FireCannon()
	}
	return false
}

func (self *Scene) onScan(client *conn.Client) []string {
	if bot := self.bots[client.Id]; bot != nil {
		self.serv.Broadcast("scan", bot.client.Id)
		return bot.Scan()
	}
	return []string{}
}

func (self *Scene) onShield(client *conn.Client) bool {
	if bot := self.bots[client.Id]; bot != nil {
		return bot.Shield()
	}
	return false
}

func (self *Scene) Start() {
	go self.serv.Listen(map[string]string{
		"host":    "0.0.0.0:3333",
		"pattern": "/ws",
	})
}
