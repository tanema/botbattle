package scene

import (
  "fmt"
  "net"
  "bufio"
)

type server struct {
  host string
}

func NewServer(host string) server {
  return server{host}
}

func (s *server) Listen(){
  fmt.Printf("listening on %s \n", s.host)
  listen, _ := net.Listen("tcp", s.host)
  for {
    conn, err := listen.Accept()
    if err != nil {
      continue
    }
    s.handlePeer(conn)
  }
}

func (s *server) handlePeer(conn net.Conn){
  fmt.Printf("Peer %s connected", conn.RemoteAddr().String())
  defer func(){
    conn.Close()
  }()

  peer := NewPeer(conn)
  for {
    line, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      peer.Player.OnDie(false)
      return
    }
    peer.OnMessages(line)
  }
}
