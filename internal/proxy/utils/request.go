package utils

import (
	"io"
	"net/http"
	"strings"
)

type RequestInfo struct {
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	Host       string              `json:"host"`
	Headers    map[string][]string `json:"headers"`
	Cookies    map[string]string   `json:"cookies"`
	GetParams  map[string][]string `json:"get_params"`
	PostParams map[string][]string `json:"post_params"`
	Body       string              `json:"body"`
}

func ParseRequest(r *http.Request) *RequestInfo {
	var ri RequestInfo
	ri.Method = r.Method
	ri.Path = r.URL.Path
	ri.Host = r.Host
	h := r.Header
	headers := make(map[string][]string)
	for key, values := range h {
		headers[key] = values
	}
	ri.Headers = headers

	c := r.Cookies()
	cookies := make(map[string]string)
	for _, cookie := range c {
		cookies[cookie.Name] = cookie.Value
	}
	ri.Cookies = cookies

	queryParams := r.URL.Query()
	params := make(map[string][]string)
	for key, values := range queryParams {
		params[key] = values
	}
	ri.GetParams = params

	if err := r.ParseForm(); err == nil {
		params := make(map[string][]string)
		for key, values := range r.PostForm {
			params[key] = values
		}
		ri.PostParams = params
	} else {
		body := &strings.Builder{}
		defer r.Body.Close()
		if _, err := io.Copy(body, r.Body); err == nil {
			ri.Body = body.String()
		}
	}
	return &ri
}
