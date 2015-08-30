package pipeline

import (
  "github.com/dereulenspiegel/go-fasterd/packet"
  "github.com/dereulenspiegel/go-fasterd/udp"
)

type PacketPipe interface {
  Process(pChan chan packet.Packet) chan packet.Packet
}

type PacketPipeline struct {
  head chan *udp.UDPUnit
  tail chan packet.Packet
}

func (pipeline PacketPipeline) Enqueue (unit *udp.UDPUnit) {
  pipeline.head <- unit
}

func (pipeline PacketPipeline) Dequeue(handler func(packet.Packet)) {
  for packet := range pipeline.tail {
    handler(packet)
  }
}

func NewPacketPipeline(pipes... PacketPipe) *PacketPipeline{
  head := make(chan *udp.UDPUnit)
  unmarshaller := packet.PacketUnmarshaller{}
  next_chan := unmarshaller.Process(head)
  for _, pipe := range pipes {
    next_chan = pipe.Process(next_chan)
  }
  return &PacketPipeline{head: head, tail: next_chan }
}
