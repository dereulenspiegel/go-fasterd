package packet

import (
  "fmt"
  "encoding/binary"
  "net"
)

type PayloadHeader struct {
  flags byte
  sequenceNumber uint64
  authenticationTag []byte
}

func(header *PayloadHeader) PacketType() PacketType {
  return PAYLOAD_PACKET_TYPE
}

func(header PayloadHeader) Length() int {
  if(header.sequenceNumber == 0) {
    return 1
  } else {
    return 24
  }
}

func(header PayloadHeader) Marshall() []byte {
  buf := make([]byte, header.Length())
  buf[0] = byte(header.PacketType())
  if(header.Length() > 1){
    sequenceBuf := make([]byte,6)
    binary.PutUvarint(sequenceBuf,header.sequenceNumber)
    copy(buf[1:8], sequenceBuf)
    copy(buf[8:23],header.authenticationTag)
  }
  return buf
}

type PayloadPacket struct {
  header *PayloadHeader
  Payload []byte
  peerAddr *net.UDPAddr
}

func (packet PayloadPacket) PeerAddress() *net.UDPAddr {
  return packet.peerAddr
}

func (packet PayloadPacket) Header() Header {
  return packet.header
}

func(packet PayloadPacket) Length() int {
    return packet.header.Length() + len(packet.Payload)
}

func(packet PayloadPacket) Marshall() []byte {
  totalSize := packet.header.Length() + len(packet.Payload)
  buf := make([]byte, totalSize)
  copy(buf[0:packet.header.Length()-1],packet.header.Marshall())
  copy(buf[packet.header.Length():],packet.Payload)
  return buf
}

func UnmarshallPayloadPacket(buf []byte, addr *net.UDPAddr, nullMethod bool) (packet PayloadPacket,err error) {
  if buf[0] != byte(PAYLOAD_PACKET_TYPE) {
    err = fmt.Errorf("Invalid packet type %d ",buf[0])
    return
  }
  header := &PayloadHeader{}
  packet = PayloadPacket{header: header}
  if(nullMethod){
    header.flags = 0x00
    header.sequenceNumber = 0
    header.authenticationTag = make([]byte,0)
    packet.Payload = buf[1:len(buf)-1]
  } else {
    header.flags = buf[1]
    sequenceNumber, n := binary.Uvarint(buf[2:8])
    if n < 1 {
      err = fmt.Errorf("Problem reading Uvarint")
    }
    header.sequenceNumber = sequenceNumber
    header.authenticationTag = buf[8:24]
    packet.Payload = buf[24:]
    packet.peerAddr = addr
  }
  return
}
