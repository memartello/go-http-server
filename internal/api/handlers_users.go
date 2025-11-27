package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) {
		var user NewUserRequest
		
		w.Header().Set("Content-Type","application/json")
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&user); err != nil{
			log.Printf("Error decoding parameters: %s",err)
			RespondWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		
		db_user, err := api.dbQueries.CreateUser(r.Context(), user.Email)
		
		user_created := NewUserResponse{
			ID: db_user.ID,
			CreatedAt: db_user.CreatedAt.Time,
			UpdatedAt: db_user.UpdatedAt.Time,
			Email: db_user.Email,
		}
		if err != nil {
			log.Printf("Error creating user: %s",err)
			RespondWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}
		
		RespondWithJSON(w, http.StatusCreated, user_created)
	}