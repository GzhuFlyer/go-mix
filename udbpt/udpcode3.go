package udbpt

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func sendPackage() []byte {
	startBits := uint32(0x8A808988)     //起始位：4个Byte
	protocolVersion := uint32(9002)     //协议号：4个Byte
	pkgLen := uint32(48)                //包长： 4个Byte -> 28 + 4 + 4 + 4
	command := uint32(0x0D)             //命令字： 4个Byte
	timestamp := time.Now().UnixMilli() //时间戳： 8个Byte
	//信息内容：N个Byte = 4 + 8 + 4 + 4 + 4 + 4 = 28
	msgCode := uint32(0x01)
	msgTimeStamp := time.Now().UnixMilli()
	msgCmd := uint32(0x01)
	msgLevel := uint32(20)
	msgType := uint32(0x0)
	msgResver := uint32(0x01)

	msgSeq := uint32(0x1) //信息序列号 4个 Byte
	// checkSum := uint32(0x32) //错误校验: 4个 Byte msgTimeStamp timdstamp
	checkSum := uint32(int64(pkgLen+command+msgSeq+msgCode+msgCmd+msgLevel+msgType+msgResver)+msgTimeStamp+timestamp) & 0xFFFFFFFF
	stopBits := uint32(0x8B8A8089) //停止位 4个 Byte

	buffer := make([]byte, 65)

	binary.LittleEndian.PutUint32(buffer[0:4], startBits)
	binary.LittleEndian.PutUint32(buffer[4:8], protocolVersion)
	binary.LittleEndian.PutUint32(buffer[8:12], pkgLen)
	binary.LittleEndian.PutUint32(buffer[12:16], command)
	binary.LittleEndian.PutUint64(buffer[16:24], uint64(timestamp))
	binary.LittleEndian.PutUint32(buffer[24:28], msgCode)
	binary.LittleEndian.PutUint64(buffer[28:36], uint64(msgTimeStamp))
	binary.LittleEndian.PutUint32(buffer[36:40], msgCmd)
	binary.LittleEndian.PutUint32(buffer[40:44], msgLevel)
	binary.LittleEndian.PutUint32(buffer[44:48], msgType)
	binary.LittleEndian.PutUint32(buffer[48:52], msgResver)
	binary.LittleEndian.PutUint32(buffer[52:56], msgSeq)
	binary.LittleEndian.PutUint32(buffer[56:60], checkSum)
	binary.LittleEndian.PutUint32(buffer[60:64], stopBits)

	fmt.Println(buffer)
	return buffer
}

func UDPT3() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		// IP:   net.IPv4(192, 168, 1, 4),
		IP:   net.IPv4(192, 168, 1, 4),
		Port: 9966,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	defer socket.Close()
	// sendData := []byte("Hello server")
	sendData := []byte{0x88, 0x89, 0x80, 0x8A, 0x2A, 0x23, 0x0, 0x0, 0x30, 0x0, 0x0, 0x0, 0xD, 0x0, 0x0, 0x0, 0x37, 0xAE, 0x21, 0xF3, 0x8C, 0x1, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x37, 0xAE, 0x21, 0xF3, 0x8C, 0x1, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xC3, 0x5C, 0x43, 0xE6, 0x89, 0x80, 0x8A, 0x8B, 0x0}

	// sendData := sendPackage()

	fmt.Printf("%c", '{')
	for _, v := range sendData {
		fmt.Printf("0x%X,", v)
	}
	fmt.Printf("%c\n", '}')
	for k, v := range sendData {
		if k%8 == 0 {
			fmt.Println("")
		}
		fmt.Printf("[%d:%2X] ", k, v)
	}
	fmt.Println("")
	_, err = socket.Write(sendData) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	data := make([]byte, 4096)

	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
