package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)


type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		cfg.fileServerHits.Add(1)
		log.Println(cfg.getHits())
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) getHits() int {
	hitsMessage := int(cfg.fileServerHits.Load())
	return hitsMessage
}

func (cfg *apiConfig) resetHits() {
	cfg.fileServerHits.Swap(0)
}

var config *apiConfig

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func noCacheHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Cache-Control, Pragma, and Expires headers to prevent caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0") // Or a date in the past, like time.Unix(0, 0).Format(time.RFC1123)

		// Serve the actual file
		h.ServeHTTP(w, r)
	})
}

func main(){
	server_mux := http.NewServeMux()
 	server := &http.Server{
		Addr: ":8080",
		Handler: server_mux,
 	}
	config = &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	wrappedHandler := middlewareLog(noCacheHandler(config.middlewareMetricsInc(http.StripPrefix("/app/",http.FileServer(http.Dir("./"))))))
	
	server_mux.Handle("/app/",wrappedHandler)

	server_mux.HandleFunc("GET /api/healthz",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		
	})

	server_mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		strResponse := fmt.Sprintf(`<html>
		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>`, config.getHits())
		w.Write([]byte(strResponse))
	})

	server_mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		config.resetHits()
	})
	err := server.ListenAndServe()

	defer server.Close()

	if err != nil {
		log.Fatal("Error initializing the server" + err.Error())
		os.Exit(1)
	}

}