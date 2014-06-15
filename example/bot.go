package main

import (
	"github.com/tanema/botbattle/client"
  "fmt"
)

type Bot struct {
  *client.BotClient
}

func (self *Bot) ScanArena(){
  first_status, _ := self.Status()
  last_x := first_status.X
  last_y := first_status.Y
  for {
    if status, err := self.MoveForward(); err == nil && (last_x != status.X || last_y != status.Y) {
      self.ScanAndShoot()
      self.ScanAndShoot()
      self.ScanAndShoot()
      self.ScanAndShoot()
      last_x = status.X
      last_y = status.Y
    } else {
      break
    }
  }
}

func (self *Bot) ScanAndShoot(){
  bots, _ := self.Scan()
  if len(bots) > 0 {
    self.FireCannon()
  }
  self.RotRight()
}


func (self *Bot) Chase(){
}

func main() {
  botclient, err := client.NewBotClient("localhost:3333", "Tim") 
  if err != nil {
    fmt.Println(err)
    return
  }

  bot := &Bot{botclient}
  for {
    bot.Shield()
    bot.ScanArena()
    bot.RotRight()
  }
}
