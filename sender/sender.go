package sender

import (
	"net"
	// "github.com/jaysoncena/bla/common"
)

// Sender TODO comment
type Sender struct {
	Conn       *net.UDPConn
	targetAddr string
	UDPAddr    *net.UDPAddr
}

// NewSender creates Sender object
func NewSender(addr string) (*Sender, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &Sender{
		targetAddr: addr,
		Conn:       conn,
		UDPAddr:    udpAddr,
	}, nil
}
