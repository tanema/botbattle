package client

import (
  "fmt"
  "net"
  "strconv"
  "time"
)

type Client struct {
  host string
  ident string
  conn net.Conn
}

func NewClient(host string, events chan string) Client {
  fmt.Printf("connecting to %s \n", host)
  conn, _ := net.Dial("tcp", host)
  go handleMessage(events)
  return Client{host, strconv.FormatInt(time.Now().Unix(), 10), conn}
}

func (c *Client) sendMessage(cmd, arguments string) {
  _, err := c.conn.Write([]byte(c.ident + " " + cmd + " " + arguments + "\n"))
  if err != nil {
    fmt.Printf("error writing out to connection: %s \n", err)
  }
}

func handleMessage(events chan string){
  for {
    line, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      println("disconnected")
      return
    }
    events <- line
  }
}
