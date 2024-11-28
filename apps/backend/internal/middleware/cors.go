package middleware

import "net/http"

func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-headers", "content-type, authorization, x-requested-with")
		w.Header().Set("access-control-allow-methods", "post, get, options, put, delete")
		w.Header().Set("access-control-allow-origin", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
