package middleware

import (
	"net/http"

	"github.com/joshuaschlichting/gocms/db"
)

type MiddlewareWithDB interface {
	AddUserToCtx(h http.Handler) http.Handler
}

func New(db db.QueriesInterface) DBMiddleware {
	r := DBMiddleware{
		DB: db,
	}
	return r
}

type DBMiddleware struct {
	DB db.QueriesInterface
}

func (m DBMiddleware) AddUserToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user from db
		m.DB.GetUser(r.Context(), 1)
		h.ServeHTTP(w, r)
	})
}
