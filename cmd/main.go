package main

import (
	"log"
	"net/http"

	"github.com/ssofiica/proxy-hw/internal/proxy"
)

var (
	PROXY_PORT = ":8080"
)

func main() {
	server := &http.Server{
		Addr: PROXY_PORT,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				// handler for https
			} else {
				proxy.HandlerHTTP(w, r)
			}
		}),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server wasn't started:", err)
	}
}
