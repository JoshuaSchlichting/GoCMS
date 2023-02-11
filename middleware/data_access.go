package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/joshuaschlichting/gocms/db"
)

type MiddlewareWithDB interface {
	AddUserToCtx(h http.Handler) http.Handler
}

func NewMiddlewareWithDB(db db.QueriesInterface, jwtSecretKey string) DBMiddleware {
	r := DBMiddleware{
		db:           db,
		jwtSecretKey: jwtSecretKey,
	}
	return r
}

type DBMiddleware struct {
	db           db.QueriesInterface
	jwtSecretKey string
}

func (m DBMiddleware) AddUserToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Context().Value(JWTEncodedString).(string)
		// unmarshal jwt

		// cast to *jwtauth.JWTAuth
		tokenAuth := jwtauth.New("HS256", []byte(m.jwtSecretKey), nil)
		jwtAuthToken, err := tokenAuth.Decode(jwt)
		if err != nil {
			log.Printf("Unable to decode JWT: %s", jwt)
			log.Printf("Error: %s", err.Error())
			h.ServeHTTP(w, r)
			return
		}
		username, _ := jwtAuthToken.Get("userInfo")
		user, err := m.db.GetUserByName(r.Context(), username.(string))
		if err != nil {
			log.Printf("Unable to get user from db: %s", username)
			log.Printf("Error: %s", err.Error())
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), User, user)))
	})
}
