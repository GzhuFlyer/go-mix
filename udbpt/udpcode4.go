package udbpt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// 包通用字段定义
const (
	startBits   uint32 = 0x8A808988
	protocolVer uint32 = 9002
	stopFlag    uint32 = 0x8B8A8089
	sunCode     uint32 = 0x01
	reserved    uint32 = 0x00
	extraLength uint32 = 20
)
const (
	sysCode uint32 = 0x0
)

// 自增序列号
var (
	sequence = uint32(0)
	mutex    sync.Mutex
)

func incrementSeq() uint32 {
	mutex.Lock()
	defer mutex.Unlock()
	sequence++
	return sequence
}

// 包命令字长度
const (
	//旋转
	directionCmd uint32 = 0x0D
	zoomCmd      uint32 = 0x09
)

// 光电编号
const (
	code uint32 = 0
)

// 控制方向命令字内容定义
const (
	stop      uint32 = 0
	left      uint32 = 1
	right     uint32 = 2
	up        uint32 = 3
	down      uint32 = 4
	lup       uint32 = 5
	ldown     uint32 = 6
	rup       uint32 = 7
	rdown     uint32 = 8
	keepMov   uint32 = 0
	singleMov uint32 = 1
)

// /================================================================
// Zoom镜头控制命令字内容定义
const (
	zoomStop     uint32 = 0  //镜头(包含zoom和focus)停止运动
	jumpPhy      uint32 = 1  //跳转到指定物理焦距
	jumpMultiple uint32 = 2  //跳转到指定倍数
	jumpZoFo     uint32 = 3  //跳转到指定zoom和focus位置
	zoomkeepOut  uint32 = 4  //镜头持续推远(直到收到停止指令)
	zoomkeepIn   uint32 = 5  //镜头持续拉近(直到收到停止指令)
	zoomOut      uint32 = 6  //镜头推远(点动)
	zoomIn       uint32 = 7  //镜头拉近(点动)
	focusKeepOut uint32 = 8  //聚焦持续推远
	focusKeepIn  uint32 = 9  //聚焦持续拉近
	focusOut     uint32 = 10 // 0x0A: 聚焦推远(点动)
	focusIn      uint32 = 11 // 0x0B: 聚焦拉近(点动)
)

// 物理焦距/倍数
const (
	phyFocus     uint32 = 0x01 // 0x01时表示镜头的物理焦距, 单位mm 如150mm
	multiple     uint32 = 0x02 // 命令为0x02时表示镜头倍数，如42倍
	zoomReserved uint32 = 0x03 //控制命令为0x03时此字段作为保留不使用
	// 控制指令为0x04 ~ 0x0B, 该字段用来表示速度，取值范围0-254
)

// ================================================================

// 运动速度等级 0-254

func timestamp() uint64 {
	return uint64(time.Now().Unix())
}

func formDirectionlPayload(cmd, speed, movement uint32) ([]byte, uint32, uint64) {
	fmt.Println("cmd  = ", cmd)
	fmt.Println("speed = ", speed)
	fmt.Println("movement = ", movement)
	len := 28
	payload := make([]byte, len)
	ts := timestamp()
	binary.LittleEndian.PutUint32(payload[0:4], sunCode)
	binary.LittleEndian.PutUint64(payload[4:12], ts)
	binary.LittleEndian.PutUint32(payload[12:16], cmd)
	binary.LittleEndian.PutUint32(payload[16:20], speed)
	binary.LittleEndian.PutUint32(payload[20:24], movement)
	binary.LittleEndian.PutUint32(payload[24:28], reserved)
	sum := uint64(sunCode+cmd+speed+movement+reserved) + ts
	return payload, uint32(len), sum
}

func formZoomlPayload(cmd, focus, zoomsite, focusSite uint32) ([]byte, uint32, uint64) {
	fmt.Println("cmd  = ", cmd)
	fmt.Println("focus = ", focus)
	fmt.Println("zoomsite = ", zoomsite)
	fmt.Println("focusSite = ", focusSite)
	len := 36
	payload := make([]byte, len)
	ts := timestamp()
	binary.LittleEndian.PutUint32(payload[0:4], sunCode)

	binary.LittleEndian.PutUint32(payload[4:8], sysCode)
	binary.LittleEndian.PutUint64(payload[8:16], ts)
	binary.LittleEndian.PutUint32(payload[16:20], cmd)
	binary.LittleEndian.PutUint32(payload[20:24], focus)
	binary.LittleEndian.PutUint32(payload[24:28], zoomsite)
	binary.LittleEndian.PutUint32(payload[28:32], focusSite)
	binary.LittleEndian.PutUint32(payload[32:36], reserved)
	sum := uint64(sunCode+sysCode+cmd+focus+zoomsite+focusSite+reserved) + ts
	return payload, uint32(len), sum
}

