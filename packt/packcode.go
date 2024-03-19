package packt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/sigurn/crc16"
)

type basicPackage interface {
	// Write(io.Writer) error
	// Unpack(io.Reader) error
}

const (
	radar = 0x03
)
const (
	slinkv1 = 0
	slinkv2 = 1
)

type ClientOptions struct {
	Servers []*url.URL
	proto   int
}

func NewClientOptions() *ClientOptions {
	o := &ClientOptions{
		Servers: nil,
		proto:   slinkv1,
	}
	return o
}

func NewClient(o *ClientOptions) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 1800})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		fmt.Println("remoteAddr = ", remoteAddr, "read length = ", n)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Println("")

		// fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		for i := 0; i < n; i++ {
			fmt.Printf("%X", data[i])
		}
		fmt.Println()
		a, _, _ := unpack(data[0:n])
		for _, v := range a {
			fmt.Println("")
			fmt.Printf("%+v", v)
			fmt.Println("")
			fmt.Printf("v.Header.msgid = %X\n", v.Header.msgid)
			if v.Header.msgid == 0xBC {
				var udpResp UdpBroadcastConfirmResponse
				err := binary.Read(bytes.NewReader(v.PayLoad), binary.LittleEndian, &udpResp)
				if err != nil {
					fmt.Println("解析 UDP 广播确认消息失败:", err)
					return
				}
				fmt.Printf("解析的 UDP 广播确认消息: %+v\n", udpResp)
				fmt.Println("sn = ", udpResp.Sn, " addr = ", udpResp.Addr, " port = ", udpResp.Port, " conntype = ", udpResp.ConnType)
			}
			// UdpBroadCastMessage
			if v.Header.msgid == 0xBA {
				var udpResp UdpBroadCastMessage
				err := binary.Read(bytes.NewReader(v.PayLoad), binary.LittleEndian, &udpResp)
				if err != nil {
					fmt.Println("-------> 解析 UDP 广播确认消息失败:", err)
					return
				}
				fmt.Printf("-------> 解析的 UDP 广播确认消息: %+v\n", udpResp)
				fmt.Println("-------> sn = ", udpResp.Sn, " Protocol = ", udpResp.Protocol)
				fmt.Printf("str sn = %s", udpResp.Sn)
			}

		}
		fmt.Println("")

		// _, err = listener.WriteToUDP([]byte("world"), remoteAddr)

		sendpack := sendBCAck()
		fmt.Println("0XBB sendpack = ")
		fmt.Printf("%X", sendpack)
		fmt.Println()
		// 1810
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteAddr.IP, 1810))
		_, err = listener.WriteToUDP(sendpack, addr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

type UdpBroadCastMessage struct {
	Protocol uint32   //协议版本
	Sn       [32]byte //
}

type UdpBroadcastConfirmResponse struct {
	Sn       [32]uint8
	Addr     [4]uint8
	Port     uint16
	ConnType uint16
}

type UdpBroadCastAckMessage struct {
	Protocol     uint32    //协议版本：数值越大版本越高
	Sn           [32]uint8 //设备唯一标识
	DeviceType   uint16    //设备类型
	DeviceStatus uint16    // 设备状态：等待连接、已连接
	Version      [64]byte  //版本信息
}

func (d *UdpBroadCastAckMessage) Encode() []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
		fmt.Errorf("UdpBroadCastMessage Encode error")
		return nil
	}
	return buf.Bytes()
}

func NewMavHeader(sourceid, destid, msgid uint8, len uint16) *Header {
	return &Header{
		magic:    magic,
		seq:      0,
		destid:   destid,
		sourceid: sourceid,
		msgid:    msgid,
		lenL:     uint8(len & 0x00FF),
		lenH:     uint8((len & 0xFF00) >> 8),
	}
}

// 0xBB
func sendBCAck() []byte {
	ssn := "go-util-sn"
	var arrsn [32]uint8

	for i, c := range ssn {
		arrsn[i] = uint8(c)
	}
	res := UdpBroadCastAckMessage{
		Protocol:     2,
		Sn:           arrsn,
		DeviceType:   radar,
		DeviceStatus: 1,
	}
	resBuff := res.Encode()
	header := NewMavHeader(0x03, 0x06, 0xbb, uint16(len(resBuff)))
	bcPacket := &Packet{
		Header:  *header,
		PayLoad: resBuff,
	}
	fmt.Printf("bcPacket = %+v", bcPacket)
	fmt.Println()
	buff, err := bcPacket.Encode()
	if err != nil {
		panic(err)
	}
	return buff
}

type Header struct {
	magic    uint8 //帧头
	lenL     uint8 //有效数据长度低8位
	lenH     uint8 //有效数据长度高8位
	seq      uint8 //帧序号
	destid   uint8 //接收者的设备ID
	sourceid uint8 //发送者的设备ID
	msgid    uint8 //消息ID
	ank      uint8 //是否需要应答
	checksum uint8 //包头校验值
}

type Packet struct {
	Header  Header
	PayLoad []byte //具体消息内容
	crc     uint16 //整帧数据校验
}

func (mavpacket *Packet) Encode() ([]byte, error) {
	mavpacket.Header.checkSum()
	buff := make([]byte, 0)
	buff = append(buff, mavpacket.Header.encode()...)
	buff = append(buff, mavpacket.PayLoad...)
	crc := make([]byte, crcLen)
	binary.BigEndian.PutUint16(crc, calcCrc(buff[crcStartLoc:]))
	for i := len(crc) - 1; i >= 0; i-- {
		buff = append(buff, crc[i])
	}
	return buff, nil
}

