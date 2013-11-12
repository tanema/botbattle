package main

import (
  "flag"
  "runtime"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
  "github.com/tanema/botbattle/scene"
  "github.com/tanema/botbattle/server"
)

func main() {
  var player *scene.BotController
	runtime.GOMAXPROCS(8)
	defer func() {
		engine.Terminate()
	}()
	engine.StartEngine()
	engine.LoadScene(scene.MainSceneGeneral)

  var host string
  flag.StringVar(&host, "h", "localhost:4569", "Your localhost that you are listening on")
  flag.Parse()

  game_server := server.NewServer(host)
  go game_server.Listen()
	for engine.MainLoop() {
		if input.KeyPress('`') && player == nil {
      player = scene.SpawnBot("tim")
    } else if input.KeyPress('`') {
      scene.SpawnBot("tim2")
    }
    if player != nil {
      speed := float32(5.0)
      t := player.Transform()
      var move engine.Vector = player.Transform().WorldPosition()
      if input.KeyDown('W') {
        t.SetRotationf(0.0)
        move.Y += speed
      }
      if input.KeyDown('S') {
        t.SetRotationf(180.0)
        move.Y += -speed
      }
      if input.KeyDown('A') {
        t.SetRotationf(90.0)
        move.X += -speed
      }
      if input.KeyDown('D') {
        t.SetRotationf(270.0)
        move.X += speed
      }
      t.SetWorldPosition(move)

      if input.MouseDown(input.MouseLeft) {
        player.Shoot()
      }
    }
  }
}
