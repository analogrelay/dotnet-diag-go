package client

import (
	"encoding/binary"
	"io"
	"unsafe"
)

const magic_v1 = "DOTNET_IPC_V1"
const magic_size = 14
const header_size = 20

func getHostEndian() binary.ByteOrder {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		return binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

type IpcHeader struct {
	Magic      [14]uint8
	Size       uint16
	CommandSet uint8
	CommandId  uint8
	Reserved   uint16
}

type IpcMessage struct {
	Header  IpcHeader
	Payload []byte
}

func ReadMessage(r io.Reader) (IpcMessage, error) {
	// Read the header, which is a fixed size
	header, err := readHeader(r)
	if err != nil {
		return IpcMessage{}, err
	}

	// Calculate the payload size, it's the header's size value minus the header size
	payloadSize := header.Size - header_size
	buf := make([]byte, 0, payloadSize)

	// Read the payload
	if _, err := io.ReadFull(r, buf); err != nil {
		return IpcMessage{}, err
	}

	return IpcMessage{header, buf}, nil
}

func WriteMessage(message *IpcMessage, w io.Writer) error {
	// Compute the correct size and update the header value
	message.Header.Size = uint16(header_size + len(message.Payload))

	// Write the header
	if err := writeHeader(&message.Header, w); err != nil {
		return err
	}

	// Write the payload
	if _, err := w.Write(message.Payload); err != nil {
		return err
	}
	return nil
}

func readHeader(r io.Reader) (IpcHeader, error) {
	var header IpcHeader
	err := binary.Read(r, getHostEndian(), &header)
	if err != nil {
		return IpcHeader{}, err
	} else {
		return header, nil
	}
}

func writeHeader(header *IpcHeader, w io.Writer) error {
	return binary.Write(w, getHostEndian(), *header)
}
