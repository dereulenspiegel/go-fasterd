package main

import (
  "fmt"
  "net"
)

const MAX_PACKET_SIZE uint32 = 8192

type Server struct {
  port int
}

func NewServer(port int)(server *Server){
  return &Server{port: port}
}

func (server *Server) Listen() (err error){
  ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d", server.port))
  if err != nil {
    return
  }

  ServerConn, err := net.ListenUDP("udp", ServerAddr)
  if err != nil {
    return
  }
  defer ServerConn.Close()

  go func(){
    for {
      buf := make([]byte, MAX_PACKET_SIZE)
      _,_,err := ServerConn.ReadFromUDP(buf)
      if err != nil{
        // TODO Log this error, disconnect the client etc.
        continue
      }
    }
  }()
  return
}
