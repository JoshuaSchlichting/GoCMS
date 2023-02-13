package api

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
	"github.com/joshuaschlichting/gocms/middleware"
)

func InitGetRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.QueriesInterface, middlewareMap map[string]func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		jwtAuth := jwtauth.New("HS256", []byte(config.Auth.JWT.SecretKey), nil)
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(middleware.AddClientJWTStringToCtx)
		r.Use(middleware.AuthenticateJWT)
		r.Use(middlewareMap["addUserToCtx"])

	})

}
