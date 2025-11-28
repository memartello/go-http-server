package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error){
	hashed_password, err :=argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		return "", err
	}

	return hashed_password, nil
}

func CheckPassword(password, hash string) (bool, error){
	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false, err
	}
	return  match, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	ss, err := token.SignedString(mySigningKey)

	 if err != nil {
		return "", err
	 }
	 return  ss, nil
}

func ValidateJWT (tokenString, tokenSecret string) (uuid.UUID, error){

	claims:= &jwt.RegisteredClaims{}

	_,err := jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	
	return  uuid.MustParse(claims.Subject), nil
}


func GetBearerToken (headers http.Header) (string, error){
	//TODO: Check that is Bearer
	authorization_header := headers.Get("Authorization")

	if authorization_header == "" {
		return "", fmt.Errorf("authorization header is not present")
	}

	stripped_header := strings.Split(authorization_header, " ")[1]

	return stripped_header, nil
}


func MakeRefreshToken() (string, error ){
	key := make([]byte, 32)
	rand.Read(key)

	refresh := hex.EncodeToString(key)
	
	return  refresh, nil
}


func GetAPIKey(headers http.Header) (string, error){
	authorization_header := headers.Get("Authorization")
	//TODO: Check that is apiKey

	if authorization_header == "" {
		return "", fmt.Errorf("authorization header is not present")
	}

	stripped_header := strings.Split(authorization_header, " ")[1]

	return stripped_header, nil
}