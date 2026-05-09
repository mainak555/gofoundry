package handlers

import (
	"net/http"
	"strings"
)

func HealthProbeHandler(appName string, _callBack func(w http.ResponseWriter, r *http.Request)) func(next http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == "GET" || r.Method == "HEAD") &&
				(strings.EqualFold(r.URL.Path, "/") || strings.EqualFold(r.URL.Path, "/health") ||
					strings.EqualFold(r.URL.Path, "/health/live") || strings.EqualFold(r.URL.Path, "/health/ready")) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(appName))
				if _callBack != nil {
					_callBack(w, r)
				}
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
