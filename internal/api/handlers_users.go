package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/memartello/go-http-server/internal/auth"
	"github.com/memartello/go-http-server/internal/database"
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

		hashed_password, _ := auth.HashPassword(user.Password)
		
		db_user, err := api.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
			Email: user.Email,
			HashedPassword: hashed_password,
		})

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

func (api *API) Login(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type","application/json")

		var parameters UserLoginRequest

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&parameters); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Error in the params for login")
			return
		}

		//TODO: Password Policy

		db_user, err := api.dbQueries.GetByEmail(r.Context(), parameters.Email)

		if err != nil {
			RespondWithError(w, http.StatusNotFound, "User not found.")
			return
		}

		match, err := auth.CheckPassword(parameters.Password, db_user.HashedPassword)
		if (err != nil || !match){
			RespondWithError(w, http.StatusUnauthorized, "Password is incorrect.")
			return
		}

		refresh_token, _ := auth.MakeRefreshToken()

		
		expiresAt := time.Now().Add(24 * 60 * time.Hour) 

		log.Printf("Expires at: ", expiresAt)
		api.dbQueries.CreateRefreshToken(r.Context(),database.CreateRefreshTokenParams{
			Token: refresh_token,
			UserID: db_user.ID,
			ExpiresAt: sql.NullTime{
				Time: expiresAt,
				Valid: true,
			},
		})

		token, err := auth.MakeJWT(db_user.ID,api.secret, time.Second * 60)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error ocurred on login")
			return
		}
		user := UserResponse{
			NewUserResponse{
				ID: db_user.ID,
				CreatedAt: db_user.CreatedAt.Time,
				UpdatedAt: db_user.UpdatedAt.Time,
				Email: db_user.Email,
			},
			token,
			refresh_token,
		}

		RespondWithJSON(w, http.StatusOK, user)

}

func (api *API) Refresh(w http.ResponseWriter, r *http.Request){
	jwt, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized ,"No authorization header is present")
		return
	}
	
	db_user_token, err  := api.dbQueries.GetUserByToken(r.Context(), jwt)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Wrong credenitals")
		return
	}

	if db_user_token.ExpiresAt.Time.Before(time.Now())  {
		RespondWithError(w, http.StatusUnauthorized, "Token Expired")
		return
	}

	if db_user_token.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Token Revoked")
		return
	}

	new_jwt, _ := auth.MakeJWT(db_user_token.UserID, api.secret, time.Minute * 60)

	RespondWithJSON(w, http.StatusOK, TokenRespone{
		Token: new_jwt,
	})

}
func (api *API) Revoke(w http.ResponseWriter, r *http.Request){
	jwt, err := auth.GetBearerToken(r.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized ,"No authorization header is present")
		return
	}
	
	_, err  = api.dbQueries.GetUserByToken(r.Context(), jwt)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Wrong credenitals")
		return
	}

	err = api.dbQueries.RevokeToken(r.Context(), jwt)

	if err != nil {
		log.Print(err)
		RespondWithError(w, http.StatusInternalServerError, "An error ocurred updating Refresh token")
		return
	}


	RespondWithJSON(w, http.StatusNoContent, TokenRespone{
		Token: jwt,
	})
}