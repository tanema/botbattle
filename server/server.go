package server

import (
  "github.com/tanema/botbattle/scene"
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
  for {
    conn, err := listen.Accept()
    if err != nil {
      continue
    }
    go s.handlePeer(conn)
  }
}

func (s *server) handleMessages(conn net.Conn, message string) {
  bits := strings.SplitN(message, " ", 3)
  switch bits[0] {
  case "SET_NAME":
    scene.SpawnBot(bits[1], conn)
  }
  //_, err := conn.Write([]byte("put_message_id_here" + " " + "response_goes_here" + "\n"))
  //if err != nil { fmt.Printf("error writing out to connection: %s \n", err) }
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
    s.handleMessages(conn, line)
  }
}
