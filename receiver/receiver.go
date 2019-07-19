package receiver

import (
	"fmt"
	"net"
)

// Receiver TODO comment
type Receiver struct {
	// iface      *net.Interface
	conn       *net.UDPConn
	listenAddr string
	UDPAddr    *net.UDPAddr
}

// NewReceiver creates Receiver object
func NewReceiver(addr string) (*Receiver, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		conn:       conn,
		listenAddr: addr,
		UDPAddr:    udpAddr,
	}, nil
}

// Listen listens
func (r *Receiver) Listen() {
	for {
		b := make([]byte, 256)
		_, _, err := r.conn.ReadFromUDP(b)
		if err != nil {
			return
		}
		fmt.Println("read", string(b))
	}
}
