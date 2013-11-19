package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	host       string
	conn       net.Conn
	scan_event chan string
	pos_event  chan string
}

type scan_cb func(x, y float64, name string)
type current_pos_cb func(x, y float64)

func NewClient(host, name string) Client {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		println("could not connect to that host")
		fmt.Println(err)
		os.Exit(0)
	}
	scan_event := make(chan string)
	pos_event := make(chan string)
	go handleMessage(scan_event, pos_event, conn)
	conn.Write([]byte("REGISTER " + name + "\n"))
	return Client{host, conn, scan_event, pos_event}
}

func (c *Client) sendMessage(cmd, arguments string) {
	_, err := c.conn.Write([]byte(cmd + " " + arguments + "\n"))
	if err != nil {
		os.Exit(0)
	}
	time.Sleep(time.Millisecond)
}

func handleMessage(scan_event, pos_event chan string, conn net.Conn) {
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
		case "ON_SCAN":
			scan_event <- bits[1]
		case "ON_CURRENT_POS":
			pos_event <- bits[1]
		}
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
	for line := range c.scan_event {
		items := strings.Split(line, ":")
		x, _ := strconv.ParseFloat(items[0], 64)
		y, _ := strconv.ParseFloat(items[1], 64)
		cb(x, y, items[2])
		return
	}
}

func (c *Client) GetCurrentPosition(cb current_pos_cb) {
	_, err := c.conn.Write([]byte("CURRENT_POS \n"))
	if err != nil {
		os.Exit(0)
	}
	for line := range c.pos_event {
		items := strings.Split(line, ":")
		x, _ := strconv.ParseFloat(items[0], 64)
		y, _ := strconv.ParseFloat(items[1], 64)
		cb(x, y)
		return
	}
}
