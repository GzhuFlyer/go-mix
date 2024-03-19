package udbpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	file, err := os.OpenFile("a.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("fail to open log file ", err)
	}
	log.SetOutput(file)
}

func logger(v ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	log.Printf("%s:%d %s: %s", file, line, funcName, v)
	// log.Printf(file, ":", line, ":", funcName, ":")
	// log.Println(v...)
	// log.Printf("%s", v)
}

type NSF4000Msg struct {
	IsOnline    int     `json:"isOnline"`    //1:在线
	IsWorking   int     `json:"isWorking"`   //工作状态： 1：工作中，2：未工作
	GpsStatus   int     `json:"gpsStatus"`   //定位状态 1:已定位  2：未定位 sTsSyncSta B1:接收机定位状态
	Ephemeris   int     `json:"ephemeris"`   //星历    1:以获取  2:未获取 //4个系统任意b2为0
	TimeSync    int     `json:"timeSync"`    //时间同步状态  1:已同步  2:未同步 晶振 IOcxoSta >=2
	Longititude float64 `json:"longititude"` //经度
	Latitude    float64 `json:"latitude"`    //纬度
	Height      float64 `json:"height"`
	Enable      bool    `json:"enable"`   //功能开关 true:开启   false:关闭 默认值false
	WorkMode    int32   `json:"workMode"` //工作模式 1：主动防御 2：区域拒止 3：定向驱离 默认值1
	Radius      int32   `json:"radius"`   //防御区半径 默认值500
	Angle       int32   `json:"angle"`    //诱导角度  0/90/180/270   0正北，90正东  1  向上  -1向下  默认值0
	OpsTime     int64   `json:"opsTime"`  //记录操作时长，时间单位为秒(enable_true_time - enable_false_time ) 默认值0
}

const (
	header = "FF"
)

type updateIPPort619 struct {
	SKey  string `json:"sKey"`
	SIP   string `json:"sIP"`
	IPort int    `json:"iPort"`
}

func GetClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				strIP := ipnet.IP.String()
				logger("strIP = ", strIP)
				if strings.Contains(strIP, "192.168.10") {
					return strIP, nil
				}
			}

		}
	}

	return "", errors.New("Can not find the client ip address!")

}
func updateIPPort619Fun(skey, sip string, port int) string {
	if skey == "" && sip == "" && port == 0 {
		skey = "a57502fcdc4e7412"
		sip = "192.168.10.8"
		port = 9098
	}

	body := &updateIPPort619{
		SKey:  skey,
		SIP:   sip,
		IPort: port,
	}
	return FormPkg(body, "619")
}

func FormPkg(v any, c string) string {
	data, err := json.Marshal(v)
	if err != nil {
		logger("data")
	}
	// pkgLen := strconv.Itoa(len(data))
	// pl := len(pkgLen)
	// var fix string
	// for i := pl; i < 4; i++ {
	// 	fix = fix + "0"
	// }
	// pkgLen = fix + pkgLen
	pkgLen := fmt.Sprintf("%04d", len(data))

	logger("pkgLen")
	pkg := header + pkgLen + c + string(data)
	return pkg
}

// 发射功率
type update603 struct {
	SKey      string `json:"sKey"`
	IType     int    `json:"iType"`
	FAttenGPS int    `json:"fAttenGPS"`
	FAttenBDS int    `json:"fAttenBDS"`
	FAttenGLO int    `json:"fAttenGLO"`
	FAttenGAL int    `json:"fAttenGAL"`
}
type update615 struct {
	SKey   string `json:"sKey"`
	ICycle int    `json:"iCycle"`
}

func update615Fun() string {

	body := &update615{
		SKey:   "a57502fcdc4e7412",
		ICycle: 2,
	}
	return FormPkg(body, "615")
}

func UDP5() {
	go sendudp()
	rcvudp()
	// time.Sleep(3 * time.Second)

	time.Sleep(9999 * time.Hour)
}

