package main

import (
  "os"
  "fmt"
  "strings"
  "github.com/tanema/botbattle/client"
)

func main() {
  events := make(chan string)
  c := client.NewClient("localhost:4569", "iTim", events)
  go handleEvents(events)
  for {
    c.Forward()
    c.Rotate(1.0)
    c.Scan()
    c.Shoot()
  }
}

func handleEvents(events chan string) {
  for line := range events {
    line = strings.Replace(line, "\n", "", -1)
    bits := strings.SplitN(line, " ", 2)
    switch bits[0] {
    case "ON_SCAN":
      fmt.Println(bits[1])
    case "ON_DIE":
      os.Exit(0)
    case "ON_CURRENT_POS":
      fmt.Println(bits[1])
    }
  }
  os.Exit(0)
}
