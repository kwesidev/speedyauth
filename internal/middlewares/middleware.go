package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/kwesidev/authserver/internal/utilities"
)

// Middleware auth handler
func JwtAuth(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	const ErrorMessageInvalidToken string = "Invalid Token"
	const ErrorMessageProvideValidToken string = "Failed provide a valid token in request header as Token"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from request header
		bearerToken := r.Header.Get("Authorization")
		token := strings.Replace(bearerToken, "Bearer ", "", -1)
		if token != "" {
			claims, err := utilities.ValidateJwtAndGetClaims(token)
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

// Method Middleware to specify the type HTTP method a request can handle
func Method(method string, handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			utilities.JSONError(w, "This Method Not Allowed", http.StatusBadRequest)
			return
		}
		handler(w, r)
	})
}

// // RateLimitter limits the number of request to 10requests and 30 more requests per seconds
// func RateLimiter(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
// 	limiter := rate.NewLimiter(30, 10)
// 	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
// 		if !limiter.Allow() {
// 			utilities.JSONError(w, "Request Exceeded try again later", http.StatusTooManyRequests)
// 			return
// 		}
// 		handler(w, request)
// 	})
// }

// Log all request made
func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
