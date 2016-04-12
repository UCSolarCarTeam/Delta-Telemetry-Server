package main

import (
	"fmt"
	"net"
	//"github.com/pborman/getopt"
	//"gopkg.in/ini.v1"
	"os"
	//"os/exec"
	//"path/filepath"
	"strconv"
	//"strings"
	"time"
)

var KeepAliveMap map[string]KeepAlive

const TimeOut = 105
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
	//fmt.Println(KeepAliveMap)
	CleanMap()
}
func forwardMessagesUDP(ServerAddr *net.UDPAddr, message []byte, length int) {
	for _, clients := range KeepAliveMap {
		port, err := strconv.Atoi(os.Args[1])
		port += port + 0
		CheckError(err)
		ServerAddrr, err := net.ResolveUDPAddr("udp", ":"+string(port))
		Conn, err := net.DialUDP("udp", ServerAddrr, clients.address)
		CheckError(err)
		defer Conn.Close()
		fmt.Println("Sending message to", clients.address.String())
		message = []byte("hi\n")
		Conn.Write(message)
		CheckError(err)
	}
}

func ReceiveAndPrintUDP(ServerConn *net.UDPConn, buf []byte) (int, error) {
	n, addr, err := ServerConn.ReadFromUDP(buf)
	//fmt.Println("Received ", string(buf[0:n]), " from ", addr)
	UpdateMap(addr)
	return n, err
}

func main() {
	KeepAliveMap = make(map[string]KeepAlive)
	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+os.Args[1])
	CheckError(err)
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 2048)

	for {
		n, err := ReceiveAndPrintUDP(ServerConn, buf)
		fmt.Println("Received ", string(buf[0:n]))
		forwardMessagesUDP(ServerAddr, buf, n)
		CheckError(err)
	}
}
