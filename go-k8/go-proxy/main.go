package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var logFields = log.Fields{
	"application": "go-proxy",
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}
	proxy.ModifyResponse = modifyResponse(targetHost)
	proxy.ErrorHandler = errorHandler()

	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logFields).WithField("event", "proxyRequest").Info("Served Proxy")
		proxy.ServeHTTP(w, r)
	}
}

func modifyResponse(url string) func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "Magical")
		resp.Header.Add("Where-Am-I", resp.Request.Host)
		return nil
	}
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		errorLogger(err.Error())
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

func errorLogger(err string) {
	log.WithFields(logFields).WithField("event", "error").Error(err)
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	log.WithFields(logFields).Info("Proxy Start")
	proxy, err := NewProxy("http://localhost:8081")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	http.ListenAndServe(":8080", nil)
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
