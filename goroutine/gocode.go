package goroutine

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func reportSystemHeartbeat() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		fmt.Println("Current Time:", formattedTime)
	}
}

func mylog(v ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	log.Printf("%s:%d %s: %s", file, line, funcName, v)
	// log.Printf(file, ":", line, ":", funcName, ":")
	// log.Println(v...)
	// log.Printf("%s", v)
}

func Gorou1() {
	go Go1()
	// go Go2()
	// os.StartProcess(Go1())
}

// 获取当前线程ID
func getThreadID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseUint(idField, 10, 64)
	return id
}

// 获取当前协程ID
func getGoroutineID() uint64 {
	var buf [64]byte
	runtime.Stack(buf[:], false)
	gidField := strings.Fields(strings.TrimPrefix(string(buf[:]), "goroutine "))[0]
	gid, _ := strconv.ParseUint(gidField, 10, 64)
	return gid
}
func Go1() {
	fmt.Println("当前进程 ID:", os.Getpid())
	fmt.Println("当前线程 ID:", getThreadID())
	fmt.Println("当前协程 ID:", getGoroutineID())
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("r = ", r)
			Go1()
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	index := 0
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			nowtime := time.Now().Format("2006-01-02 15:04:05")
			mylog("Go1  ++++++++++++++++ ", nowtime)
			index++
			if index%3 == 0 {
				panic("make panic")
				// done <- true
				// ticker.Stop()
				// os.Exit(0)
			}
		}
	}
	// mylog("GO1 exit ....")
	// ticker.Stop()
	// done <- true
}

func Go2() {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			nowtime := time.Now().Format("2006-01-02 15:04:05")
			mylog("Go2------------- = ", nowtime)

		}
	}
	// ticker.Stop()
	// done <- true
}
