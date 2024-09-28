package utils

import (
	"io"
	"net/http"
	"strings"
)

type ResponseInfo struct {
	Code    int                 `json:"code"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func ParseResponse(r *http.Response) *ResponseInfo {
	var ri ResponseInfo
	ri.Code = r.StatusCode
	h := r.Header
	headers := make(map[string][]string)
	for key, values := range h {
		headers[key] = values
	}
	ri.Headers = headers
	body := &strings.Builder{}
	defer r.Body.Close()
	if _, err := io.Copy(body, r.Body); err == nil {
		ri.Body = body.String()
	}
	return &ri
}
