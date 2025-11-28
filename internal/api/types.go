package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/memartello/go-http-server/internal/database"
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

type UpdateUserRequest NewUserRequest

type UserLoginRequest NewUserRequest
type NewUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
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


type Event string 
const (
	UserUpgraded Event = "user.upgraded"
)


type User struct {
	UserID string `json:"user_id"`
}
type HookEvent struct {
	Event Event `json:"event"`
	Data User `json:"data"`
}

func ConvertToResponseUser(dbUser *database.User) (*NewUserResponse) {
	user := &NewUserResponse{
		ID: dbUser.ID,
		Email: dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	return  user
}