// checkSum 计算头部的校验和
func (m *Header) checkSum() {
	buff := []byte{m.magic, m.lenL, m.lenH, m.seq, m.destid, m.sourceid, m.msgid, m.ank}
	var sum uint32
	for _, v := range buff {
		sum += uint32(v)
	}
	m.checksum = uint8(sum & 0xFF)
}

// encode 头部编码
func (m *Header) encode() []byte {
	return []byte{m.magic, m.lenL, m.lenH, m.seq, m.destid, m.sourceid, m.msgid, m.ank, m.checksum}
}

const (
	checkSumLen = 8           //checksum 位校验长度
	headerLen   = 9           //header 的长度
	crcLen      = 2           //CRC 校验位长度
	crcStartLoc = 1           //CRC 开始位
	frameLoc    = 0           //帧头在协议包中字节位置
	magic       = uint8(0xFD) //帧头
	lenLoc      = 3           //包长度最大位置
)

func isMavLink(first byte) bool {
	return first == magic
}

// getPacketLen 获取包长度,默认buff第1、2位为包头的lenL, lenH
func getPacketLen(buff []byte) uint32 {
	if len(buff) < lenLoc {
		return 0
	}
	return uint32(uint16(buff[2])<<8|uint16(buff[1])) + headerLen + crcLen
}

// GetSourceID 获取源端设备ID
func (m *Packet) GetSourceID() uint8 {
	return m.Header.sourceid
}

// GetDestID 获取宿端设备ID
func (m *Packet) GetDestID() uint8 {
	return m.Header.destid
}

// GetMsgID 获取消息ID
func (m *Packet) GetMsgID() uint8 {
	return m.Header.msgid
}

// GetSeq 获取消息Seq号
func (m *Packet) GetSeq() uint8 {
	return m.Header.seq
}

// Len 包长度
func (m *Packet) Len() uint32 {
	return uint32(uint16(m.Header.lenH)<<8|uint16(m.Header.lenL)) + headerLen + crcLen
}

// getCrc 包CRC检验
func (m *Packet) getCrc(buff []byte) uint16 {
	if len(buff) != crcLen {
		return 0
	}
	return uint16(buff[1])<<8 | uint16(buff[0])
}

// payloadLen 载荷长度
func (m *Packet) payloadLen() uint32 {
	return uint32(uint16(m.Header.lenH)<<8 | uint16(m.Header.lenL))
}

// decode 头部解码
func (m *Header) decode(buff []byte) error {
	if len(buff) != headerLen {
		return errors.New("head len error")
	}
	m.magic = buff[0]
	m.lenL = buff[1]
	m.lenH = buff[2]
	m.seq = buff[3]
	m.destid = buff[4]
	m.sourceid = buff[5]
	m.msgid = buff[6]
	m.ank = buff[7]
	m.checksum = buff[8]
	return nil
}

// Decode 解码
func Decode(buff []byte) (*Packet, error) {
	if buff == nil {
		return nil, errors.New("buff empty")
	}
	if !isMavLink(buff[frameLoc]) {
		return nil, errors.New("not mavlink")
	}
	packet := &Packet{}
	err := packet.Header.decode(buff[:headerLen])
	if err != nil {
		return nil, err
	}
	if !isCheckSumEqual(buff[:checkSumLen], packet.Header.checksum) {
		return nil, errors.New("check sum error")
	}
	packet.crc = packet.getCrc(buff[(packet.Len() - crcLen):])
	if !isCheckCrcEqual(buff[crcStartLoc:(packet.Len()-crcLen)], packet.crc) {
		return nil, errors.New("check crc error")
	}
	packet.PayLoad = make([]byte, packet.payloadLen())
	copy(packet.PayLoad, buff[headerLen:(packet.Len()-crcLen)])
	return packet, nil
}

// isCheckSumEqual 校验校验和
func isCheckSumEqual(buff []byte, sum uint8) bool {
	if len(buff) != checkSumLen {
		return false
	}
	var check uint32
	for _, v := range buff {
		check += uint32(v)
	}
	return sum == uint8(check&0xFF)
}

// 处理接收到的buff
func unpack(buff []byte) ([]*Packet, []byte, error) {
	packets := make([]*Packet, 0)
	for len(buff) > 0 {
		//找到mavlink头开始计算包长度
		//fmt.Printf("receive buff len:%d, buff:% x", len(buff), buff)
		headLoc := 0
		for ; headLoc < len(buff); headLoc++ {
			if isMavLink(buff[headLoc]) {
				break
			}
			fmt.Printf("buff loc:%d not head, buff:%x", headLoc, buff[headLoc])
		}
		buff = buff[headLoc:]
		packetLen := getPacketLen(buff)
		//buff不足3字节无法取出包长度，或者包还没有接收完整，需要继续收包
		if packetLen <= 0 || len(buff) < int(packetLen) {
			break
		}
		packet, err := Decode(buff[:packetLen])
		if err != nil {
			fmt.Errorf("decode packet error %v", err)
			buff = buff[packetLen:]
			continue
		}
		//fmt.Printf(" mavlink %+v", packet)
		packets = append(packets, packet)
		buff = buff[packetLen:]
	}
	return packets, buff, nil
}

// isCheckCrcEqual 校验校验值
func isCheckCrcEqual(buff []byte, crc uint16) bool {
	if buff == nil {
		return false
	}
	return calcCrc(buff) == crc
}

// calcCrc 获取校验值
func calcCrc(buff []byte) uint16 {
	table := crc16.MakeTable(crc16.CRC16_MCRF4XX)
	return crc16.Checksum(buff, table)
}
