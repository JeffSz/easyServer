# Example1:
```
import (
    "net/http"
    "io"
    "time"
    "easyServer"
    "net"
    "fmt"
    "io/ioutil"
)

server := *easyServer.EasyServer

server.AddRoute("/test", easyServer.HTTP_GET, func(w http.ResponseWriter, r *http.Request){
    io.WriteString(w, "Ok")
})

server.AddRoute("/test", easyServer.HTTP_POST, func(w http.ResponseWriter, r *http.Request){
    io.WriteString(w, r.Method)
})

if err := (&http.Server{
    Addr: fmt.Sprintf("%s:%s", ip, *port),
    Handler: server,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 <<20,
}).ListenAndServe(); err != nil{
    panic(err)
}
```