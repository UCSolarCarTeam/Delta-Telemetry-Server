package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func listener(Conn *net.UDPConn) {
	buf := make([]byte, 2048)
	for {
		fmt.Println("Looking for UDP")
		n, addr, err := Conn.ReadFromUDP(buf)
		fmt.Println("Received something?")
		CheckError(err)
		fmt.Println("Received", string(buf[0:n]), "from", addr.String())
	}
}

func main() {
	fmt.Println("args: <localIP> <ServerIP> <ServerPort>")
	ServerAddr, err := net.ResolveUDPAddr("udp", os.Args[3]+":"+os.Args[4])
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", os.Args[1]+":"+os.Args[2])
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()
	i := 0

	go listener(Conn)
	for {
		fmt.Println("Writing Message")
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg + "\n")
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}
