# TCP Proxy for Alpha Jet Evo

This Go application acts as a simple TCP proxy between clients and an Alpha Jet Evo device. It listens on a configurable port for incoming client connections and forwards all data bidirectionally between the client and the device. All traffic is logged for debugging and monitoring purposes.

## Features

- **Bidirectional Proxying**: Forwards data in both directions (Client → Device and Device → Client).
- **Configurable Endpoints**: Specify the device address and the listen port via command-line flags.
- **Connection Logging**: Logs connection events, bytes transferred, and direction of data flow.
- **Concurrent Handling**: Supports multiple simultaneous clients using goroutines.
- **Graceful Cleanup**: Closes connections and logs when a client or device disconnects.

## Prerequisites

- Go 1.16 or newer installed.
- Network access to the Alpha Jet Evo device (default `192.168.1.184:3000`) from the machine running the proxy.

## Installation

1. Clone the repository (or copy the source file) to your local machine.
2. Ensure your `GOPATH` or Go modules are set up.

## Building

```bash
# If using Go modules:
go mod init example.com/alpha-jet-proxy  # only if starting a new module
# Place the source in the module directory, e.g., main.go

go build -o alpha-jet-proxy main.go
