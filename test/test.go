package main

import (
	"net/http"
	"io"
	"time"
	"easyServer"
	"net"
	"fmt"
	"flag"
	"github.com/smartystreets/assertions"
	"io/ioutil"
	"context"
)

func main() {
	port := flag.String("port", "8001", "server port")
	flag.Parse()

	ip := "localhost"
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		panic(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	server := *easyServer.EasyServer

	server.AddRoute("/test", easyServer.HTTP_GET, func(w http.ResponseWriter, r *http.Request){
		io.WriteString(w, "Ok")
	})

	server.AddRoute("/test", easyServer.HTTP_POST, func(w http.ResponseWriter, r *http.Request){
		io.WriteString(w, r.Method)
	})

	webServer := &http.Server{
		Addr: fmt.Sprintf("%s:%s", ip, *port),
		Handler: server,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 <<20,
	}

	go func(){
		time.Sleep(2 * time.Second)

		fmt.Println("====:====testing====:====")

		resp, err := http.Get(fmt.Sprintf("http://%s:%s/test", ip, *port))
		assertions.ShouldBeNil(err)
		assertions.ShouldNotBeNil(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		assertions.ShouldBeNil(err)
		resp.Body.Close()
		assertions.ShouldEqual(body, "Ok")
		fmt.Println(body)
		fmt.Println("Get test passed")

		resp, err = http.Get(fmt.Sprintf("http://%s:%s/test", ip, *port))
		assertions.ShouldBeNil(err)
		assertions.ShouldNotBeNil(resp.Body)
		body, err = ioutil.ReadAll(resp.Body)
		assertions.ShouldBeNil(err)
		resp.Body.Close()
		assertions.ShouldEqual(body, "POST")
		fmt.Println("Post test passed")

		fmt.Println("====:====tested====:====")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		webServer.Shutdown(ctx)
	}()

	fmt.Printf("Server start on %s: %s \n", ip, *port)

	if err = webServer.ListenAndServe(); err != nil{
		panic(err)
	}
}
