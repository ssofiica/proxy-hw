package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ssofiica/proxy-hw/internal/proxy"
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

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	req, err := h.repo.GetRequestList(r.Context())
	if err != nil {
		InternalServerError(w)
		return
	}
	response, err := json.Marshal(req)
	if err != nil {
		InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	if param == "" {
		log.Print("no id")
		BadRequest(w)
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Print("invalid id")
		BadRequest(w)
		return
	}
	res, err := h.repo.GetRequestByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *Handler) Repeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	if param == "" {
		log.Print("no id")
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Print("invalid id")
		BadRequest(w)
		return
	}
	res, err := h.repo.GetRequestByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	var reqInfo utils.RequestInfo
	err = json.Unmarshal(res, &reqInfo)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	req, err := utils.MakeRequest(&reqInfo)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	response, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	io.Copy(w, response.Body)
	if err = response.Body.Close(); err != nil {
		log.Print(err)
	}
	proxy.CopyHeader(response.Header, w.Header())
}

func (h *Handler) Scan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	if param == "" {
		log.Print("no id")
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Print("invalid id")
		BadRequest(w)
		return
	}
	res, err := h.repo.GetRequestByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}
	var reqInfo utils.RequestInfo
	err = json.Unmarshal(res, &reqInfo)
	if err != nil {
		log.Print(err)
		InternalServerError(w)
		return
	}

	payloads := []string{
		";cat /etc/passwd;",
		"|cat /etc/passwd|",
		"`cat /etc/passwd`",
	}
	for _, payload := range payloads {
		req := changeParams(reqInfo, payload)
		request, err := utils.MakeRequest(&req)
		if err != nil {
			log.Print(err)
			InternalServerError(w)
			return
		}
		response, err := http.DefaultTransport.RoundTrip(request)
		if err != nil {
			log.Print(err)
			InternalServerError(w)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
			InternalServerError(w)
			return
		}

		if strings.Contains(string(body), "root:") {
			w.Header().Set("Content-Type", "application/json")
			str := fmt.Sprintf("Данный %s-параметр уязвим", reqInfo.Method)
			w.Write([]byte(str))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Уязвимостей нет"))
}

func changeParams(req utils.RequestInfo, payload string) utils.RequestInfo {
	if req.Method == "GET" {
		for key, values := range req.GetParams {
			if len(req.GetParams[key]) != 0 {
				req.GetParams[key][0] = values[0] + payload
			}
		}
	} else if req.Method == "POST" {
		for key, values := range req.PostParams {
			if len(req.PostParams[key]) != 0 {
				req.PostParams[key][0] = values[0] + payload
			}
		}
	}
	return req
}
