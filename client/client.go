package client

import (
  "bufio"
  "fmt"
  "net"
  "time"
)

type Client struct {
  host string
  conn net.Conn
  events chan string
}

func NewClient(host, name string, events chan string) Client {
  conn, _ := net.Dial("tcp", host)
  go handleMessage(events, conn)
  conn.Write([]byte("REGISTER " + name + "\n"))
  return Client{host, conn, events}
}

func (c *Client) sendMessage(cmd, arguments string) {
  _, err := c.conn.Write([]byte(cmd + " " + arguments + "\n"))
  if err != nil {
    c.events <- "DISCONNECT"
    return
  }
  time.Sleep(time.Millisecond)
}

func handleMessage(events chan string, conn net.Conn){
  for {
    line, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      events <- "DISCONNECT"
      return
    }
    events <- line
  }
}

func (c *Client) Forward() {
  c.sendMessage("FORWARD", "")
}

func (c *Client) Backward() {
  c.sendMessage("BACKWARD", "")
}

func (c *Client) Stop() {
  c.sendMessage("STOP", "")
}

func (c *Client) Shoot() {
  c.sendMessage("SHOOT", "")
}

func (c *Client) Scan() {
  c.sendMessage("SCAN", "")
}

func (c *Client) Rotate(deg float32) {
  c.sendMessage("ROTATE", fmt.Sprintf("%f", deg))
}

func (c *Client) RotateTo(deg float32) {
  c.sendMessage("ROTATE_TO", fmt.Sprintf("%f", deg))
}

func (c *Client) GetCurrentPosition() {
  c.sendMessage("CURRENT_POS", "")
}
