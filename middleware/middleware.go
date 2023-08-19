package middleware

import (
	"log"
	"net/http"
	"strings"
)

func LogAllButStaticRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/static/") {
			h.ServeHTTP(w, r)
			return
		}
		log.Printf("%s %s %s FROM %s", r.Method, r.RequestURI, r.Proto, r.RemoteAddr)
		log.Printf("Cookies: %v", r.Cookies())
		h.ServeHTTP(w, r)
	})
}
