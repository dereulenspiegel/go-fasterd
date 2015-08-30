package packet

import (
  "encoding/binary"
  "fmt"
  "net"
)

type HandshakeHeader struct {
  tlvRecordLength uint16
}

func (header *HandshakeHeader) PacketType() PacketType {
  return HANDSHAKE_PACKET_TYPE
}

func (header *HandshakeHeader) Length() int{
  return 4
}

func (header *HandshakeHeader) Marshall() []byte {
  buf := make([]byte,4)
  buf[0] = byte(HANDSHAKE_PACKET_TYPE)
  buf[1] = 0x00
  binary.PutUvarint(buf[2:4],uint64(header.tlvRecordLength))
  return buf
}

type HandshakePacket struct {
  header *HandshakeHeader
  TLVRecords map[TLVRecordType]*TLVRecord
  peerAddr *net.UDPAddr
}

func (packet HandshakePacket) PeerAddress() *net.UDPAddr{
  return packet.peerAddr
}

func (packet *HandshakePacket) Header() Header {
  return packet.header
}

func (packet *HandshakePacket) Marshall() []byte {
  headerBuf := packet.header.Marshall()
  recordsBuf := make([]byte, packet.header.tlvRecordLength)

  var pointer uint16 = 0
  for _, record := range packet.TLVRecords {
    recordBuf := make([]byte,record.Length + 4)
    binary.PutUvarint(recordBuf[0:2], uint64(record.Type))
    binary.PutUvarint(recordBuf[2:4], uint64(record.Length))
    copy(recordBuf[4:len(recordBuf)-1], record.Body)
    pointer = pointer + record.Length + 4
    copy(recordsBuf[pointer:record.Length+4],recordBuf)
  }

  return append(headerBuf, recordsBuf...)
}

func (packet HandshakePacket)Length() int {
  return packet.header.Length() + int(packet.header.tlvRecordLength)
}

func (packet *HandshakePacket) AddTLVRecord(record *TLVRecord) {
  length := record.Length + 4
  packet.TLVRecords[record.Type] = record
  packet.header.tlvRecordLength = packet.header.tlvRecordLength + length
}

func UnmarshallHandshakePacket(buf []byte, addr *net.UDPAddr)(packet *HandshakePacket, err error) {
  if buf[0] != byte(HANDSHAKE_PACKET_TYPE) {
    err = fmt.Errorf("Invalid packet type %d", buf[0])
    return
  }
  tlvRecordLength, _ := binary.Uvarint(buf[2:4])
  header := &HandshakeHeader{
    tlvRecordLength: uint16(tlvRecordLength),
  }
  packet = &HandshakePacket{
    header: header,
    peerAddr: addr,
    TLVRecords: make(map[TLVRecordType]*TLVRecord),
  }
  for pointer := 4; pointer < int(tlvRecordLength) ; {
    typeValue, _ := binary.Uvarint(buf[pointer:pointer+2])
    lengthValue, _ := binary.Uvarint(buf[pointer+2:pointer+4])
    bodyBuf := make([]byte, lengthValue)
    copy(bodyBuf,buf[pointer+4:(pointer+4+int(lengthValue))])
    record := &TLVRecord{
      Type: TLVRecordType(typeValue),
      Length: uint16(lengthValue),
      Body: bodyBuf,
    }
    packet.TLVRecords[record.Type] = record
    pointer = pointer + int(lengthValue) + 4
  }
  return
}
