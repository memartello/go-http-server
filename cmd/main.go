package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

type Parameters struct {
	Body string `json:"body"`
}

type ErrorMessage struct{
	Error string `json:"error"`
}

type Response struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	msgError := ErrorMessage{
		Error: msg,
	}
	dat, _ := json.Marshal(msgError)
	
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data, _ := json.Marshal(payload)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

var profaneWords = []string{"kerfuffle","sharbert","fornax"}

func cleanedString (str string) string{
	str_list := strings.Split(str, " ")
	for i, v := range str_list{
		for _, prof := range profaneWords{
			if strings.ToLower(v) == prof {
				str_list[i] = "****"
				break
			}
		}
	}

	str = strings.Join(str_list, " ")

	return  str
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

	
	server_mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {

		var parameters Parameters

		w.Header().Set("Content-Type","application/json")
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&parameters); err != nil{
			log.Printf("Error decoding parameters: %s",err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		if parameters.Body == "" {
			respondWithError(w, http.StatusBadRequest, "You need to specify the body.")
			return
		}

		if len(parameters.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long.")
			return
		}
		
		res := Response{
			CleanedBody: cleanedString(parameters.Body),
		}

		respondWithJSON(w, http.StatusOK, res)
	})
	err := server.ListenAndServe()

	defer server.Close()

	if err != nil {
		log.Fatal("Error initializing the server" + err.Error())
		os.Exit(1)
	}

}