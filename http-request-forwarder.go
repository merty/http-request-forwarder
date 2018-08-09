package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var (
	listen  = flag.String("l", "8080", "source port to listen to")
	targets = flag.String("h", "127.0.0.1:8081,127.0.0.1:8082", "target hosts to forward the requests to")
	timeout = flag.Int("t", 3000, "request forwarding timeout in milliseconds")
	verbose = flag.Bool("v", false, "verbose output mode")
)

var t *http.Transport

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error {
	return nil
}

type handler struct {
	TargetHosts []string
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)

	for _, targetHost := range h.TargetHosts {
		forwardRequest(targetHost, req, requestBody)
	}

	defer req.Body.Close()
	defer w.WriteHeader(200)
}

func forwardRequest(targetHost string, req *http.Request, requestBody []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[err] Recovered from panic! %s", r)
		}
	}()

	targetRequest := createRequest(targetHost, req, requestBody)
	targetResponse, _ := t.RoundTrip(targetRequest)

	defer targetResponse.Body.Close()
}

func createRequest(targetHost string, req *http.Request, requestBody []byte) (targetRequest *http.Request) {
	targetURL, _ := url.Parse("http://" + targetHost + req.URL.String())

	b := new(bytes.Buffer)
	b.Write(requestBody)

	targetRequest = &http.Request{
		Method:        req.Method,
		URL:           targetURL,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Header:        req.Header,
		Body:          readCloser{b},
		ContentLength: req.ContentLength,
		Close:         true,
		Host:          req.Host,
	}

	return
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var ln net.Listener
	var err error

	ln, err = net.Listen("tcp", ":"+*listen)

	if err != nil {
		log.Fatalf("Cannot listen to port: %s", *listen)
	}

	log.Printf("Forwarding incoming traffic from port %s to the following hosts: %s", *listen, *targets)

	h := handler{
		TargetHosts: strings.Split(*targets, ","),
	}

	timeout := time.Duration(*timeout) * time.Millisecond
	t = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 0,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       timeout,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 0,
	}

	s := &http.Server{
		Handler: h,
	}
	s.SetKeepAlivesEnabled(false)
	s.Serve(ln)
}
