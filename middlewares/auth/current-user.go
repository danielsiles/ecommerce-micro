package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
)

var (
	SessionStore = sessions.NewCookieStore([]byte("SESSION_SECRET"))
)

// MiddlewareCurrentUser return the current user logged in the session.
func MiddlewareCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		session, err := SessionStore.Get(r, "SESSION_SECRET")
		if err != nil || session == nil {
			fmt.Println("Session expired")
		}

		fmt.Println("TESTE2")

		tokenValue, ok := session.Values["token"]
		if !ok {
			fmt.Println("Session expired")
			return
		}

		token, err := jwt.Parse(tokenValue.(string), func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("secret"), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("Invalid token")
			return
		}
		// add the product to the context
		ctx := context.WithValue(r.Context(), "currentUser", claims)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
