package main

import (
  "github.com/dereulenspiegel/go-fasterd/packet"
  "github.com/dereulenspiegel/go-fasterd/pipeline"
  "github.com/dereulenspiegel/go-fasterd/peer"
)

type HandshakePipe struct {
  SendPipeline *pipeline.PacketPipeline
}

func (pipe *HandshakePipe) Process(in chan packet.Packet) chan packet.Packet {
  out := make(chan packet.Packet)
  go func() {
    for pkg := range in {
      if pkg.Header().PacketType() == packet.PAYLOAD_PACKET_TYPE {
        // Don't process any payload packets, just push them though
        out <- pkg
      } else {
        pipe.ProcessHandshake(pkg)
      }
    }
  }()
  return out
}

func (pipe *HandshakePipe) ProcessHandshake(pkg packet.Packet) {
  peerState, exists := peer.GlobalPeerState[pkg.PeerAddress()]
  if !exists {
    peerState = &peer.PeerState{}
  }
  
}
