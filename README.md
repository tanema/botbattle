botbattle
=========

A battle arena built in go for a upcoming go hack meetup.

##Server
To build the server you will need glew and glfw installed (along with a c++ compiler to install them) 

The default port it will run on is 4569 but you can change that with the flag -h for the host

##Client/Bot
It is best to start off by building and learning from the example_bot directory. 

###Client API
To connect and register your bot you will need to do the following

```go
package main

import (
  "github.com/tanema/botbattle/client"
)

func main() {
  bot := &Bot{client: client.NewClient("localhost:4569", "iRobot")}
  for {}
}
```

The General API is 

```go

NewClient(host, name string)
Client
  Forward() // will go forward until Stop() is called
  Backward() // will go backward until Stop() is called 
  Stop()
  Shoot()
  Rotate(degrees float32)
  RotateTo(degress float32)
  Scan(func(x, y float64, name string))
  GetCurrentPosition(func(x, y float64))

```
