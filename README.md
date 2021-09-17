# sdmux-remote

An HTTP API wrapping the [sdmux control utility](https://git.tizen.org/cgit/tools/testlab/sd-mux/) for switching microSD-Cards between Tester and DUT.
A server application listens for incoming GET requests and and calls sd-mux-ctrl appropriately.
A client application sends GET requests to a given URI.

sdmux-remote queues all requests to avoid concurrent invocations of the underlying sd-mux-ctrl.

## Server

```
go get git.home.jm0.eu/josua/sdmux-remote/cmd/sdrd

sdrd --help
Usage:
  sdrd [OPTIONS]

Application Options:
  -l, --listen-host= address to listen on for incoming connections (default: localhost)
  -p, --listen-port= port to listen on for incoming connections (default: 80)

Help Options:
  -h, --help         Show this help message
```

## Client

```
go get git.home.jm0.eu/josua/sdmux-remote/cmd/sdr

sdr --help
Usage:
  sdr [OPTIONS]

Application Options:
  -u, --uri=    dd-remote server address (default: http://localhost/)
  -t, --tester  connect SD to tester
  -d, --dut     connect SD to device

Help Options:
  -h, --help    Show this help message
```
