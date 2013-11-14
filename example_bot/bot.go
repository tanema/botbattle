package main

import (
  "github.com/tanema/botbattle/client"
)

func main() {
  events := make(chan string)
  c := client.NewClient("localhost:4569", events)
  c.SetName("Tim")

  for {}
}
