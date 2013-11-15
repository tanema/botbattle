package main

import (
  "github.com/tanema/botbattle/client"
)

func main() {
  c := client.NewClient("localhost:4569", "iTim")
  for {
    c.Forward()
    c.Rotate(1.0)
    go c.Scan(func(x, y float64, name string){
      println(name)
    })
    go c.GetCurrentPosition(func(x, y float64){
      println("pos")
    })
    c.Shoot()
  }
}
