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
		if input.KeyPress('`') {
      scene.SpawnBot("tim")
    }
  }
}
