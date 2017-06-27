package easyServer


import (
	"time"
	"regexp"
	"net/http"
	"fmt"
	"io"
	"errors"
	logger "github.com/JeffSz/logger"
	"strings"
	"net"
)
func time_sub(t time.Time) int64 {
	return int64(time.Now().Sub(t) / time.Millisecond)
}
type Method int
var HTTPMethods = make(map[string]Method)

const(
	HTTP_GET = 0x01
	HTTP_POST = 0x02
	HTTP_PUT = 0x04
	HTTP_DELETE = 0x08
	HTTP_ALL = 0xFF
)

type MyRouter struct{
	method Method
	pattern *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request)
}
type Server struct{
	routes []MyRouter
	errorHandler func(w http.ResponseWriter, error Error) error
}
func (server *Server) AddRoute(url string, method Method, handler func(http.ResponseWriter, *http.Request)) error{
	if method != HTTP_GET && method != HTTP_POST && method != HTTP_ALL{
		return errors.New("Method not allow")
	}
	if pattern, err := regexp.Compile(url); err != nil{
		return err
	}else{
		server.routes = append(server.routes, MyRouter{method: method, pattern: pattern, handler: handler})
		return nil
	}
}

func (server Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start_time := time.Now()
	defer func() {
		if err := recover(); err != nil {
			logger.Info(fmt.Sprintf("%s\n", stack()))
			if er, is := err.(Error); is{
				server.errorHandler(w, er)
				logger.Debug(time.Now(), "request:", r, time_sub(start_time), http.StatusBadRequest)
			}else{
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Server Error")
				logger.Error(err)
				logger.Error(time.Now(), "request:", r, time_sub(start_time), http.StatusInternalServerError)
			}
		}
	}()
	for _, router := range server.routes{
		if router.pattern.MatchString(r.URL.Path) && (router.method & HTTPMethods[strings.ToUpper(r.Method)] != 0x00){
			router.handler(w, r)
			logger.Debug(time.Now(), "request:", r, time_sub(start_time), http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Page not found")
	logger.Warn(time.Now(), "request:", r, time_sub(start_time), http.StatusNotFound)
}

func init() {
	HTTPMethods["GET"] = HTTP_GET
	HTTPMethods["POST"] = HTTP_POST
	HTTPMethods["PUT"] = HTTP_PUT
	HTTPMethods["DELETE"] = HTTP_DELETE
	HTTPMethods["ALL"] = HTTP_ALL
}

func LocalIP() string{
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		panic(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("No such device")
}

func NewServer(errorHandler func(http.ResponseWriter, Error) error) *Server{
	return &Server{routes:make([]MyRouter, 0), errorHandler: errorHandler}
}

func easyHander(w http.ResponseWriter, error Error) error{
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, error.Error())
	return nil
}

var EasyServer = NewServer(easyHander)
