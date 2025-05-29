package main

import (
	"fmt"
	"log"
	"net"
)

var messages = make(chan Packet)

// this function currently use `sendPacket` but instead it should use the current
func handleConnection(conn net.Conn) error {
	var err error
	defer conn.Close()
	defer fmt.Println("End of connection")
	defer fmt.Println(err)
	fmt.Println("Connection incoming")

	response := make([]byte, 1024)
	var n int
	n, err = conn.Read(response)
	if err != nil {
		return nil
	}
	bytesResponse := response[:n]
	var packet Packet
	packet, err = Parse(bytesResponse)
	if err != nil {
		return err
	}
	fmt.Println(packet)

	// Handle packet in a switch
	switch packet.(type) {
	case *AskPacket:
		// give command in stock
		select {
		case message := <-messages:
			_, err = conn.Write(message.ToBytes())
			if err != nil {
				return err
			}

		default:
			// No message in pocket
			empty := EmptyPacket{
				Header: Header{
					PacketType: EMPTY,
					Source:     "",
					Hop:        0,
				},
			}
			_, err = conn.Write(empty.ToBytes())
			if err != nil {
				return err
			}
		}
	default:
		messages <- packet
		fmt.Println(packet)
	}

	_, err = fmt.Fprintf(conn, "%s", response)
	if err != nil {
		return err
	}
	return nil
}

func listen(host Host) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host.ip, host.port))
	if err != nil {
		log.Fatal("Failed to start listener:", err)
	}
	defer listener.Close()
	fmt.Printf("Listening on %s:%s...\n", host.ip, host.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
