package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"

	"github.com/jessevdk/go-flags"
)

// Define CLI options
type Options struct {
	ListenHost string `short:"l" long:"listen-host" default:"localhost" description:"address to listen on for incoming connections"`
	ListenPort uint   `short:"p" long:"listen-port" default:"80" description:"port to listen on for incoming connections"`
}

func main() {
	var err error
	var options Options
	var ok bool

	// parse cli options
	_, err = flags.Parse(&options)
	if err != nil {
		os.Exit(1)
	}

	// start server
	ok = startMuxServer(options.ListenHost, options.ListenPort)
	if !ok {
		os.Exit(1)
	}
}

// start http server
func startMuxServer(host string, port uint) bool {
	var err error
	http.Handle("/", &requestHandler{})
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return false
	}

	return true
}

// define http server state
type requestHandler struct {
	Busy sync.Mutex
}

// define handler function for all http requests
func (requestHandler *requestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		requestHandler.handleGet(responseWriter, request)
	default:
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// define handler function for get requests
func (requestHandler *requestHandler) handleGet(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	var command *exec.Cmd
	var values url.Values
	var statearg string
	// serialize requests
	requestHandler.Busy.Lock()
	defer requestHandler.Busy.Unlock()

	values = request.URL.Query()
	stateval, ok := values["state"]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: request does not have \"state\" value in query")

		// send error response
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(stateval) > 1 {
		fmt.Fprintln(os.Stderr, "Error: request repeated \"state\" value in query")

		// send error response
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	switch stateval[0] {
	case "false":
		statearg = "--ts"
	case "true":
		statearg = "--dut"
	default:
		fmt.Fprintf(os.Stderr, "Error: request has invalid \"state\" value in query: \"%s\"\n", stateval[0])

		// send error response
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	// call cli sd-mux-ctrl
	command = exec.Command("/usr/sbin/sd-mux-ctrl", "-e", "sdwire-24", statearg)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		// send error response
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send success response
	responseWriter.WriteHeader(http.StatusNoContent)
}
