package packet

import (
  "fmt"
  "net"
  "github.com/dereulenspiegel/go-fasterd/udp"
)

type PacketType byte

const (
  HANDSHAKE_PACKET_TYPE PacketType = 0x01
  PAYLOAD_PACKET_TYPE PacketType = 0x02
)

type Header interface {
  PacketType() PacketType
  Length() int
  Marshall() []byte
}

type Packet interface {
  Marshall() []byte
  Length() int
  Header() Header
  PeerAddress() *net.UDPAddr
}

func Unmarshall(data []byte, addr *net.UDPAddr) (Packet, error) {
  switch(PacketType(data[0])){
  case HANDSHAKE_PACKET_TYPE:
    return UnmarshallHandshakePacket(data, addr)
  case PAYLOAD_PACKET_TYPE:
    return UnmarshallPayloadPacket(data, addr, false)
  default:
    return nil, fmt.Errorf("Unknwon Packet type %d", data[0])
  }
}

type PacketUnmarshaller struct {}

func (um PacketUnmarshaller) Process(in chan *udp.UDPUnit) chan Packet {
  out := make(chan Packet)
  go func(){
    for unit := range in {
      packet, _ := Unmarshall(unit.Data, unit.Address)
      out <- packet
    }
  }()
  return out
}
