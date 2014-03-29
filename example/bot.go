package main

import (
	"botbattle/client"
)

type Bot struct {
  *client.BotClient
}

func (self *Bot) ScanArena(){
  first_status := self.Status()
  last_x := first_status.X
  last_y := first_status.Y
  self.MoveForward()
  for {
    if status := self.Status(); last_x != status.X || last_y != status.Y {
      self.ScanAndShoot()
      self.ScanAndShoot()
      self.ScanAndShoot()
      self.ScanAndShoot()
      self.MoveForward()
      last_x = status.X
      last_y = status.Y
    } else {
      break
    }
  }
}

func (self *Bot) ScanAndShoot(){
  bots := self.Scan()
  if len(bots) > 0 {
    self.FireCannon()
  }
  self.RotRight()
}


func (self *Bot) Chase(){
}

func main() {
  bot := &Bot{client.NewBotClient("localhost:3333", "Tim")}
  for {
    bot.ScanArena()
    bot.RotRight()
  }
}
