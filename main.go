package main

import (
	"MyIOSWebServer/router"
	ser "MyIOSWebServer/server"
	"os"
	"fmt"
	"net/http"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8082" //mux address
)

var addr = flag.String("addr", ":8083", "websocket address")

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for http listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}
	server := router.NewServer()
	go server.Run(":" + port)

	//socketf服务
	http.HandleFunc("/ws", ser.ServeWebSocket)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Println("error in socket server")
	}

}
