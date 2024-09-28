package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssofiica/proxy-hw/internal/proxy"
	"github.com/ssofiica/proxy-hw/internal/repo"
)

var (
	PROXY_PORT    = ":8080"
	HOST          = "localhost"
	POSTGRES_CONN = "postgres://svalova:mydbpass@localhost:5432/test-gaz"
)

func main() {
	var crt, key string
	flag.StringVar(&crt, "crt", "certs/ca.crt", "")
	flag.StringVar(&key, "key", "certs/ca.key", "")
	var protocol string
	flag.StringVar(&protocol, "protocol", "https", "")
	flag.Parse()

	db, err := pgxpool.New(context.Background(), POSTGRES_CONN)
	if err != nil {
		fmt.Println("error wih db", err)
	}
	repository := repo.NewRepo(db)
	proxy := proxy.NewHandler(repository)

	proxyServer := &http.Server{
		Addr: PROXY_PORT,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print(r.Method, r.URL)
			if r.Method == http.MethodConnect {
				proxy.HandlerConnect(w, r)
			} else {
				proxy.HandlerHTTP(w, r)
			}
		}),
		TLSConfig: &tls.Config{ServerName: HOST},
	}
	if protocol == "https" {
		if err := proxyServer.ListenAndServeTLS(crt, key); err != nil {
			log.Fatal("proxyServer wasn't started:", err)
		}
	} else if protocol == "http" {
		if err := proxyServer.ListenAndServe(); err != nil {
			log.Fatal("proxyServer wasn't started:", err)
		}
	}

}
