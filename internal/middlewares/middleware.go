package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/kwesidev/authserver/internal/utilities"
)

// Middleware auth handler
func JwtAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	const ErrorMessageInvalidToken string = "Invalid Token"
	const ErrorMessageProvideValidToken string = "Failed provide a valid token in request header as Token"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from request header
		providedToken := r.Header.Get("token")
		if providedToken != "" {
			claims, err := utilities.ValidateJwtAndGetClaims(providedToken)
			if err != nil {
				utilities.JSONError(w, ErrorMessageInvalidToken, http.StatusForbidden)
				log.Println(ErrorMessageInvalidToken)
				return
			}
			ctx := context.WithValue(r.Context(), "claims", claims)
			handler(w, r.WithContext(ctx))
		} else {
			utilities.JSONError(w, ErrorMessageProvideValidToken, http.StatusForbidden)
			log.Println(ErrorMessageProvideValidToken)
			return
		}
	})
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
