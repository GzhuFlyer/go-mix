package websock

import (
	"fmt"
	"log"
	"mixtest/websock/deppb/deppb"
	"net/url"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// 定义消息结构体
type Message struct {
	MessageType int32  `protobuf:"varint,1,opt,name=message_type,json=messageType" json:"message_type,omitempty"`
	Message     string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

const (
	systemHeart                        = 3000001
	ClientMsgIDTracerDetect            = 1000008
	ClientMsgIdAgxCalibrationData      = 1000037 // mavlink.agxCalibrationReport
	ClientMsgIdAgxHeartData            = 1000031 // mavlink.AgxHeartData
	ClientMsgIDSFLHeartBeat            = 1000024 // mavlink.SflHeartMsgReport		302
	ClientMsgIDSFLDetect               = 1000025 // mavlink.SflDetectMsgReport		300
	ClientMsgIDTracerDroneRemoteDetect = 1000040 // mavlink.ClientMsgIDTracerDroneRemoteDetect		234
)

func WS1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover() err = ", err)
			time.Sleep(1 * time.Second)
			WS1()
		}

	}()
	// 设置 WebSocket 服务地址
	// whost := "172.20.208.1:9907"
	// whost := "127.0.0.1:9907"
	whost := "10.10.49.234:9907"
	wPath := "/ws/v1/messagebox"
	u := url.URL{Scheme: "ws", Host: whost, Path: wPath}

	// 连接 WebSocket 服务
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Fatal("连接 WebSocket 服务失败：", err)
		fmt.Println("连接 WebSocket 服务失败：", err)
	}
	defer conn.Close()

	// 循环读取 WebSocket 数据
	for {
		// 读取 WebSocket 消息
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败：", err)
			time.Sleep(1 * time.Second)
			break
		}
		// fmt.Println("p = ", p)
		// 解码消息
		message := deppb.ClientReport{}
		err = proto.Unmarshal(p, &message)
		if err != nil {
			log.Println("解码消息失败：", err)
			continue
		}
		// fmt.Printf("message = %s\n", message)
		// fmt.Printf("message type = %d\n", message.MsgType)

		if systemHeart == message.MsgType {
			msg2 := deppb.SflDetectInfo{}
			// err = proto.Unmarshal(message.Data, &msg2)
			// if err != nil {
			// 	log.Println("解码消息失败：", err)
			// 	continue
			// }
			// fmt.Println(time.Now().Format("2006-01-02 15:04:05.00"))
			fmt.Println("systemHeart msg2 = ", msg2)

		}

		if ClientMsgIDSFLHeartBeat == message.MsgType {
			msg2 := deppb.SflDetectInfo{}
			err = proto.Unmarshal(message.Data, &msg2)
			if err != nil {
				log.Println("解码消息失败：", err)
				continue
			}
			fmt.Println(time.Now().Format("2006-01-02 15:04:05.00"))
			fmt.Println("msg2 = ", msg2)

		}

		if ClientMsgIDTracerDetect == message.MsgType {
			msg2 := deppb.TracerDetectInfo{}
			err = proto.Unmarshal(message.Data, &msg2)
			if err != nil {
				log.Println("解码消息失败：", err)
				continue
			}
			fmt.Println(time.Now().Format("2006-01-02 15:04:05.00"))
			fmt.Println("TracerDetectInfo msg2 = ", msg2)
		}

		if ClientMsgIDTracerDroneRemoteDetect == message.MsgType {
			msg2 := deppb.TracerDroneIdRemoteIdDetectInfo{}
			err = proto.Unmarshal(message.Data, &msg2)
			if err != nil {
				log.Println("解码消息失败：", err)
				continue
			}

			// msg3 := deppb.EquipmentMessageBoxEntity{}
			// err = proto.Unmarshal(msg2.Header, &msg3)
			// if err != nil {
			// 	log.Println("解码消息失败：", err)
			// 	continue
			// }

			fmt.Println(time.Now().Format("2006-01-02 15:04:05.00"))
			fmt.Println("TracerDroneIdRemoteIdDetectInfo header = ", msg2.Header)
			fmt.Println("TracerDroneIdRemoteIdDetectInfo Data.Sn = ", msg2.Data.Sn)
			fmt.Println("TracerDroneIdRemoteIdDetectInfo Data.Info = ", msg2.Data.Info)
		}

		if ClientMsgIdAgxHeartData == message.MsgType {
			fmt.Println("message = ", message.MsgType)
			msg2 := deppb.AgxHeartBeatInfo{}
			err = proto.Unmarshal(message.Data, &msg2)
			if err != nil {
				log.Println("解码消息失败：", err)
				continue
			}
			fmt.Println("msg2 = ", msg2)
		}

		if ClientMsgIdAgxCalibrationData == message.MsgType {
			now := time.Now()

			// 打印年月日时分秒
			fmt.Println(now.Format("2006-01-02 15:04:05"))
			fmt.Println("--------------------------------------")
			fmt.Println("message = ", message.MsgType)

			msg3 := deppb.AgxPTZCalibrationInfo{}
			err = proto.Unmarshal(message.Data, &msg3)
			if err != nil {
				log.Println("解码消息失败：", err)
				continue
			}
			fmt.Println("msg3 = ")
			fmt.Println(msg3)
			fmt.Println("header = ", msg3.Header)
			fmt.Println("data = ", msg3.Data)

			// 			{{{} [] [] 0xc000162670} 0 [] name:"AGX-12345678"  sn:"AGX-12345678"  equipType:12  msgType:239 flag:1}
			// header =  name:"AGX-12345678"  sn:"AGX-12345678"  equipType:12  msgType:239
			// data =  flag:1
			fmt.Println("--------------------------------")
			fmt.Println("-------------+-------------------------")
		}
	}
}
