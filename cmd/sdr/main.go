package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
)

// Define CLI options
type Options struct {
	URI    string `short:"u" long:"uri" default:"http://localhost/" description:"dd-remote server address"`
	Tester bool   `short:"t" long:"tester" description:"connect SD to tester"`
	DUT    bool   `short:"d" long:"dut" description:"connect SD to device"`
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

	if options.DUT && options.Tester {
		fmt.Fprintln(os.Stderr, "Error: can not select both dut and tester")
		os.Exit(1)
	}
	if !options.DUT && !options.Tester {
		fmt.Fprintln(os.Stderr, "Error: must select one of dut and tester")
		os.Exit(1)
	}

	// request remote action
	ok = setMuxRemote(options.URI, options.DUT)
	if !ok {
		os.Exit(1)
	}
}

func setMuxRemote(uri string, state bool) bool {
	var err error
	var response *http.Response

	// execute http GET request to given uri
	response, err = http.Get(fmt.Sprintf("%s?state=%t", uri, state))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return true
	} else {
		fmt.Fprintf(os.Stderr, "server returned %d\n", response.StatusCode)
		return false
	}
}
