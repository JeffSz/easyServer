package easyServer

import (
	"time"
	"regexp"
	"net/http"
	"fmt"
	"io"
	"errors"
	"github.com/JeffSz/logger"
	"strings"
	"net"
	"runtime/debug"
)

func timeSub(t time.Time) int64 {
	return int64(time.Now().Sub(t) / time.Millisecond)
}

type Method int
type ErrorHandler func(w http.ResponseWriter, error error) error
type RouteHandler func(http.ResponseWriter, *http.Request)

var HTTPMethods = make(map[string]Method)

const (
	HTTP_GET    = 0x01
	HTTP_POST   = 0x02
	HTTP_PUT    = 0x04
	HTTP_DELETE = 0x08
	HTTP_ALL    = 0xFF
)

type MyRouter struct {
	method  Method
	pattern *regexp.Regexp
	handler RouteHandler
}
type Server struct {
	routes       []MyRouter
	errorHandler ErrorHandler
}

func (server *Server) AddRoute(url string, method Method, handler RouteHandler) error {
	if method != HTTP_GET && method != HTTP_POST && method != HTTP_ALL {
		return errors.New("method not allow")
	}
	if pattern, err := regexp.Compile(url); err != nil {
		return err
	} else {
		server.routes = append(server.routes, MyRouter{method: method, pattern: pattern, handler: handler})
		return nil
	}
}

func (server *Server) SetErrorHandler(errorHandler ErrorHandler) {
	server.errorHandler = errorHandler
}

func (server Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Sprintf("Error found: %s\n", debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Server Error")
			logger.Error(err)
			logger.Error(time.Now(), "request:", r, timeSub(startTime), http.StatusInternalServerError)

			if er, is := err.(error); is {
				server.errorHandler(w, er)
			}
		}
	}()
	for _, router := range server.routes {
		if router.pattern.MatchString(r.URL.Path) && (router.method&HTTPMethods[strings.ToUpper(r.Method)] != 0x00) {
			router.handler(w, r)
			logger.Debug(time.Now(), "request:", *r, timeSub(startTime), http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Page not found")
	logger.Warn(time.Now(), "request:", *r, timeSub(startTime), http.StatusNotFound)
}

func init() {
	HTTPMethods["GET"] = HTTP_GET
	HTTPMethods["POST"] = HTTP_POST
	HTTPMethods["PUT"] = HTTP_PUT
	HTTPMethods["DELETE"] = HTTP_DELETE
	HTTPMethods["ALL"] = HTTP_ALL
}

func LocalIP() string {
	address, err := net.InterfaceAddrs()

	if err != nil {
		panic(err)
	}

	for _, address := range address {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	panic("No such device")
}

func NewServer(errorHandler ErrorHandler) *Server {
	return &Server{routes: make([]MyRouter, 0), errorHandler: errorHandler}
}

func easyHandler(w http.ResponseWriter, error error) error {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, error.Error())
	return nil
}

var EasyServer = NewServer(easyHandler)
