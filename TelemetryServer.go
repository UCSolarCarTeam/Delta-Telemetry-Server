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
	//"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	fmt.Print("Hello\n")

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
