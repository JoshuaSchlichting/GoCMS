package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/joshuaschlichting/gocms/data/db"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
}

type MiddlewareWithDB interface {
	AddUserToCtx(h http.Handler) http.Handler
}

func NewMiddlewareWithDB(db db.DBCache, jwtSecretKey string) DBMiddleware {
	r := DBMiddleware{
		db:           db,
		jwtSecretKey: jwtSecretKey,
	}
	return r
}

type DBMiddleware struct {
	db           db.DBCache
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
			logger.Error(fmt.Sprintf("error decoding JWT from context: %s: JWT from context: %v", err.Error(), jwt))
			h.ServeHTTP(w, r)
			return
		}
		username, _ := jwtAuthToken.Get("userInfo")
		user, err := m.db.GetUserByName(r.Context(), username.(string))
		if err != nil {
			log.Printf("unable to get user '%s' from db due to error: %v", username, err.Error())
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
				UserID:        user.ID,
				UsergroupName: group,
			})
			if err != nil {
				logger.Error("unable to check if user '%s' is in group '%s' due to error: %v", user.Name, group, err.Error())
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
