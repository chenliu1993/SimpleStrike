package utils

import (
	"fmt"
	"net"
)

const (
	// ClientPort is where client.
	ClientPort int = 8080
	// ServerPort is service.
	ServerPort int = 8808
)

// TransferServer serves as a mid-server for transferring messages.
type TransferServer struct {
	ClientLis   *net.TCPListener
	TransferLis *net.TCPListener
	Channels    map[int]*Channel
	CurChID     int
}

// Channel used as marking channels for TransferServer.
type Channel struct {
	ID              int
	Client          net.Conn
	Transfer        net.Conn
	ClientRecvMsg   chan []byte
	TransferSendMsg chan []byte
}

// NewServer news a transferserver
func NewServer() *TransferServer {
	return &TransferServer{
		Channels: make(map[int]*Channel),
		CurChID:  0,
	}
}

// StartService starts msg transfer.
func (t *TransferServer) StartService(clientPort, transferPort int) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", clientPort))
	if err != nil {
		return err
	}
	t.ClientLis, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", transferPort))
	if err != nil {
		return err
	}
	t.TransferLis, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	go t.AcceptLoop()
	return nil
}

// StopService stops service.
func (t *TransferServer) StopService() {
	t.ClientLis.Close()
	t.TransferLis.Close()
	for _, v := range t.Channels {
		v.Transfer.Close()
		v.Client.Close()
	}
}

// DelChannel removes a channel
func (t *TransferServer) DelChannel(id int) {
	chs := t.Channels
	delete(chs, id)
	t.Channels = chs
}

// AcceptLoop deals with conn.
func (t *TransferServer) AcceptLoop() {
	transfer, err := t.TransferLis.Accept()
	if err != nil {
		return
	}
	for {
		client, err := t.ClientLis.Accept()
		if err != nil {
			continue
		}
		ch := &Channel{
			ID:              t.CurChID,
			Client:          client,
			Transfer:        transfer,
			ClientRecvMsg:   make(chan []byte),
			TransferSendMsg: make(chan []byte),
		}
		t.CurChID++
		chs := t.Channels
		chs[ch.ID] = ch
		t.Channels = chs
		go t.ClientMsgLoop(ch)
		go t.TransferMsgLoop(ch)
		go t.MsgLoop(ch)
	}
}

// ClientMsgLoop deals client msg.
func (t *TransferServer) ClientMsgLoop(ch *Channel) {
	defer func() {
		fmt.Println("ClientMsgLoop exists")
	}()
	for {
		select {
		case data, isClose := <-ch.TransferSendMsg:
			{
				if !isClose {
					return
				}
				_, err := ch.Client.Write(data)
				if err != nil {
					return
				}
			}
		}
	}
}

// TransferMsgLoop deals with backend service
func (t *TransferServer) TransferMsgLoop(ch *Channel) {
	defer func() {
		fmt.Println("TransferMsgLoop exit")
	}()
	for {
		select {
		case data, isClose := <-ch.ClientRecvMsg:
			{
				if !isClose {
					return
				}
				_, err := ch.Transfer.Write(data)
				if err != nil {
					return
				}
			}
		}
	}
}

// MsgLoop deals with all msg.
func (t *TransferServer) MsgLoop(ch *Channel) {
	defer func() {
		close(ch.ClientRecvMsg)
		close(ch.TransferSendMsg)
		t.DelChannel(ch.ID)
		fmt.Println("MsgLoop exit")
	}()
	buf := make([]byte, 1024)
	for {
		n, err := ch.Client.Read(buf)
		if err != nil {
			return
		}
		ch.ClientRecvMsg <- buf[:n]
		n, err = ch.Transfer.Read(buf)
		if err != nil {
			return
		}
		ch.TransferSendMsg <- buf[:n]
	}
}
