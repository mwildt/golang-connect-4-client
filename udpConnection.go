package main

import (
	"fmt"
	"net"
)

type Handler func(sender chan string, receiver chan string)

type UdpConnection struct {
	remoteAdr *net.UDPAddr
	localAdr *net.UDPAddr
	connection *net.UDPConn
}

func NewUdpConnection(serverAddress string, port int) *UdpConnection {
	remoteAdr, err := net.ResolveUDPAddr("udp4", serverAddress)
	checkError(err)
	localAddress, err :=  net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	checkError(err)
	connection, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	return &UdpConnection{
		remoteAdr: remoteAdr,
		localAdr: localAddress,
		connection: connection,
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
	for {
		reply := make([]byte, 1024)
		n, _, err := connection.connection.ReadFromUDP(reply)
		checkError(err)
		reply = reply[:n]
		receiver <- string(reply)
	}
}

func (connection *UdpConnection) registerSender(sender chan string) {
	for {
		msg := <-sender
		connection.send(msg)
	}
}


