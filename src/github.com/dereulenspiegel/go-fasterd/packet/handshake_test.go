package packet

import (
  "encoding/binary"
  "testing"
  "github.com/stretchr/testify/assert"
)

func createValidPacket() []byte {
  buf := make([]byte,9)
  buf[0] = 0x01
  binary.PutUvarint(buf[2:4],5)
  binary.PutUvarint(buf[4:6],0x0000)
  binary.PutUvarint(buf[6:8], 1)
  buf[8] = byte(1)
  return buf
}

func TestUnmarshallHandshake(t *testing.T) {
  assert := assert.New(t)
  validPacket := createValidPacket()
  packet, err := UnmarshallHandshakePacket(validPacket, nil)

  assert.Nil(err)
  assert.NotNil(packet)
  assert.Equal(1, len(packet.TLVRecords))
  assert.Equal(uint16(5), packet.header.tlvRecordLength)
  assert.Equal(HandshakeType, packet.TLVRecords[0].Type)
  assert.Equal(uint16(1), packet.TLVRecords[0].Length)
  assert.Equal(byte(1), packet.TLVRecords[0].Body[0])
}
