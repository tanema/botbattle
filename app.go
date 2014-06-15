package main

import (
	"fmt"
	"net/http"

	"github.com/tanema/botbattle/game"
)

func main() {
	game.NewScene().Start()
	http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Println("Listening")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
