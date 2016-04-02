package main

import (
	"fmt"
	"net"
	//"github.com/pborman/getopt"
	//"gopkg.in/ini.v1"
	"os"
	//"os/exec"
	//"path/filepath"
	//"strconv"
	//"strings"
	"time"
)

var KeepAliveMap map[string]KeepAlive

const TimeOut = 5
const HeartBeatPacket = "KeepMeAlive!"

type KeepAlive struct {
	address   *net.UDPAddr
	timestamp time.Time
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

/* cleans out any entries that have been in the Map for longer than the timeout */
func CleanMap() {
	for _, contents := range KeepAliveMap {
		fmt.Println(contents.timestamp)
		if time.Since(contents.timestamp).Seconds() > TimeOut {
			fmt.Println(contents.address, "timed out")
			delete(KeepAliveMap, contents.address.String())
		}
	}
}

func UpdateMap(addr *net.UDPAddr) {
	time_receieved := time.Now()
	temp := KeepAlive{addr, time_receieved}
	KeepAliveMap[addr.String()] = temp
	fmt.Println(KeepAliveMap)
	CleanMap()
}

func ReceiveAndPrintUDP(ServerConn *net.UDPConn, buf []byte) error {
	n, addr, err := ServerConn.ReadFromUDP(buf)
	fmt.Println("Received ", string(buf[0:n]), " from ", addr)
	UpdateMap(addr)
	return err
}

func main() {
	KeepAliveMap = make(map[string]KeepAlive)
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 2048)

	for {
		err := ReceiveAndPrintUDP(ServerConn, buf)
		CheckError(err)
	}
}
