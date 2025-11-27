package api

import (
	"fmt"
	"log"
	"net/http"
)


func (api *API) resetHits() {
	api.fileServerHits.Swap(0)
}
func (api *API) getHits() int {
	hitsMessage := int(api.fileServerHits.Load())
	return hitsMessage
}


func (api *API) Reset(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type","text/plain; charset=utf-8")
			err := api.dbQueries.DeleteUsers(r.Context())
			if err != nil {
				log.Printf("Error deleting users: %s",err)
				RespondWithError(w, http.StatusInternalServerError, "Something went wrong!")
				return
			}
			w.WriteHeader(http.StatusOK)
			api.resetHits()
		}

func (api *API) GetMetrics(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		strResponse := fmt.Sprintf(`<html>
		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>`, api.getHits())
		w.Write([]byte(strResponse))
	}