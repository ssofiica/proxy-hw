package proxy

import (
	"io"
	"log"
	"net/http"
)

func HandlerHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Print(err)
	}

	io.Copy(w, resp.Body)
	if err = resp.Body.Close(); err != nil {
		log.Print(err)
	}

	copyHeader(resp.Header, w.Header())
}

func copyHeader(from, where http.Header) {
	for key, values := range from {
		for _, v := range values {
			where.Add(key, v)
		}
	}
}
