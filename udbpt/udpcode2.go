package udbpt

import (
	"log"
	"net"
)

func UDPT2() {
	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", "10.240.34.12:9966")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	buf[2] |= 0x80 // Set QR bit

	pc.WriteTo(buf, addr)
}
