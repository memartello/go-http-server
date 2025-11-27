package api

import (
	"log"
	"net/http"

	"github.com/memartello/go-http-server/internal/auth"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NoCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Cache-Control, Pragma, and Expires headers to prevent caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0") // Or a date in the past, like time.Unix(0, 0).Format(time.RFC1123)

		// Serve the actual file
		h.ServeHTTP(w, r)
	})
}


func (api *API) MetricsIncMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		api.fileServerHits.Add(1)
		log.Println(api.getHits())
		next.ServeHTTP(w,r)
	})
}


// type ctxKey string

// const userCtxKey ctxKey = "user"
// func UserFromContext(ctx context.Context) (string, bool) {
//     userID, ok := ctx.Value(userCtxKey).(string)
//     return userID, ok
// }
// TODO Save user in the context.auth.UserFromContext(r.Context())
func (api *API) AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		jwt, err := auth.GetBearerToken(r.Header)

		if err != nil {
			RespondWithError(w, http.StatusUnauthorized ,"No authorization header is present")
			return
		}

		_, err = auth.ValidateJWT(jwt, api.secret)

		if err != nil {
			RespondWithError(w, http.StatusUnauthorized ,"Invalid JWT")
			return
		}
		//TODO: Check user in DB and send to
		
		next.ServeHTTP(w,r)
	})
}