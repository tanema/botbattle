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
	botclient, err := client.NewBotClient("localhost:3333", "sir killalot")
  ... //check out example/bot.go for an example on how to use it
}
```

###Client API
```go
func NewBotClient(host, botname string) (*BotClient, error)
type Status
  Id           int
  Name         string
	X            int
	Y            int
	Rotation     int
	Health       int
  KillCount    int
type BotClient
  ArenaHeight int
  ArenaWidth  int
  func Reconnect() (*BotClient, error)
  func MoveForward() (*Status, error)
  func MoveBackward() (*Status, error)
  func FireGun() (bool, error)
  func FireCannon() (bool, error)
  func RotateLeft() (*Status, error)
  func RotateRight() (*Status, error)
  func Scan() ([]*Status, error)
  func Status() (*Status, error)
  func Shield() (bool, error)
```

####func MoveForward() (\*Status, error)

- Moves forward in the direction that you are facing
- has a delay of 500 milleseconds
- returns 
    - current status of bot 
    - error if you have been killed or disconnected

####func MoveBackward() (\*Status, error)

- Moves backward in the direction that you are facing
- has a delay of 500 milleseconds
- returns 
    - current status of bot 
    - error if you have been killed or disconnected

####func FireGun() (bool, error)

- Will return true if the bullet hit somethin
- has a damage of 25
- has a delay of 1000 milleseconds
- returns 
    - bool of the success of the shot
    - error if you have been killed or disconnected

####func FireCannon() (bool, error)

- Will return true if the bullet hit somethin
- has a damage of 50
- has a delay of 3000 milleseconds
- returns 
    - bool of the success of the shot
    - error if you have been killed or disconnected

####func RotateLeft() (\*Status, error)

- Rotates -90 degrees
- has a delay of 500 milleseconds
- returns 
    - current status of bot 
    - error if you have been killed or disconnected

####func RotateRight() (\*Status, error)

- Rotates 90 degrees
- has a delay of 500 milleseconds
- returns 
    - current status of bot 
    - error if you have been killed or disconnected

####func Scan() ([]\*Status, error)

- Will return array of status's of the bots you can see
- returns empty array if you see nothing
- has a delay of 500 milleseconds
- returns 
    - current status of bots that you can see
    - error if you have been killed or disconnected

####func Status() (\*Status, error)

- Will return the status of your own bot please to refer to the status object for what info is included
- returns 
    - current status of bot 
    - error if you have been killed or disconnected

####func Shield() (bool, error)

- Will return true if the shield was enabled
- has a warmup time of 5000 millesecond
- will remain on for 3000 milleseconds
- returns 
    - bool of the success in enabling the shield
    - error if you have been killed or disconnected
