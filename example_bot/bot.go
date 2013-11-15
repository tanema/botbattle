package main

import (
  "github.com/tanema/botbattle/client"
  "math"
  "time"
)

type Bot struct {
  X, Y float64
  client client.Client
}

var last_found *Bot

func main() {
  bot := &Bot{client: client.NewClient("localhost:4569", "iTim")}
  bot.scan()
  for {}
}

func (b *Bot) scan(){
  b.client.Stop()
  last_found = nil
  b.client.GetCurrentPosition(func(x, y float64){
    b.X = x
    b.Y = y
  })
  for i := 1; i < 360; i++ {
    b.client.RotateTo(float32(i))
    go b.client.Scan(func(x, y float64, name string){
      if name != "" {
        last_found = &Bot{X: x, Y:y}
      }
      return
    })
    time.Sleep(time.Millisecond)
  }
  if last_found != nil {
    b.chase(last_found.X, last_found.Y)
  } else {
    b.scan()
  }
}

func (b *Bot) chase(x, y float64){
  deltaX := last_found.X - b.X
  deltaY := last_found.Y - b.Y
  deg := (math.Abs(math.Atan2(deltaY, deltaX)) * (180 / math.Pi))
  if b.Y < last_found.Y {
    b.client.RotateTo(float32(deg) - 90)
  } else {
    b.client.RotateTo(-float32(deg) - 90)
  }
  b.client.Forward()
  for i := 0; i < 10; i++ {
    b.client.Shoot()
    time.Sleep(100 * time.Millisecond)
  }
  b.scan()
}
