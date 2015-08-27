package packet

import (
  "encoding/binary"
  "testing"
  "github.com/stretchr/testify/assert"
)

func createValidPayload() []byte {
  buf := make([]byte, 1024+24)
  buf[0] = 0x02
  binary.PutUvarint(buf[1:7], 5)
  for  i := 8; i< 24; i++ {
    buf[i] = byte(i-8)
  }
  for i := 24; i<1024; i++ {
    buf[i] = 0xFF
  }
  return buf
}

func TestUnmarshallValidPayloadPacket(t *testing.T){
  assert := assert.New(t)

  validPacket := createValidPayload()
  packet, err := UnmarshallPayloadPacket(validPacket, false)
  assert.Nil(err)
  assert.NotNil(packet)
  assert.Equal(PAYLOAD_PACKET_TYPE, packet.Header.PacketType())
}
