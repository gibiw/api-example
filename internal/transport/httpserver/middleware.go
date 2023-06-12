package httpserver

import (
	"net/http"
)

func setResponseHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodPost, http.MethodPut:
				w.Header().Add("Content-type", "application/json")
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
