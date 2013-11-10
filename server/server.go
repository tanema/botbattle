package server

import (
  "fmt"
  "net"
  "bufio"
  "strings"
)

type server struct {
  host string
  send chan string
  connection_list map[net.Addr]net.Conn
}

func NewServer(host string) server {
  return server{host, make(chan string), make(map[net.Addr]net.Conn)}
}

func (s *server) Listen(){
  fmt.Printf("listening on %s \n", s.host)
  listen, _ := net.Listen("tcp", s.host)
  go s.handleMessages()
  for {
    conn, err := listen.Accept()
    if err != nil {
      continue
    }
    go s.handlePeer(conn)
  }
}

func (s *server) handleMessages() {
  for s := range s.send {
    bits := strings.SplitN(s, " ", 3)
    if len(bits) != 3 {
      fmt.Printf("Error: invalid line: %v\n", bits)
      continue
    }
    //handle action
  }
}

func (s *server) handlePeer(conn net.Conn){
  fmt.Printf("Peer %s connected", conn.RemoteAddr().String())
  s.connection_list[conn.RemoteAddr()] = conn
  defer func(){
    conn.Close()
    delete(s.connection_list, conn.RemoteAddr())
  }()

  for {
    line, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      println("disconnecting peer")
      return
    }
    s.send <- line
  }
}