// “包长度”到“信息序列号”
func sendPkg(payload []byte, len uint32, sum uint64, cmd uint32) []byte {
	fmt.Println("len = ", len)
	fmt.Println("cmd = ", cmd)
	buf := new(bytes.Buffer)
	pkgLen := len + extraLength
	binary.Write(buf, binary.LittleEndian, startBits)
	binary.Write(buf, binary.LittleEndian, protocolVer)
	binary.Write(buf, binary.LittleEndian, pkgLen)
	binary.Write(buf, binary.LittleEndian, cmd)
	binary.Write(buf, binary.LittleEndian, timestamp())
	buf.Write(payload)
	sequence++
	binary.Write(buf, binary.LittleEndian, incrementSeq())
	// checksum := uint32(uint64(packetLength+command+sequence+msgCode+msgCmd+msgLevel+msgType+msgResver)+msgTimeStamp+timestamp) & 0xFFFFFFFF
	checksum := uint32(sum + uint64(pkgLen))
	binary.Write(buf, binary.LittleEndian, checksum)
	binary.Write(buf, binary.LittleEndian, stopFlag)
	fmt.Println("pkgLen = ", pkgLen)
	return buf.Bytes()
}

func UDPT4() {
	// 定义服务器地址和端口号
	serverAddr := "192.168.1.4:9966"
	// serverAddr := "127.0.0.1:9966"

	// 解析服务器地址
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("解析服务器地址出错：", err)
		os.Exit(1)
	}

	// 建立UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("建立UDP连接出错：", err)
		os.Exit(1)
	}
	defer conn.Close()
	payload, len, sum := formDirectionlPayload(left, 20, keepMov)
	buf := sendPkg(payload, len, sum, directionCmd)
	// payload, len, sum := formZoomlPayload(0x05, 10, 0, 0)
	// buf := sendPkg(payload, len, sum, zoomCmd)
	// 发送数据包
	// fmt.Println("buf = ", buf)
	// fmt.Printf("%02x", buf)
	//left
	// 88 89 80 8a 2a 23 00 00 30 00 00 00 0d 00 00 00 22 59 a2 65 00 00 00 00 01 00 00 00 22 59 a2 65 00 00 00 00 01 00 00 00 7b 00 00 00 00 00 00 00 01 00 00 00 02 00 00 00 d0 59 a2 65 89 80 8a 8b
	// cmd  =  13
	// child =  1
	// speed =  0
	// movement =  0
	// len =  28
	// cmd =  13
	// 88 89 80 8a 2a 23 00 00 30 00 00 00 0d 00 00 00 e0 e5 a4 65 00 00 00 00 01 00 00 00 e0 e5 a4 65 00 00 00 00 01 00 00 00 00 00 00 00 00 00 00 00 01 00 00 00 02 00 00 00 13 e6 a4 65 89 80 8a 8b
	//stop
	// 88 89 80 8a 2a 23 00 00 30 00 00 00 0d 00 00 00 2f 59 a2 65 00 00 00 00 01 00 00 00 2f 59 a2 65 00 00 00 00 00 00 00 00 7b 00 00 00 00 00 00 00 01 00 00 00 02 00 00 00 dc 59 a2 65 89 80 8a 8b
	for _, v := range buf {
		fmt.Printf("%02x ", v)
	}
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("发送数据出错：", err)
		os.Exit(1)
	}

	fmt.Println("数据包发送成功！")

}

//88 89 80 8a 2a 23 00 00 38 00 00 00 09 00 00 00 fd 1d a5 65 00 00 00 00 01 00 00 00 00 00 00 00 fd 1d a5 65 00 00 00 00 05 00 00 00 0a 00 00 00 00 00 00 00 00 00 00 00 01 00 00 00 0c 00 00 00 46 1e a5 65 89 80 8a 8b
//88 89 80 8a 2a 23 00 00 38 00 00 00 09 00 00 00 35 1e a5 65 00 00 00 00 01 00 00 00 00 00 00 00 35 1e a5 65 00 00 00 00 05 00 00 00 0a 00 00 00 00 00 00 00 00 00 00 00 01 00 00 00 0e 00 00 00 7e 1e a5 65 89 80 8a 8b

//88 89 80 8a 2a 23 00 00 38 00 00 00 09 00 00 00 89 1d a5 65 00 00 00 00 01 00 00 00 00 00 00 00 89 1d a5 65 00 00 00 00 05 00 00 00 0a 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02 00 00 00 d1 1d a5 65 89 80 8a 8b
//88 89 80 8a 2a 23 00 00 38 00 00 00 09 00 00 00 dc 1d a5 65 00 00 00 00 01 00 00 00 00 00 00 00 dc 1d a5 65 00 00 00 00 05 00 00 00 0a 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02 00 00 00 24 1e a5 65 89 80 8a 8b
