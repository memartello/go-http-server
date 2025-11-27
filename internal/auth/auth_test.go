package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestNewJWT(t *testing.T){
	secret := "AllYourBase"
	ttl := time.Minute * 1
	userId := uuid.MustParse("03449a98-b974-4d61-bded-a5689c04c17d")

	tokenStr, err := MakeJWT(userId,secret, ttl)

	if err  != nil {
		t.Fatalf("makeJWT returned an error: %v", err)
	}

	if tokenStr == "" {
		t.Fatalf("make jwt return empty tokenstr")
	}	

	claim_uuid, err := ValidateJWT(tokenStr, secret)

	if err != nil {
		t.Fatalf("Validate jwt throws an error %v",err)
	}

	if claim_uuid != userId{
		t.Fatalf("Expected %v, got %v", userId, claim_uuid)
	}

}

//TODO Add functionality to this test
func TestValdiateJWT(t *testing.T){

}


func TestAuthorizationHeader(t *testing.T){
	mocked_header := http.Header{}

	mocked_header.Set("Authorization","Bearer dindare")

	auth_string, err := GetBearerToken(mocked_header)

	if err != nil {
		t.Fatalf("An error ocurred %v \n", err)
	}
	if (auth_string != "dindare"){
		t.Fatalf("Value of header dosen't match mocked header")
	}

	mocked_header.Del("Authorization")

	_, err =  GetBearerToken(mocked_header)

	if err == nil {
		t.Fatalf("An error is expected when Authorization is not setted")
	}
}


func TestRefreshToken(t *testing.T){
	_, err := MakeRefreshToken()
	
	if err != nil {
		t.Fatalf("Token was not received.")
	}
}