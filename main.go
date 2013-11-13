package main

import (
  "flag"
  "runtime"
  "math"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
  "github.com/tanema/botbattle/scene"
  "github.com/tanema/botbattle/server"
)

const (
  rotSpeed = float32(250.0)
  RadianConst = math.Pi / 180
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
      speed := float64(0.0)
      t := player.Transform()

      if input.KeyDown('W') {
        speed = 5.0
      }
      if input.KeyDown('S') {
        speed = -5.0
      }

      rot := player.Transform().Rotation()
	    delta := float32(engine.DeltaTime())
      if input.KeyDown('D') {
        t.SetRotationf(rot.Z - rotSpeed*delta)
      }
      if input.KeyDown('A') {
        t.SetRotationf(rot.Z + rotSpeed*delta)
      }

      rot = player.Transform().Rotation()
      move := player.Transform().WorldPosition()
      move.X = float32(-math.Sin(float64(rot.Z)*RadianConst)*speed + float64(move.X))
      move.Y = float32(math.Cos(float64(rot.Z)*RadianConst)*speed + float64(move.Y))
      t.SetWorldPosition(move)

      if input.MouseDown(input.MouseLeft) {
        player.Shoot()
      }
    }
  }
}
