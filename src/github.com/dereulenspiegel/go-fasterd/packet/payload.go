package packet

import (
  "fmt"
  "encoding/binary"
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
  Header *PayloadHeader
  Payload []byte
}

func(packet PayloadPacket) Length() int {
    return packet.Header.Length() + len(packet.Payload)
}

func(packet PayloadPacket) Marshall() []byte {
  totalSize := packet.Header.Length() + len(packet.Payload)
  buf := make([]byte, totalSize)
  copy(buf[0:packet.Header.Length()-1],packet.Header.Marshall())
  copy(buf[packet.Header.Length():],packet.Payload)
  return buf
}

func UnmarshallPayloadPacket(buf []byte, nullMethod bool) (packet PayloadPacket,err error) {
  if buf[0] != byte(PAYLOAD_PACKET_TYPE) {
    err = fmt.Errorf("Invalid packet type %d ",buf[0])
    return
  }
  header := &PayloadHeader{}
  packet = PayloadPacket{Header: header}
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
  }
  return
}
