package middleware

import (
	"fmt"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	nh := func(w http.ResponseWriter, req *http.Request) {
		isSkip := SkipFavicon(w, req)

		if isSkip {
			return
		}

		next.ServeHTTP(w, req)
		LogRequest(w, req)
	}

	return http.HandlerFunc(nh)
}

func SkipFavicon(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/favicon.ico" {
		return true
	}
	return false
}

func LogRequest(w http.ResponseWriter, req *http.Request) {
	request := fmt.Sprintf("%s %s", req.URL.String(), req.Method)
	fmt.Println(request)
}
