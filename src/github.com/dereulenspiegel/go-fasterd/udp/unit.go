package udp

import (
  "net"
)

type UDPUnit struct {
  Data []byte
  Address *net.UDPAddr
  Error error
}

func NewUDPUnit(data []byte, address *net.UDPAddr, err error) *UDPUnit {
  return &UDPUnit {
    Data: data,
    Address: address,
    Error: err,
  }
}
