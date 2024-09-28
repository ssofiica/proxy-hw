package proxy

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/ssofiica/proxy-hw/internal/api"
	"github.com/ssofiica/proxy-hw/internal/proxy/utils"
	"github.com/ssofiica/proxy-hw/internal/repo"
)

type Handler struct {
	repo repo.Repo
}

func NewHandler(repo repo.Repo) Handler {
	return Handler{
		repo: repo,
	}
}

func (h *Handler) HandlerHTTP(w http.ResponseWriter, r *http.Request) {
	//r.Header.Del("Proxy-Connection")
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Print(err)
		api.InternalServerError(w)
	}

	io.Copy(w, resp.Body)
	if err = resp.Body.Close(); err != nil {
		log.Print(err)
	}
	copyHeader(resp.Header, w.Header())

	reqInfo := utils.ParseRequest(r)
	req, err := json.Marshal(reqInfo)
	if err != nil {
		log.Print(err)
		api.InternalServerError(w)
		return
	}
	respInfo := utils.ParseResponse(resp)
	res, err := json.Marshal(respInfo)
	if err != nil {
		log.Print(err)
		api.InternalServerError(w)
		return
	}
	id, err := h.repo.SaveRequest(context.Background(), req)
	if err != nil {
		log.Print(err)
		api.InternalServerError(w)
		return
	}

	err = h.repo.SaveResponse(context.Background(), res, id)
	if err != nil {
		log.Print(err)
		api.InternalServerError(w)
		return
	}
}

func (h *Handler) HandlerConnect(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)
	conn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println("failed to dial to target", r.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Fatal("http hijacking failed")
	}

	w.WriteHeader(http.StatusOK)

	go tunnelConn(conn, clientConn)
	go tunnelConn(clientConn, conn)
}

func tunnelConn(dst io.WriteCloser, src io.ReadCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}

func copyHeader(from, where http.Header) {
	for key, values := range from {
		for _, v := range values {
			where.Add(key, v)
		}
	}
}
