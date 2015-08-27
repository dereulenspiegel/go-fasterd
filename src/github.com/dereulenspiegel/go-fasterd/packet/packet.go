package packet

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
  Header() Header
  Marshall() []byte
}
