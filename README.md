botbattle
=========

A battle arena built in go for a upcoming go hack meetup.

##Server

- To run the webserver/game arena run `go run app.go`
- open your web browser to `localhost:3000`


##Client/Bot

```go
package main

import (
	"botbattle/client"
)

func main() {
	botclient := client.NewBotClient("localhost:3333", "sir killalot")
  ... //check out example/bot.go for an example on how to use it
}
```

###Client API
```go
func NewBotClient(host, botname string) (*BotClient)
type Status
  Id           int
	X            int
	Y            int
	Rotation     int
	Health       int
type BotClient
  ArenaHeight int
  ArenaWidth  int
  func MoveForward() (x, y int)
  func MoveBackward() (x, y int)
  func FireGun() bool
  func FireCannon() bool
  func RotateLeft() (rotation int)
  func RotateRight() (rotation int)
  func Scan() ([]*Status)
  func Status() (*Status)
  func Shield() bool
```

####func MoveForward() (x, y int)

- Moves forward in the direction that you are facing
- has a delay of 500 milleseconds
  

####func MoveBackward() (x, y int)

- Moves backward in the direction that you are facing
- has a delay of 500 milleseconds

####func FireGun() bool

- Will return true if the bullet hit somethin
- has a damage of 25
- has a delay of 1000 milleseconds

####func FireCannon() bool

- Will return true if the bullet hit somethin
- has a damage of 50
- has a delay of 3000 milleseconds

####func RotateLeft() (rotation int)

- Rotates -90 degrees
- has a delay of 500 milleseconds

####func RotateRight() (rotation int)

- Rotates 90 degrees
- has a delay of 500 milleseconds

####func Scan() ([]*Status)

- Will return array of status's of the bots you can see
- returns empty array if you see nothing
- has a delay of 500 milleseconds

####func Status() (*Status)

- Will return the status of your own bot please to refer to the status object for what info is included

####func Shield() bool

- Will return true if the shield was enabled
- has a warmup time of 5000 millesecond
- will remain on for 3000 milleseconds

