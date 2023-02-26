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

func NewMiddlewareWithDB(db db.Queries, jwtSecretKey string) DBMiddleware {
	r := DBMiddleware{
		db:           db,
		jwtSecretKey: jwtSecretKey,
	}
	return r
}

type DBMiddleware struct {
	db           db.Queries
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

// AuthorizeUserGroup checks if the user is in the specified group
func (m DBMiddleware) AuthorizeUserGroup(group string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			user := r.Context().Value(User).(db.User)
			if user.Name == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			userInGroup, err := m.db.GetUserIsInGroup(r.Context(), db.GetUserIsInGroupParams{
				UserID:    user.Name,
				GroupName: group,
			})
			if err != nil {
				log.Printf("Unable to check if user is in group: %s", user.Name)
				log.Printf("Error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !userInGroup {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			// pass through
			h.ServeHTTP(w, r)
		})
	}
}