// FF0063619{"sKey":"a57502fcdc4e7412","sIP":"192.168.10.123","iPort":9098}
// FF0063619{"sKey":"a57502fcdc4e7412","sIP":"192.168.10.123","iPort":9098}
func sendudp() {
	svraddr := "192.168.10.238:9099"
	udpaddr, err := net.ResolveUDPAddr("udp", svraddr)
	if err != nil {
		logger("err = ", err)
	}

	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		logger("DialUDP err = ", err)
	}
	defer conn.Close()

	// paload := "FF0062619{\"sKey\":\"a57502fcdc4e7412\",\"sIP\":\"192.168.10.88\",\"iPort\":9098}"

	// localIp, err := GetClientIp()
	// if err != nil {
	// 	logger("err = ", err)
	// 	return
	// }
	// payload := updateIPPort619Fun("a57502fcdc4e7412", localIp, 9098)
	payload := update615Fun()

	// payload := update603Fun()

	n, err := conn.Write([]byte(payload))
	if err != nil {
		logger("write err = ", err)
	}
	logger("n = ", n)
	logger("paload = ", payload)

	readbuf := make([]byte, 1024*1024)
	for {
		n, addr, err := conn.ReadFromUDP(readbuf)
		if err != nil {
			logger("err = ", err)
			continue
		}
		logger("sendudp()  addr = ", addr)
		logger("sendudp()  recv , len = ", n)
		logger("sendudp()  read buf = ", readbuf[:n])
	}
}

// "sDevNum":"4758340280353A63"
func rcvudp() {
	// addr, err := net.ResolveUDPAddr("udp", "192.168.10.88:9098")
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9098")
	if err != nil {
		logger("resolveudp err = ", err)
	}
	logger("add = ", addr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger("err = ", err)
	}
	defer conn.Close()
	readbuf := make([]byte, 1024*1024)
	heartData := NsfGXWHeart{}
	for {
		n, addr, err := conn.ReadFromUDP(readbuf)
		if err != nil {
			logger("err = ", err)
			continue
		}
		logger("addr = ", addr)
		logger("recv , len = ", n)
		logger("read buf = ", string(readbuf[:n]))
		dataLen, _ := strconv.Atoi(string(readbuf[2:6]))

		logger("data buf = ", string(readbuf[9:dataLen+9]))
		logger("dataLen = ", dataLen)
		err = json.Unmarshal(readbuf[9:dataLen+9], &heartData)
		if err != nil {
			logger("err umar json = ", err)
		}
		// logger("heart data = ", heartData)
		var s NSF4000Msg
		s.IsOnline = 1
		s.IsWorking = 1
		if heartData.IOcxoSta >= 2 {
			// s.GpsStatus = 1
			// s.Ephemeris = 1
			s.TimeSync = 1
		}
		// tmp1 := string([]byte(heartData.STsSyncSta)[1:2])

		if getByteNum(heartData.STsSyncSta, 1) == 0 {
			s.GpsStatus = 1
		}
		if getByteNum(heartData.SWorkStaGPS, 2) == 0 ||
			getByteNum(heartData.SWorkStaBDS, 2) == 0 ||
			getByteNum(heartData.SWorkStaGLO, 2) == 0 ||
			getByteNum(heartData.SWorkStaGAL, 2) == 0 {
			s.Ephemeris = 1

		}

		s.Longititude = heartData.DbFixLon
		s.Latitude = heartData.DbFixLat
		s.Height = heartData.DbFixLat
		logger("NSF4000Msg = ", s)
	}
}

func getByteNum(s string, site int) int {
	num, err := strconv.Atoi(string([]byte(s)[site : site+1]))
	if err == nil {
		return num
	}
	logger("err = ", err)
	return 1
}

