package conn4api

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func adr(port int) *net.UDPAddr {
	serverAdr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	checkError(err)
	return serverAdr
}

func connection(port int) *net.UDPConn {
	serverAdr := adr(port)
	serverConn, err := net.ListenUDP("udp", serverAdr)
	checkError(err)
	return serverConn
}

func readFromUDP(connection *net.UDPConn, channel chan string) {
	reply := make([]byte, 1024)
	n, _, err := connection.ReadFromUDP(reply)
	checkError(err)
	reply = reply[:n]
	channel <- string(reply)
}

func TestUdpConnectionSend(t *testing.T) {
	serverPort := 13337
	localPort := 13338
	server := connection(serverPort)

	channel := make(chan string)
	go readFromUDP(server, channel)

	conn := NewUdpConnection(fmt.Sprintf(":%d", serverPort), localPort)

	conn.Run(func(sender chan string, receiver chan string) {
		sender<- "TEST"
	})

	select {
		case <-time.After(5 * time.Second):
			t.Error("timeout")
		case msg := <-channel : {
			if msg != "TEST" {
				t.Error("ERROR")
			}
		}
	}
	close(channel)
	server.Close()
	conn.Close()
}

func TestUdpConnectionReceive(t *testing.T) {
	serverPort := 13331
	localPort := 13332
	server := connection(serverPort)
	channel := make(chan string)

	conn := NewUdpConnection(fmt.Sprintf(":%d", serverPort), localPort)
	go conn.registerReceiver(channel)

	server.WriteToUDP([]byte("TEST"), adr(localPort))

	select {
		case <-time.After(1 * time.Second):
			t.Error("timeout")
		case msg := <-channel : {
			if msg != "TEST" {
				t.Error("ERROR")
			}
		}
	}

	close(channel)
	server.Close()
	conn.Close()
}

