package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

// API GW WS max payload size is 128 KB.
// Make it less than that to account for JSON metadata.
// int32 = 4 bytes
const maxPayloadSize = 127000

func unixSocketServer(c net.Conn) {
	for {
		// Incoming buffer length
		prefix := make([]byte, 4)
		n, err := io.ReadFull(c, prefix)
		if err != nil || n < 4 {
			return
		}
		length := binary.BigEndian.Uint32(prefix)
		if length <= 0 {
			return
		}

		// Read 'length' bytes
		buf := make([]byte, length)
		n, err = io.ReadFull(c, buf)
		if uint32(n) == length {
			// Times four since one uint32 = 4 bytes
			if length*4 > maxPayloadSize {
				go LargeOutput(&buf, length)
			} else {
				go Output(string(buf))
			}
		}
	}
}

// InitUnixSocket ...
func InitUnixSocket() {
	l, err := net.Listen("unix", "/tmp/out")
	if err != nil {
		log.Fatal(err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go unixSocketServer(fd)
	}
}