type NsfGXWHeart struct {
	ISysSta        int     `json:"iSysSta"`        // 1. 系统工作状态
	STime          string  `json:"sTime"`          // 2. 时间
	SDevNum        string  `json:"sDevNum"`        // 3. 设备编号
	IDevType       int     `json:"iDevType"`       // 4. 设备类型
	ISysRunTime    int     `json:"iSysRunTime"`    // 5. 系统运行时长
	FEnvTemp       float32 `json:"fEnvTemp"`       // 6. 环境温度
	STsSyncSta     string  `json:"sTsSyncSta"`     // 7. 授时同步状态字
	DbFixLon       float64 `json:"dbFixLon"`       // 8. 当前定位位置经度
	DbFixLat       float64 `json:"dbFixLat"`       // 9. 当前定位位置纬度
	DbFixAlt       float64 `json:"dbFixAlt"`       // 10. 当前定位位置高度
	IOcxoSta       int     `json:"iOcxoSta"`       // 11. 晶振状态
	FTimeAccur     float32 `json:"fTimeAccur"`     // 12. 时间精度
	SWorkStaGPS    string  `json:"sWorkStaGPS"`    // 13. GPS工作状态字
	SWorkStaBDS    string  `json:"sWorkStaBDS"`    // 14. BDS工作状态字
	SWorkStaGLO    string  `json:"sWorkStaGLO"`    // 15. GLONASS工作状态字
	SWorkStaGAL    string  `json:"sWorkStaGAL"`    // 16. GALILEO工作状态字
	ISatNumGPS     int     `json:"iSatNumGPS"`     // 17. GPS转发卫星颗数
	ISatNumBDS     int     `json:"iSatNumBDS"`     // 18. BDS转发卫星颗数
	ISatNumGLO     int     `json:"iSatNumGLO"`     // 19. GLONASS转发卫星颗数
	ISatNumGAL     int     `json:"iSatNumGAL"`     // 20. GALILEO转发卫星颗数
	ISwitchGPS     int     `json:"iSwitchGPS"`     // 21. GPS发射开关
	ISwitchBDS     int     `json:"iSwitchBDS"`     // 22. BDS发射开关
	ISwitchGLO     int     `json:"iSwitchGLO"`     // 23. GLONASS发射开关
	ISwitchGAL     int     `json:"iSwitchGAL"`     // 24. GALILEO发射开关
	IPASwitch      int     `json:"iPASwitch"`      // 25. 功放开关
	DbSimuLon      float64 `json:"dbSimuLon"`      // 26. 当前诱骗位置经度
	DbSimuLat      float64 `json:"dbSimuLat"`      // 27. 当前诱骗位置高度
	DbSimuAlt      float64 `json:"dbSimuAlt"`      // 28. 当前诱骗位置高度
	FAttenGPS      float32 `json:"fAttenGPS"`      // 29. GPS功率衰减值
	FAttenBDS      float32 `json:"fAttenBDS"`      // 30. BDS功率衰减值
	FAttenGLO      float32 `json:"fAttenGLO"`      // 31. GLONASS功率衰减值
	FAttenGAL      float32 `json:"fAttenGAL"`      // 32. GALILEO功率衰减值
	FDelayGPS      float32 `json:"fDelayGPS"`      // 33. GPS通道时延
	FDelayBDS      float32 `json:"fDelayBDS"`      // 34. BDS通道时延
	FDelayGLO      float32 `json:"fDelayGLO"`      // 35. GLONASS通道时延
	FDelayGAL      float32 `json:"fDelayGAL"`      // 36. GALILEO通道时延
	IAutoTimingSw  int     `json:"iAutoTimingSw"`  // 37. 整点授时开关
	FInitSpeedVal  float32 `json:"fInitSpeedVal"`  // 38. 模拟初速度大小
	FInitSpeedHead float32 `json:"fInitSpeedHead"` // 39. 模拟初速度方向
	FAccSpeedVal   float32 `json:"fAccSpeedVal"`   // 40. 模拟加速度大小
	FAccSpeedHead  float32 `json:"fAccSpeedHead"`  // 41. 模拟加速度方向
	FCirRadius     float32 `json:"fCirRadius"`     // 42. 模拟圆周运动半径
	FCirCycle      float32 `json:"fCirCycle"`      // 43. 模拟圆周运动周期
	ICirRotDir     int     `json:"iCirRotDir"`     // 44. 模拟圆周运动方向
	IAutoTranSw    int     `json:"iAutoTranSw"`    // 45. 上电自动发射开关
	IFirstTim      int     `json:"iFirstTim"`      // 46. 上电首次校时标志
	AGPSPRN        []int   `json:"aGPSPRN"`        // 47. GPS转发卫星PRN号
	ABDSPRN        []int   `json:"aBDSPRN"`        // 48. BDS转发卫星PRN号
	AGLOPRN        []int   `json:"aGLOPRN"`        // 49. GLONASS转发卫星PRN号
	AGALPRN        []int   `json:"aGALPRN"`        // 50. GALILEO转发卫星PRN号
	FMaxSpeed      float32 `json:"fMaxSpeed"`      // 51. 模拟最大速度值
	ITranCtrlMode  int     `json:"iTranCtrlMode"`  // 52. 发射控制模式
	IAntSup        int     `json:"iAntSup"`        // 53. 天线馈电控制开关
	IGnssScreenSw  int     `json:"iGnssScreenSw"`  // 54. GNSS屏蔽功能开关
	IPinLevel      int     `json:"iPinLevel"`      // 55. 外部引脚输出电平
	IGalUsedNum    int     `json:"iGalUsedNum"`    // 56. 伽利略可用卫星数
	ISwitchToneGPS int     `json:"iSwitchToneGPS"` // 57. GPS单载波模式开关
	ISwitchToneGLO int     `json:"iSwitchToneGLO"` // 58. GLO单载波模式开关
	IToneGLOChan   int     `json:"iToneGLOChan"`   // 59. GLO单载波通道号
}
