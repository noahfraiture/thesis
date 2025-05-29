package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Parse command-line flags
	serverAddr := flag.String("server", "127.0.0.1:8080", "server address (ip:port)")
	t := flag.String("type", "", "Type of message")
	hop := flag.Int("hop", 1, "")
	flag.Parse()

	var packet Packet
	switch *t {
	case "sleep":
		packet = &SleepPacket{
			Header: Header{PacketType: SLEEP, Hop: byte(*hop)},
			Time:   11,
		}
	case "spread":
		packet = &SpreadPacket{
			Header:   Header{PacketType: SPREAD, Hop: byte(*hop)},
			Quantity: 1,
		}
	case "exec":
		packet = &ExecPacket{
			Header:  Header{PacketType: EXEC, Hop: byte(*hop)},
			Command: "touch test",
		}
	case "quit":
		packet = &QuitPacket{
			Header: Header{PacketType: QUIT, Hop: byte(*hop)},
		}
	default:
		log.Fatal(*t)
	}

	packetBytes := packet.ToBytes()

	// Establish a TCP connection to the server
	conn, err := net.DialTimeout("tcp", *serverAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send the packet bytes
	_, err = conn.Write(packetBytes)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	fmt.Println("Message sent to", *serverAddr)
}
