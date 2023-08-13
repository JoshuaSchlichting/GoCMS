package api

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/data/db"
)

type filesystem interface {
	GetFileContents(path string) ([]byte, error)
	WriteFileContents(path string, contents []byte) error
}

type API struct {
	tmpl       *template.Template
	config     *config.Config
	data       db.DBCache
	fs         filesystem
	router     *chi.Mux
	middleware map[string]func(http.Handler) http.Handler
}

func InitAPI(r *chi.Mux, tmpl *template.Template, config *config.Config, data db.DBCache, fs filesystem) *API {
	api := &API{
		tmpl:   tmpl,
		config: config,
		data:   data,
		fs:     fs,
		router: r,
	}
	api.initPostRoutes()
	api.initGetRoutes()
	return api
}
