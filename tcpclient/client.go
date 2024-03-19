package tcpclient

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func Myconn() {
	// 启动第一个TCP客户端
	go func() {
		conn, err := net.Dial("tcp", "114.55.244.179:31001")
		if err != nil {
			fmt.Println("Error connecting:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Client 1 connected to TCP server")

		// 构造要发送的十六进制数据
		hexData := "78781300010132373633483831303031503339362020000D0A"
		data, err := hex.DecodeString(hexData)
		if err != nil {
			fmt.Println("Error decoding hex data:", err)
			return
		}

		// 发送数据给服务器
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		// 接收并打印服务器发送的数据
		recvData := make([]byte, 1024)
		n, err := conn.Read(recvData)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			return
		}
		fmt.Printf("Client 1 received data from server: % X\n", recvData[:n])
		fmt.Println("conn 1 = ", conn)
		fmt.Println("conn 1 = ", conn.LocalAddr())
	}()
	time.Sleep(3 * time.Second)

	// 启动第一个TCP客户端
	go func() {
		conn, err := net.Dial("tcp", "114.55.244.179:31001")
		if err != nil {
			fmt.Println("Error connecting:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Client 2 connected to TCP server")

		// 构造要发送的十六进制数据
		hexData := "78781300010132373633483831303031503339362020000D0A"
		data, err := hex.DecodeString(hexData)
		if err != nil {
			fmt.Println("Error decoding hex data:", err)
			return
		}

		// 发送数据给服务器
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		// 接收并打印服务器发送的数据
		recvData := make([]byte, 1024)
		n, err := conn.Read(recvData)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			return
		}
		fmt.Printf("Client 2 received data from server: % X\n", recvData[:n])
		fmt.Println("conn 2 = ", conn)
		fmt.Println("conn 2 = ", conn.LocalAddr())
	}()
	fmt.Println("task start done")
	// 保持主goroutine不退出
	fmt.Scanln()

}
