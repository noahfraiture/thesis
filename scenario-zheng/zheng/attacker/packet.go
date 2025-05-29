package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type PACKET_TYPE byte

const (
	_      PACKET_TYPE = iota
	SLEEP              // parent -> child
	SPREAD             // parent -> child
	EXEC               // parent -> child
	QUIT               // parent -> child
	ASK                // child -> parent
	EMPTY              // empty packet
)

type Packet interface {
	ToBytes() []byte
}

type Header struct {
	PacketType PACKET_TYPE
	Source     string
	Hop        byte
}

func (p *Header) ToBytes() []byte {
	bytes := []byte{byte(p.PacketType)}
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(len(p.Source)))
	bytes = append(bytes, []byte(p.Source)...)
	bytes = append(bytes, p.Hop)
	return bytes
}

type SleepPacket struct {
	Header
	Time int
}

func (p *SleepPacket) ToBytes() []byte {
	bytes := p.Header.ToBytes()
	return binary.BigEndian.AppendUint32(bytes, uint32(p.Time))
}

type SpreadPacket struct {
	Header
	Quantity int
}

func (p *SpreadPacket) ToBytes() []byte {
	bytes := p.Header.ToBytes()
	return binary.BigEndian.AppendUint32(bytes, uint32(p.Quantity))
}

type ExecPacket struct {
	Header
	Command string
}

func (p *ExecPacket) ToBytes() []byte {
	bytes := p.Header.ToBytes()
	return append(bytes, []byte(p.Command)...)
}

type QuitPacket struct {
	Header
}

func (p *QuitPacket) ToBytes() []byte {
	return p.Header.ToBytes()
}

type AskPacket struct {
	Header
}

func (p *AskPacket) ToBytes() []byte {
	return p.Header.ToBytes()
}

type EmptyPacket struct {
	Header
}

func (p *EmptyPacket) ToBytes() []byte {
	return p.Header.ToBytes()
}

func Parse(data []byte) (Packet, error) {
	reader := bytes.NewReader(data)

	// Read packet type
	packetTypeByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	packetType := PACKET_TYPE(packetTypeByte)

	// Read source length
	srcLenBuf := make([]byte, 2)
	_, err = io.ReadFull(reader, srcLenBuf)
	if err != nil {
		return nil, err
	}
	srcLen := binary.BigEndian.Uint16(srcLenBuf)

	// Read source
	sourceBytes := make([]byte, srcLen)
	_, err = io.ReadFull(reader, sourceBytes)
	if err != nil {
		return nil, err
	}
	source := string(sourceBytes)

	// Read hop
	hop, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	// Create header
	header := Header{
		PacketType: packetType,
		Source:     source,
		Hop:        hop,
	}

	// Parse based on packet type
	switch packetType {
	case SLEEP:
		timeBuf := make([]byte, 4)
		_, err = io.ReadFull(reader, timeBuf)
		if err != nil {
			return nil, err
		}
		sleepTime := int(binary.BigEndian.Uint32(timeBuf))
		return &SleepPacket{Header: header, Time: sleepTime}, nil

	case SPREAD:
		quantityBuf := make([]byte, 4)
		_, err = io.ReadFull(reader, quantityBuf)
		if err != nil {
			return nil, err
		}
		quantity := int(binary.BigEndian.Uint32(quantityBuf))
		return &SpreadPacket{Header: header, Quantity: quantity}, nil

	case EXEC:
		// Read remaining bytes as the command string
		commandBytes := make([]byte, reader.Len())
		_, err = io.ReadFull(reader, commandBytes)
		if err != nil {
			return nil, err
		}
		return &ExecPacket{Header: header, Command: string(commandBytes)}, nil

	case EMPTY:
		return &EmptyPacket{Header: header}, nil

	case QUIT:
		return &QuitPacket{Header: header}, nil

	case ASK:
		return &AskPacket{Header: header}, nil

	default:
		return nil, errors.New("unknown packet type")
	}
}
