package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	api_service "github.com/memartello/go-http-server/internal/api"
	"github.com/memartello/go-http-server/internal/database"
)


func main(){
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	token_env := os.Getenv("JWT_TOKEN")
	polka_key := os.Getenv("POLKA_KEY")
	fmt.Println(dbUrl)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database connection: " + err.Error())
	}
	defer db.Close()
	dbQueries := database.New(db)
	
	
	api := api_service.NewAPI(dbQueries, token_env, polka_key)
	
	server_mux := http.NewServeMux()
 	server := &http.Server{
		Addr: ":8080",
		Handler: server_mux,
 	}

	wrappedHandler := api_service.LogMiddleware(api_service.NoCacheMiddleware(api.MetricsIncMiddleware(http.StripPrefix("/app/",http.FileServer(http.Dir("./"))))))
	

	// Static files
	server_mux.Handle("/app/",wrappedHandler)
	
	// Admin
	server_mux.HandleFunc("GET /admin/metrics", api.GetMetrics)
	if os.Getenv("PLATFORM") == "dev" {
		server_mux.HandleFunc("POST /admin/reset", api.Reset)
	}
	
	// Api
	server_mux.HandleFunc("GET /api/healthz",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		
	})
	server_mux.HandleFunc("POST /api/users", api.CreateUser)
	server_mux.Handle("PUT /api/users", api.AuthMiddleware(http.HandlerFunc(api.UpdateUser)))
	server_mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)
	server_mux.Handle("POST /api/chirps", api.AuthMiddleware(http.HandlerFunc(api.CreateChirp)))
	server_mux.Handle("DELETE /api/chirps/{chirpID}", api.AuthMiddleware(http.HandlerFunc(api.DeleteChirp)))
	server_mux.HandleFunc("GET /api/chirps", api.GetChirps)
	server_mux.HandleFunc("GET /api/chirps/{chirpID}", api.GetChirp)
	server_mux.HandleFunc("POST /api/login", api.Login)
	server_mux.HandleFunc("POST /api/refresh", api.Refresh)
	server_mux.Handle("POST /api/revoke", api.AuthMiddleware(http.HandlerFunc(api.Revoke)))
	server_mux.HandleFunc("POST /api/polka/webhooks", api.UpgradeUser)
	
	
	server.ListenAndServe()

	defer server.Close()

	if err != nil {
		log.Fatal("Error initializing the server" + err.Error())
		os.Exit(1)
	}

}