package mylock

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Data struct {
	value int
	sync.RWMutex
}

var data = Data{}

func (d *Data) Write(newValue int) {
	fmt.Printf("start Goroutine %d wrote value: %d\n", getGoroutineID(), newValue)
	d.Lock()
	defer d.Unlock()
	time.Sleep(1 * time.Second)
	d.value = newValue
	fmt.Printf("--------------> Goroutine %d wrote value: %d\n", getGoroutineID(), newValue)
}
func (d *Data) Read() {
	for {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Goroutine %d read value: %d\n", getGoroutineID(), d.value)
	}

}

func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idStr := strings.TrimPrefix(string(buf[:n]), "goroutine ")
	idStr = idStr[:strings.IndexByte(idStr, ' ')]
	goroutineID, _ := strconv.Atoi(idStr)
	return goroutineID
}

func Show() {

	go func() {
		data.Read()
	}()
	go func() {
		data.Read()
	}()
	go func() {
		data.Read()
	}()

	go func() {
		data.Write(42)
	}()

	go func() {
		data.Write(24)
	}()
	go func() {
		data.Write(33)
	}()

	go func() {
		data.Write(11)
	}()

	time.Sleep(100 * time.Second)

	fmt.Println("Final value:", data.value)
}
