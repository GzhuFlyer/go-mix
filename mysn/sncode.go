package mysn

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

<<<<<<< HEAD
=======
	"log"

>>>>>>> 286a0e1 ([feat]:udp)
	"github.com/denisbrodbeck/machineid"
)

func Sn1() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	// 遍历每个网络接口
	for _, iface := range ifaces {
		// 获取接口的名称
		name := iface.Name

		// 获取接口的 MAC 地址
		mac := iface.HardwareAddr.String()

		// 获取接口的 IP 地址
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		var ips []string
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ips = append(ips, ipnet.IP.String())
			}
		}
		ip := strings.Join(ips, ", ")

		// 打印接口的信息
		fmt.Printf("%s\t%s\t%s\n", name, mac, ip)
	}
	fmt.Println("")
	fmt.Printf(" runtime.GOOS  = %s\n", runtime.GOOS)
}

func Sn2() {
	// filePath := "/mnt/vendor/persist/devSN"
	filePath := "E:\\workspace\\com\\project\\websocket\\go-util\\b.go"

	// 判断文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在

		// 读取文件内容
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("读取文件失败:", err)
			return
		}

		// 打印文件内容
		fmt.Println("data = ", data)
		fmt.Println("文件内容:", string(data))
	} else if os.IsNotExist(err) {
		// 文件不存在
		fmt.Println("文件不存在")
	} else {
		// 其他错误
		fmt.Println("发生错误:", err)
	}

}
<<<<<<< HEAD
=======

var snFilePath = "./sn"

func GetDevSn() (string, error) {
	if _, err := os.Stat(snFilePath); err == nil {
		data, err := os.ReadFile(snFilePath)
		if err != nil {
			fmt.Println("read snFile fail:", err)
			return "", err
		}
		if len(data) == 0 {
			errMsg := fmt.Sprintf("%s not contain sn\n", snFilePath)
			return "", errors.New(errMsg)
		}
		strData := string(data)
		return strData, err
	} else if os.IsNotExist(err) {
		// 文件不存在
		fmt.Println("snFil not exit")
		return "", err
	} else {
		// 其他错误
		fmt.Println("unkonw:", err)
		return "", err
	}
}

>>>>>>> 286a0e1 ([feat]:udp)
func WinSn() {
	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
