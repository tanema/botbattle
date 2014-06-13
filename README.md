botbattle
=========

A battle arena built in go for a upcoming go hack meetup.

##Server
just run `go run app.go`

##Client/Bot

```go
package main

import (
	"botbattle/client"
)

func main() {
	botclient := client.NewBotClient("localhost:3333", "sir killalot")
}
```

###Client API
The General API is 

```go

func NewBotClient(host, botname string) (*BotClient)
type Status
	X            int
	Y            int
	Rotation     int
	Health       int
type BotClient
  ArenaHeight int
  ArenaWidth  int
  func Register(name string) (arena_width, arena_height int)
  func MoveForward() (x, y int)
  func MoveBackward() (x, y int)
  func RotLeft() (rotation int)
  func RotRight() (rotation int)
  func Scan() ([]*Status)
  func Status() (*Status)
  func Shield() bool

```


TODO
bot collision damage
reconnect option?
kill counter
