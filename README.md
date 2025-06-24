# TCP Proxy for Alpha Jet Evo

This Go application acts as a simple TCP proxy between clients and an Alpha Jet Evo device. It listens on a configurable port for incoming client connections and forwards all data bidirectionally between the client and the device. All traffic is logged for debugging and monitoring purposes.

## Features

* **Bidirectional Proxying**: Forwards data in both directions (Client → Device and Device → Client).
* **Configurable Endpoints**: Specify the device address and the listen port via command-line flags.
* **Connection Logging**: Logs connection events, bytes transferred, and direction of data flow.
* **Concurrent Handling**: Supports multiple simultaneous clients using goroutines.
* **Graceful Cleanup**: Closes connections and logs when a client or device disconnects.

## Prerequisites

* Go 1.16 or newer installed.
* Network access to the Alpha Jet Evo device (default `192.168.1.184:3000`) from the machine running the proxy.

## Installation

1. Clone the repository (or copy the source file) to your local machine.
2. Ensure your `GOPATH` or Go modules are set up.

## Building

```bash
# If using Go modules:
go mod init example.com/alpha-jet-proxy  # only if starting a new module
# Place the source in the module directory, e.g., main.go

go build -o alpha-jet-proxy main.go
```

This produces an executable named `alpha-jet-proxy`.

## Usage

```bash
./alpha-jet-proxy [flags]
```

### Command-Line Flags

* `-device`: The IP address and port of the Alpha Jet Evo device. Default: `192.168.1.184:3000`.
* `-listen_port`: The local port on which the proxy listens for incoming client connections. Default: `4001`.

Example:

```bash
./alpha-jet-proxy -device=192.168.1.184:3000 -listen_port=4001
```

Upon starting, you will see a message:

```
Proxy listening on :4001 -> 192.168.1.184:3000
```

Clients can then connect to this machine on port `4001`, and their traffic will be forwarded to the Alpha Jet Evo device.

## Logging

The application logs to standard output and standard error:

* **Startup and Shutdown**:

  * Logs when the proxy starts listening.
  * Logs each client connection and disconnection with timestamps.

* **Data Transfer**:

  * Logs each read/write event with the direction (C→D or D→C), number of bytes, and a quoted representation of the data.
  * Example:

    ```
    [127.0.0.1:54321] C->D 128 bytes: "..."
    [127.0.0.1:54321] D->C 64 bytes: "..."
    ```

* **Errors**:

  * Logs errors when failing to accept a connection, connect to the device, or read/write operations.

### Log Format

* Timestamps for connection events use the format `YYYY-MM-DD HH:MM:SS`, e.g., `2025-06-24 14:05:30`.
* Logs include the client's remote address to distinguish between multiple clients.

## Behavior Details

1. **Listening**: The proxy listens on `0.0.0.0:<listen_port>`, accepting TCP connections.
2. **Accepting a Client**: For each client connection:

   * Logs the connection with remote address and timestamp.
   * Establishes a new TCP connection to the device address (`deviceAddr`).
   * Starts two goroutines:

     * **Client → Device**: Reads from the client, logs data, and writes to the device.
     * **Device → Client**: Reads from the device, logs data, and writes back to the client.
   * Uses a `done` channel: when either side closes or errors occur, signals to clean up both connections.
   * Logs when the connection is closed with timestamp.
3. **Concurrency**: Each client is handled in its own goroutines, allowing multiple clients simultaneously. Connections are independent.

## Error Handling

* **Listen Errors**: If binding to the listen port fails, the application exits with an error message.
* **Accept Errors**: Logs an error and continues accepting new connections.
* **Dial Errors to Device**: Logs error and returns, closing the client connection.
* **Read/Write Errors**: When a read or write fails on either side, logs the error and closes connections.

## Configuration and Environment

* Ensure firewall rules allow inbound connections on the chosen `listen_port` and outbound connections to the device address.
* If running on a systemd or similar service manager, you can create a service file to manage the proxy.

### Example Systemd Service (optional)

```ini
[Unit]
Description=Alpha Jet Evo TCP Proxy
After=network.target

[Service]
Type=simple
ExecStart=/path/to/alpha-jet-proxy -device=192.168.1.184:3000 -listen_port=4001
Restart=on-failure
User=proxyuser

[Install]
WantedBy=multi-user.target
```

## Extensibility

* **TLS Support**: If secure connections are needed, wrap `net.Listen` and `net.Dial` with `tls.Listen`/`tls.Dial` and provide certificates.
* **Authentication**: Add authentication or access control for clients before proxying.
* **Metrics**: Integrate Prometheus or another metrics system to record bytes transferred, connection counts, errors, etc.
* **Config File**: Replace command-line flags with a configuration file or environment variables.
* **Buffer Size**: Adjust the buffer size (currently 4096 bytes) if needed for performance tuning.

## Troubleshooting

* **Cannot Bind to Port**: Check if the port is in use or if you lack permissions (ports <1024 require root).
* **Connection to Device Fails**: Verify network connectivity and that the device is listening on the specified address/port.
* **High Latency or Drops**: Ensure network quality; consider enabling TCP keepalive or tuning OS networking parameters.
* **Large or Binary Data**: The proxy logs raw bytes with `%q`, which escapes non-printable bytes. For heavy binary traffic, adjust or disable verbose logging.

## Example Usage Scenario

1. Start the proxy on a server or local machine:

   ```bash
   ./alpha-jet-proxy -device=192.168.1.184:3000 -listen_port=4001
   ```
2. On a client machine or application, connect to the proxy’s host on port 4001 instead of directly to the device.
3. Monitor logs on the proxy to observe traffic and debug communication with the Alpha Jet Evo device.

## Building from Source

Optionally, you can include this in a Git repository. Sample `go.mod`:

```go
module example.com/alpha-jet-proxy

go 1.18
```

Then:

```bash
go mod tidy
go build
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request with clear descriptions and tests if applicable.

## License

Specify an appropriate license (e.g., MIT, Apache 2.0) as needed.
