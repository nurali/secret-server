package metric

import (
	"net/http"
	"strings"
	"time"
)

func Recorder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		endpoint := extractEndpoint(r.URL.Path)

		recordRequestCounter(endpoint, method)

		startTime := time.Now()
		next(w, r)
		recordResponseTime(endpoint, method, startTime)
	}
}

func extractEndpoint(path string) string {
	endpoint := strings.TrimPrefix(path, "/api/")
	lastInd := strings.LastIndex(endpoint, "/")
	if lastInd == -1 {
		return path
	}
	return path[:lastInd+5] // len("/api/") = 5
}
