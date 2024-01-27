# Course Registration Tool

The Course Registration Tool is a command-line utility written in Go for automating the registration process for courses on the Sharif University of Technology's educational platform. It connects to the platform's WebSocket API to monitor course availability and attempts to register for your favorite courses automatically.

## Features

- **WebSocket Integration**: Utilizes WebSocket communication to receive real-time updates on course availability.

- **Course Registration**: Attempts to register for your favorite courses automatically.

- **Customizable Delays**: Allows you to set the delay between registration attempts to avoid overloading the server.

- **Retries**: Defines the maximum number of registration retries in case of failure.

- **Infinite Registration**: Option to keep attempting registration indefinitely until successful.

### Requirements
```bash
go mod download
go mod tidy
```

## How to Build

You can build the Course Registration Tool for different operating systems (OS) using the Go cross-compilation feature. Here are the steps to build it for various OSs:

### Prerequisites

Make sure you have Go (Golang) installed on your machine.

### Build for Windows

```bash
go build -o ./bin/cr
```

### Build for Linux

```bash
GOOS=linux GOARCH=amd64 go build -o ./bin/cr
```

### Build for macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o ./bin/cr
```

## How to run

First of all, copy the `example.json` to  `config.json`:

For windows:
```bash
copy ./example.json ./config.json
```

For Linux & macOS:
```bash
cp copy ./example.json ./config.json
```

You should replace `<YOUR-TOKEN-HERE>` with your authorization otken. You can obtain it from local storage of [my.edu.sharif.edu](https://my.edu.sharif.edu):

F12 -> Application -> Local Storage -> https://my.edu.sharif.edu -> value of token

Then you should determine the value of `fav` and `action` in `config.json`. You should always set a value for these.

```json
    "fav": [],
    "action": ""
```

The default value of `fav` will be your favorite courses on [my.edu.sharif.edu](https://my.edu.sharif.edu).

The default value of `action` will be `add`. Other options are `move` and `remove`.


After building the executable for your desired OS, you can run it from the command line. Here's the basic usage:

```bash
./cr -d [DELAY_SECONDS] -r [MAX_RETRIES] -i -config [PATH] -o [OFFSET]
```

`-d [DELAY_SECONDS]`: Specifies the delay in seconds between registration attempts (default: `5`)‍.

`-r [MAX_RETRIES]`: Sets the maximum number of registration retries (default: `5`).

`-i`: Enables infinite registration attempts until successful (default: `false`).

`-on-time`: Enable on-time registration. (default: `false`).

`-config [PATH]`: Path to the configuration file (default: `config.json`)

`-o [OFFSET]`: Offset in milliseconds before the first registration request (default: `300`)
