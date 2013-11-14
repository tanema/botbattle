package main

import (
  "github.com/tanema/botbattle/client"
)

func main() {
  events := make(chan string)
  c := client.NewClient("localhost:4569", "Tim Robot", events)
  c.Forward()
  for {
    c.Rotate(1.0)
    c.Scan()
    c.Shoot()
  }
}
