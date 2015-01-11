package main

import (
	"fmt"
	"net/http"

	"github.com/tanema/botbattle/game"
)

// main setup of the entire server/app
func main() {
  // start the arena serverside emulation
	game.NewScene().Start()
  //handle all the javascript html and css
	http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Println("Listening")
  //start the server
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
