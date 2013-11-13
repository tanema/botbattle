package main

import (
  "flag"
  "runtime"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
  "github.com/tanema/botbattle/scene"
  "github.com/tanema/botbattle/server"
)

const rotSpeed = float32(250.0)

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
		if input.KeyPress('`') && (player == nil || player.GameObject() == nil) {
      player = scene.SpawnBot("tim", nil)
    } else if input.KeyPress('`') {
      scene.SpawnBot("tim2", nil)
    }
    if player != nil && player.GameObject() != nil {
      if input.KeyDown('W') {
        player.Forward()
      }
      if input.KeyDown('Q') {
        player.Stop()
      }
      if input.KeyDown('S') {
        player.Backward()
      }
      if input.KeyDown('D') {
        player.Rotate(-5.0)
      }
      if input.KeyDown('A') {
        player.Rotate(5.0)
      }
      if input.KeyDown(' ') {
        player.Shoot()
      }
      if input.KeyDown('E') {
        player.Scan()
      }
      if input.KeyDown('G') {
        player.OnDie(false)
      }
    }
  }
}
