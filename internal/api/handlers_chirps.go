package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/memartello/go-http-server/internal/database"
)

func (api *API) ValidateChirp(w http.ResponseWriter, r *http.Request) {
		var parameters ValidateChirpRequest

		w.Header().Set("Content-Type","application/json")
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&parameters); err != nil{
			log.Printf("Error decoding parameters: %s",err)
			RespondWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		if parameters.Body == "" {
			RespondWithError(w, http.StatusBadRequest, "You need to specify the body.")
			return
		}

		if len(parameters.Body) > 140 {
			RespondWithError(w, http.StatusBadRequest, "Chirp is too long.")
			return
		}
		
		res := ValidateChirpResponse{
			CleanedBody: CleanedString(parameters.Body),
		}

		RespondWithJSON(w, http.StatusOK, res)
}

func (api *API) GetChirps(w http.ResponseWriter, r *http.Request){
	
	w.Header().Set("Content-Type","application/json")

	chirps, err := api.dbQueries.GetChirps(r.Context())

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response_chirps := make([]NewChirpResponse, 0, len(chirps))

	for _, chirp := range chirps {
		response_chirps = append(
			response_chirps,
			NewChirpResponse{
				Body: chirp.Body,
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt.Time,
				UpdatedAt: chirp.UpdatedAt.Time,
				UserID: chirp.UserID.String(),
			},
		)
	}

	RespondWithJSON(w, http.StatusOK, response_chirps)
}

func (api *API) CreateChirp(w http.ResponseWriter, r *http.Request){
	var parameters NewChirpRequest

	w.Header().Set("Content-Type","application/json")

	user_uuid, ok  := UserFromContext(r.Context())
	fmt.Printf("User uuid %s \n", user_uuid)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return;
	}

	
	decoder := json.NewDecoder(r.Body)
	
	if err := decoder.Decode(&parameters); err != nil {
		log.Printf("Error decoding parameters: %s",err)
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong!")
		return
	}

	if parameters.Body == "" {
		RespondWithError(w, http.StatusBadRequest,  "Body can not be empty")
		return
	}

	_, err := api.dbQueries.GetUser(r.Context(), uuid.MustParse(user_uuid))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "User not found")
		return
	}

	newChirp, err := api.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: uuid.MustParse(user_uuid),
		Body: parameters.Body,
	})

	response_chirp := NewChirpResponse{
		UserID: newChirp.UserID.String(),
		Body: newChirp.Body,
		CreatedAt: newChirp.CreatedAt.Time,
		UpdatedAt: newChirp.UpdatedAt.Time,
		ID: newChirp.ID,
	}
	
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,"Error creatining chirp :" + err.Error())
		return
	}

	RespondWithJSON(w,http.StatusCreated, response_chirp)
}

func (api *API) GetChirp(w http.ResponseWriter, r *http.Request){
	value := r.PathValue("chirpID")

	if value == "" {
		RespondWithError(w, http.StatusBadRequest, "Needs to specify the Chirp ID")
		return
	}

	db_chirp, err := api.dbQueries.GetChirp(r.Context(), uuid.MustParse(value))

	if err != nil {
		log.Printf("Error %v \n", err)
		RespondWithError(w, http.StatusNotFound, "Resource not found")
		return
	}

	chirpResponse := NewChirpResponse{
		UserID: db_chirp.UserID.String(),
		ID: db_chirp.ID,
		CreatedAt: db_chirp.CreatedAt.Time,
		UpdatedAt: db_chirp.UpdatedAt.Time,
		Body: db_chirp.Body,
	}

	RespondWithJSON(w, http.StatusOK, chirpResponse)
}