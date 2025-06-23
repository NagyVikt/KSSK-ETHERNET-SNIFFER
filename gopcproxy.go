package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "time"
)

var (
    deviceAddr = flag.String("device", "192.168.1.184:3000", "Alpha Jet Evo address")
    listenPort = flag.String("listen_port", "4001", "proxy listen port for clients")
)

func die(format string, v ...interface{}) {
    fmt.Fprintf(os.Stderr, format+"\n", v...)
    os.Exit(1)
}

func main() {
    flag.Parse()

    ln, err := net.Listen("tcp", ":"+*listenPort)
    if err != nil {
        die("Failed to listen on port %s: %v", *listenPort, err)
    }
    fmt.Printf("Proxy listening on :%s -> %s\n", *listenPort, *deviceAddr)

    for {
        clientConn, err := ln.Accept()
        if err != nil {
            fmt.Printf("Accept error: %v\n", err)
            continue
        }
        go handleClient(clientConn)
    }
}

func handleClient(clientConn net.Conn) {
    defer clientConn.Close()
    clientAddr := clientConn.RemoteAddr().String()
    fmt.Printf("[%s] Client connected at %s\n", clientAddr, time.Now().Format("2006-01-02 15:04:05"))

    // Dial the device for each client
    deviceConn, err := net.Dial("tcp", *deviceAddr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to device %s: %v\n", *deviceAddr, err)
        return
    }
    defer deviceConn.Close()
    fmt.Printf("[%s] Connected to device for client\n", clientAddr)

    done := make(chan struct{})

    // Client -> Device
    go func() {
        buf := make([]byte, 4096)
        for {
            n, err := clientConn.Read(buf)
            if n > 0 {
                data := buf[:n]
                fmt.Printf("[%s] C->D %d bytes: %q\n", clientAddr, n, data)
                _, werr := deviceConn.Write(data)
                if werr != nil {
                    fmt.Fprintf(os.Stderr, "Write to device failed: %v\n", werr)
                    break
                }
            }
            if err != nil {
                break
            }
        }
        done <- struct{}{}
    }()

    // Device -> Client
    go func() {
        buf := make([]byte, 4096)
        for {
            n, err := deviceConn.Read(buf)
            if n > 0 {
                data := buf[:n]
                fmt.Printf("[%s] D->C %d bytes: %q\n", clientAddr, n, data)
                _, werr := clientConn.Write(data)
                if werr != nil {
                    fmt.Fprintf(os.Stderr, "Write to client failed: %v\n", werr)
                    break
                }
            }
            if err != nil {
                break
            }
        }
        done <- struct{}{}
    }()

    // Wait for one side
    <-done
    fmt.Printf("[%s] Connection closed at %s\n", clientAddr, time.Now().Format("2006-01-02 15:04:05"))
}
