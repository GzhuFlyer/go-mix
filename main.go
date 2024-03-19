package main

import (
	"fmt"
	"mixtest/udbpt"
	"net"
	"strconv"
	"time"
)

func main() {
	// for i := 0; i < 1; i++ {
	// 	myqueue.QueueT4()
	// 	// myqueue.QueueT3()
	// // }
	// mqttTwo.M1()
	// mqttTwo.M2()
	// s := mqtt.NewAdaptorWithAuth("tcp://localhost:1883", "client123", "fzw", "12345678")
	// mqtt config : {tcp://10.240.34.35:30010 C2@1727584984128937984@c2app 1727584984128937984 e3bc38a4faa625d074664d572d810c1e}
	// s := mqtt.NewAdaptorWithAuth("tcp://10.240.34.35:30010", "C2@1727884703555842048@1334", "admin", "e3bc38a4faa625d074664d572d810c1e")

	// ret := s.Connect()
	// fmt.Println("connect ret = ", ret)
	// s.SetQoS(1)
	// mqtt.HandlerSub(s)
	// time.Sleep(500000 * time.Second)
	// mqtt.HandlerPub(s)
	// mqtt.TestCode()

	// gencode.Gen1()
	// gencode.Gen2()

	// udbpt.UDPT4()
	// websock.WS1()
	// go ffmpegt.StartPushVideo2Cloud()
	// time.Sleep(100000 * time.Minute)
	// mysn.Sn1()
	// mysn.Sn2()
	// mysn.WinSn()
	// option := packt.NewClientOptions()
	// packt.NewClient(option)
	// goroutine.Gorou1()
	// time.Sleep(1000 * time.Minute)
	// mymd5.Show()
	// mylock.Show()
	// myaws.Mytest1()
	// for {
	// 	websock.WS1()
	// }
	// tmp()
	// time.Sleep(99999 * time.Hour)
	// tcpclient.Myconn()
	// myaws.Mytest2()
	// data := []byte("这里是你的25个字节数据，可以是任意内容")

	// // 计算SHA-256哈希值
	// hash := sha256.Sum256(data)

	// // 取哈希值的前16个字节作为结果
	// result := hash[:16]
	// fmt.Println(result)

	// addr, _ := udbpt.GetClientIp()
	// fmt.Println("addr = ", addr)
	udbpt.UDP5()
}

func tmp2() {
	i := 12
	is := fmt.Sprintf("%04d", i)
	fmt.Println("is = ", is)

	j := "1065"
	js, err := strconv.Atoi(j)
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Println("js = ", js)
}

func tmp() {
	tick := time.NewTicker(5 * time.Second)
	addr := "192.168.1.21:1800"
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			fmt.Println("Run connect NFS4000")
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println("再链接失败, err:", err)
				continue
			}
			fmt.Println("conn = ", conn)
		}
	}
}
