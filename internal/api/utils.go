package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func RespondWithError(w http.ResponseWriter, code int, msg string){
	msgError := ErrorMessage{
		Error: msg,
	}
	dat, _ := json.Marshal(msgError)
	
	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Write(data)
}


var profaneWords = []string{"kerfuffle","sharbert","fornax"}

func CleanedString (str string) string{
	str_list := strings.Split(str, " ")
	for i, v := range str_list{
		for _, prof := range profaneWords{
			if strings.ToLower(v) == prof {
				str_list[i] = "****"
				break
			}
		}
	}

	str = strings.Join(str_list, " ")

	return  str
}

func UserFromContext(ctx context.Context) (string, bool) {
	
	userID, ok := ctx.Value(userCtxKey).(string)
	fmt.Printf("Gettin value from context %s \n",userID)
    return userID, ok
}