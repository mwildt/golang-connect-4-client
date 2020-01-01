package conn4api

import (
	"fmt"
	"net"
)

type Handler func(sender chan string, receiver chan string)

type UdpConnection struct {
	remoteAdr *net.UDPAddr
	localAdr *net.UDPAddr
	connection *net.UDPConn
	running bool
}

func NewUdpConnection(serverAddress string, port int) *UdpConnection {
	remoteAdr, err := net.ResolveUDPAddr("udp4", serverAddress)
	checkError(err)
	localAddress, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	checkError(err)
	connection, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	return &UdpConnection{
		remoteAdr: remoteAdr,
		localAdr: localAddress,
		connection: connection,
		running: true,
	}
}

func (connection *UdpConnection) send(message string) {
	_, err := connection.connection.WriteToUDP([]byte(message), connection.remoteAdr)
	checkError(err)
}

func (connection *UdpConnection) Run(handler Handler) {
	sender := make(chan string)
	receiver := make(chan string)
	go connection.registerReceiver(receiver)
	go connection.registerSender(sender)
	go handler(sender, receiver)
}

func (connection *UdpConnection) registerReceiver(receiver chan string) {
	for connection.running {
		reply := make([]byte, 1024)
		n, _, err := connection.connection.ReadFromUDP(reply)
		if connection.running {
			checkError(err)
			reply = reply[:n]
			receiver <- string(reply)
		}
	}
}

func (connection *UdpConnection) registerSender(sender chan string) {
	for connection.running {
		msg := <-sender
		connection.send(msg)
	}
}

func (connection *UdpConnection) Close() {
	connection.running = false
	connection.connection.Close()
}



