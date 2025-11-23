package main

import (
	"log"
	"net/http"
	"os"
)



func main(){
	server_mux := http.NewServeMux()
 	server := &http.Server{
		Addr: ":8080",
		Handler: server_mux,
 	}

	server_mux.Handle("/app/",http.StripPrefix("/app/",http.FileServer(http.Dir("./"))))

	server_mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		
	})
	err := server.ListenAndServe()

	defer server.Close()

	if err != nil {
		log.Fatal("Error initializing the server" + err.Error())
		os.Exit(1)
	}

}