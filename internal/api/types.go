package api

import (
	"time"

	"github.com/google/uuid"
)

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ValidateChirpRequest struct {
	Body string `json:"body"`
}

type TokenRespone struct {
	Token string `json:"token"`
}

type NewUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest NewUserRequest
type NewUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserResponse struct {
	NewUserResponse
	Token string `json:"token"`	
	RefreshToken string `json:"refresh_token"`
}

type ErrorMessage struct{
	Error string `json:"error"`
}


type NewChirpRequest struct {
	Body string `json:"body"`
}

type NewChirpResponse struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID string `json:"user_id"`
}

type ctxKey string

const userCtxKey ctxKey = "user_uuid"