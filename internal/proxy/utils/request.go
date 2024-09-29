package utils

import (
	"fmt"
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

func MakeRequest(r *RequestInfo) (*http.Request, error) {
	var body io.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(
		r.Method,
		fmt.Sprintf("http://%s%s", r.Host, r.Path),
		body,
	)
	if err != nil {
		return nil, err
	}

	for key, values := range r.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	for name, val := range r.Cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: val})
	}

	query := req.URL.Query()
	for key, values := range r.GetParams {
		for _, val := range values {
			query.Add(key, val)
		}
	}
	req.URL.RawQuery = query.Encode()

	for key, values := range r.PostParams {
		for _, value := range values {
			req.PostForm.Add(key, value)
		}
	}

	return req, nil
}
