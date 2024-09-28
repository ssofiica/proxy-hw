package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	req, err := h.repo.GetRequestList(context.Background())
	if err != nil {
		InternalServerError(w)
		return
	}
	response, err := json.Marshal(req)
	if err != nil {
		InternalServerError(w)
		return
	}
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
	res, err := h.repo.GetRequestByID(context.Background(), id)
	if err != nil {
		log.Print(err)
		return
	}
	w.Write(res)
}

func (h *Handler) Repeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Print("no id")
		return
	}
}

func (h *Handler) Scan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Print("no id")
		return
	}
}
