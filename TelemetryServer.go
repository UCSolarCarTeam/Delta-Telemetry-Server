package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

var KeepAliveMap map[string]KeepAlive

const TimeOut = 20

type KeepAlive struct {
    address   *net.UDPAddr
    timestamp time.Time
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(0)
    }
}

func cleanMap() {
    for _, contents := range KeepAliveMap {
        fmt.Println(contents.timestamp)
        if time.Since(contents.timestamp).Seconds() > TimeOut {
            fmt.Println(contents.address, " : TIMED OUT")
            delete(KeepAliveMap, contents.address.String())
        }
    }
}

func updateMap(addr *net.UDPAddr) {
    time_receieved := time.Now()
    temp := KeepAlive{addr, time_receieved}
    KeepAliveMap[addr.String()] = temp
    cleanMap()
}

func forwardMessagesUdp(ServerConn *net.UDPConn, message []byte, length int) {
    if(string(message) == "Heartbeat"){
        fmt.Println("NOT SENDING A HEARTBEAT...")
        return
    }
    for _, clients := range KeepAliveMap {
        fmt.Println("SENDING TO CLIENT: ", clients.address.String())
        _, err := ServerConn.WriteToUDP(message, clients.address)
        checkError(err)
    }
}

func receiveAndPrintUdp(ServerConn *net.UDPConn, buf []byte) (int, error) {
    n, addr, err := ServerConn.ReadFromUDP(buf)
    response := string(buf[0:n])
    fmt.Println("\nRECEIVED: ", response, "\nFROM: ", addr)
    updateMap(addr)
    return n, err
}

func main() {
    KeepAliveMap = make(map[string]KeepAlive)
    ServerAddr, err := net.ResolveUDPAddr("udp", ":" + os.Args[1])
    checkError(err)
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    checkError(err)
    defer ServerConn.Close()

    buf := make([]byte, 2048)

    for {
        n, err := receiveAndPrintUdp(ServerConn, buf)
        go forwardMessagesUdp(ServerConn, buf, n)
        checkError(err)
    }
}
