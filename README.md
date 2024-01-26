# Course Registration Tool

The Course Registration Tool is a command-line utility written in Go for automating the registration process for courses on the Sharif University of Technology's educational platform. It connects to the platform's WebSocket API to monitor course availability and attempts to register for your favorite courses automatically.

## Features

- **WebSocket Integration**: Utilizes WebSocket communication to receive real-time updates on course availability.

- **Course Registration**: Attempts to register for your favorite courses automatically.

- **Customizable Delays**: Allows you to set the delay between registration attempts to avoid overloading the server.

- **Retries**: Defines the maximum number of registration retries in case of failure.

- **Infinite Registration**: Option to keep attempting registration indefinitely until successful.

## How to Build

You can build the Course Registration Tool for different operating systems (OS) using the Go cross-compilation feature. Here are the steps to build it for various OSs:

### Prerequisites

Make sure you have Go (Golang) installed on your machine.

### Build for Windows

```bash
go build -o cr.exe
```

### Build for Linux

```bash
GOOS=linux GOARCH=amd64 go build
```

### Build for macOS

```bash
GOOS=darwin GOARCH=amd64 go build
```

### How to run

After building the executable for your desired OS, you can run it from the command line. Here's the basic usage:

```bash
./cr.exe -d [DELAY_SECONDS] -r [MAX_RETRIES] -i
```

`-d [DELAY_SECONDS]`: Specifies the delay in seconds between registration attempts (default: 5 seconds).

`-r [MAX_RETRIES]`: Sets the maximum number of registration retries (default: 5 retries).

`-i`: Enables infinite registration attempts until successful (default: false).

`-on-time`: Enable on-time registration. (default: false).

`-config`: Path to the configuration file (default: config.json)

`-o`: Offset in milliseconds before the first registration request (default: 300)
