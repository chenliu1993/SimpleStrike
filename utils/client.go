package utils

import (
	"fmt"
	"net"
)

const (
	// Host is the service addr.
	Host = "127.0.0.1"
	// LocalPort is the client access port.
	LocalPort int = 80
	// RemotePort is the server listen port.
	RemotePort int = 8808
)

// ClientHandler handles a conn to the server.
func ClientHandler(r net.Conn, localPort int) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			continue
		}
		data := buf[:n]
		local, err := net.Dial("tcp", fmt.Sprintf(":%d", localPort))
		if err != nil {
			continue
		}
		n, err = local.Write(data)
		if err != nil {
			continue
		}
		n, err = local.Read(buf)
		local.Close()
		if err != nil {
			continue
		}
		data = buf[:n]
		n, err = r.Write(data)
		if err != nil {
			continue
		}
	}
}
