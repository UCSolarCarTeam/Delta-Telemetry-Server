package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "time"
    "bufio"
)

func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
    }
}

func main() {
    fmt.Println("args: <ServerIP> <ServerPort>")
    ServerAddr, err := net.ResolveUDPAddr("udp", os.Args[1]+":"+os.Args[2])
    CheckError(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
    CheckError(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)

    defer Conn.Close()
    i := 0

    for {       
        /* Write to server */
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _, err := Conn.Write(buf)
        CheckError(err)

        /* Read from server */
        buf =  make([]byte, 2048)
        _, err = bufio.NewReader(Conn).Read(buf)
        if err == nil {
            fmt.Printf("%s\n", buf)
        } else {
            fmt.Printf("Some error %v\n", err)
        }

        time.Sleep(time.Second * 1)
    }
}
