package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joshuaschlichting/gocms/middleware"
)

// func (a *API) initGetRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.Queries, middlewareMap map[string]func(http.Handler) http.Handler) {
func (a *API) initGetRoutes() {

	a.router.Group(func(r chi.Router) {
		jwtAuth := jwtauth.New("HS256", []byte(a.config.Auth.JWT.SecretKey), nil)
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(middleware.AddClientJWTStringToCtx)
		r.Use(middleware.AuthenticateJWT)
		r.Use(a.middleware["addUserToCtx"])

	})

}
