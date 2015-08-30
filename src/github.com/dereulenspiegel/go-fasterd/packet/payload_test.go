package packet

import (
  "encoding/binary"
  "testing"
  "github.com/stretchr/testify/assert"
)

func createValidPayload() []byte {
  buf := make([]byte, 1024+24)
  buf[0] = 0x02
  buf[1] = 0x00
  binary.PutUvarint(buf[2:8], 5)
  for  i := 8; i< 24; i++ {
    buf[i] = byte(i-8)
  }
  for i := 24; i<1024; i++ {
    buf[i] = 0xFF
  }
  return buf
}

func createPayloadWithInvalidHeader() []byte {
  buf := createValidPayload()
  buf[0] = 0x01
  return buf
}

func TestUnmarshallWithInvalidPacketType(t *testing.T) {
  assert := assert.New(t)
  invalidPacket := createPayloadWithInvalidHeader()
  _, err := UnmarshallPayloadPacket(invalidPacket, nil, false)

  assert.NotNil(err)
}

func TestUnmarshallValidPayloadPacket(t *testing.T){
  assert := assert.New(t)

  validPacket := createValidPayload()
  packet, err := UnmarshallPayloadPacket(validPacket, nil, false)
  assert.Nil(err)
  assert.NotNil(packet)
  assert.Equal(byte(0x0), packet.header.flags)
  assert.Equal(PAYLOAD_PACKET_TYPE, packet.header.PacketType())
  assert.Equal(uint64(5), packet.header.sequenceNumber)
  assert.Equal(24, packet.header.Length())
  assert.Equal(int(1024+24),packet.Length())
}
