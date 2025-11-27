package api

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NoCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Cache-Control, Pragma, and Expires headers to prevent caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0") // Or a date in the past, like time.Unix(0, 0).Format(time.RFC1123)

		// Serve the actual file
		h.ServeHTTP(w, r)
	})
}


func (api *API) MetricsIncMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		api.fileServerHits.Add(1)
		log.Println(api.getHits())
		next.ServeHTTP(w,r)
	})
}