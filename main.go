package main

import (
	"flag"
	"github.com/caikaijie/goconnect/handler"
	"log"
	"net/http"
)

func main() {
	var ipport string
	flag.StringVar(&ipport, "P", "127.0.0.1:8080", "[:port|ip:port]")
	flag.Parse()

	s := &http.Server{
		Addr:    ipport,
		Handler: handler.New(),
	}

	log.Print("ListenAndServe on ", ipport, "\n")
	log.Fatal(s.ListenAndServe())
}
