package client

import (
  "os"
  "bufio"
  "fmt"
  "net"
  "time"
  "strings"
  "strconv"
)

type Client struct {
  host string
  conn net.Conn
  events chan string
}

type scan_cb func(x, y float64, name string)
type current_pos_cb func(x, y float64)

func NewClient(host, name string) Client {
  conn, err := net.Dial("tcp", host)
  if err != nil {
    println("could not connect to that host")
    os.Exit(0)
  }
  events := make(chan string)
  go handleMessage(events, conn)
  conn.Write([]byte("REGISTER " + name + "\n"))
  return Client{host, conn, events}
}

func (c *Client) sendMessage(cmd, arguments string) {
  _, err := c.conn.Write([]byte(cmd + " " + arguments + "\n"))
  if err != nil {
    os.Exit(0)
  }
  time.Sleep(time.Millisecond)
}

func handleMessage(events chan string, conn net.Conn){
  for {
    line, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      os.Exit(0)
    }
    line = strings.Replace(line, "\n", "", -1)
    bits := strings.SplitN(line, " ", 2)
    switch bits[0] {
    case "ON_DIE":
      println("you died")
      os.Exit(0)
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

func (c *Client) Rotate(deg float32) {
  c.sendMessage("ROTATE", fmt.Sprintf("%f", deg))
}

func (c *Client) RotateTo(deg float32) {
  c.sendMessage("ROTATE_TO", fmt.Sprintf("%f", deg))
}

func (c *Client) Scan(cb scan_cb) {
  _, err := c.conn.Write([]byte("SCAN \n"))
  if err != nil {
    os.Exit(0)
  }
  for line := range c.events {
    line = strings.Replace(line, "\n", "", -1)
    bits := strings.SplitN(line, " ", 2)
    switch bits[0] {
    case "ON_SCAN":
      items := strings.Split(bits[1], ":")
      x, _ := strconv.ParseFloat(items[0], 64)
      y, _ := strconv.ParseFloat(items[1], 64)
      cb(x, y, items[2])
      return
    }
  }
}

func (c *Client) GetCurrentPosition(cb current_pos_cb) {
  _, err := c.conn.Write([]byte("CURRENT_POS \n"))
  if err != nil {
    os.Exit(0)
  }
  for line := range c.events {
    line = strings.Replace(line, "\n", "", -1)
    bits := strings.SplitN(line, " ", 2)
    switch bits[0] {
    case "ON_CURRENT_POS":
      items := strings.Split(bits[1], ":")
      x, _ := strconv.ParseFloat(items[0], 64)
      y, _ := strconv.ParseFloat(items[1], 64)
      cb(x, y)
      return
    }
  }
}
