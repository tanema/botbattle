package scene

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Peer struct {
	name   string
	conn   net.Conn
	Player *BotController
}
type cb func()

func NewPeer(conn net.Conn) *Peer {
	return &Peer{conn: conn}
}

func (p *Peer) OnMessages(line string) {
	line = strings.Replace(line, "\n", "", -1)
	bits := strings.SplitN(line, " ", 2)
	switch bits[0] {
	case "REGISTER":
		if p.Player != nil && p.Player.GameObject() != nil {
			p.Player.OnDie(false)
		}
		p.Player = SpawnBot(bits[1], p)
	case "FORWARD":
		p.onPlayer(p.Player.Forward)
	case "BACKWARD":
		p.onPlayer(p.Player.Backward)
	case "STOP":
		p.onPlayer(p.Player.Stop)
	case "SHOOT":
		p.onPlayer(p.Player.Shoot)
	case "SCAN":
		p.onPlayer(p.Player.Scan)
	case "ROTATE":
		if p.Player != nil && p.Player.GameObject() != nil {
			deg, err := strconv.ParseFloat(bits[1], 32)
			if err == nil {
				p.Player.Rotate(float32(deg))
			}
		}
	case "ROTATE_TO":
		if p.Player != nil && p.Player.GameObject() != nil {
			deg, err := strconv.ParseFloat(bits[1], 32)
			if err == nil {
				p.Player.RotateTo(float32(deg))
			}
		}
	case "CURRENT_POS":
		p.onPlayer(p.Player.GetCurrentPosition)
	}
}

func (p *Peer) onPlayer(with_player cb) {
	if p.Player != nil && p.Player.GameObject() != nil {
		with_player()
	}
}

func (p *Peer) OnScan(x, y float32, object string) {
	p.sendMessage("ON_SCAN", fmt.Sprintf("%g:%g:%s", x, y, object))
}

func (p *Peer) OnDie() {
	p.sendMessage("ON_DIE", "sorry")
}

func (p *Peer) OnCurrentPostition(x, y float32) {
	p.sendMessage("ON_CURRENT_POS", fmt.Sprintf("%g:%g", x, y))
}

func (p *Peer) sendMessage(cmd, arguments string) {
	_, err := p.conn.Write([]byte(cmd + " " + arguments + "\n"))
	if err != nil && p.Player != nil && p.Player.GameObject() != nil {
		p.Player.OnDie(false)
	}
}
