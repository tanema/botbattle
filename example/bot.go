package main

import (
	"botbattle/client"
	"fmt"
)

func main() {
  client := client.NewBotClient("localhost:3333", "Tim")
	fmt.Println(client.ArenaWidth, client.ArenaHeight)
	status := client.Status()
	fmt.Println(status.X, status.Y, status.Rotation, status.Health)

	shot := client.FireGun()
	fmt.Println(shot)

	x, y := client.MoveForward()
	fmt.Println(x, y)
	x, y = client.MoveForward()
	fmt.Println(x, y)
	x, y = client.MoveForward()
	fmt.Println(x, y)

	bots := client.Scan()
	fmt.Println(bots)

	rot := client.RotLeft()
	fmt.Println(rot)

	x, y = client.MoveBackward()
	fmt.Println(x, y)
	x, y = client.MoveBackward()
	fmt.Println(x, y)
	x, y = client.MoveBackward()
	fmt.Println(x, y)

	rot = client.RotRight()
	fmt.Println(rot)

	shot = client.FireCannon()
	fmt.Println(shot)
}
