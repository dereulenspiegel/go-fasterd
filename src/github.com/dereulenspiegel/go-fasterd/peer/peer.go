package peer

import (
  "net"
  "sync"
  "github.com/dereulenspiegel/go-fasterd/method"
)

var GlobalPeerState map[*net.UDPAddr]*PeerState

var seqLock sync.Mutex
var peerIdSequence uint64

type PeerState struct {
  Address *net.UDPAddr
  Authenticated bool
  Method method.MethodType
  MethodState method.MethodState
  Lock       sync.Mutex
}

func GeneratePeerId() uint64 {
  // TODO this is primitive. We should be able to reuse freed ids
  seqLock.Lock()
  id := peerIdSequence
  peerIdSequence++
  seqLock.Unlock()
  return id
}
