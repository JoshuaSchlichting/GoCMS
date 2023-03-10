package api

import (
	"html/template"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
)

type filesystem interface {
	GetFileContents(path string) ([]byte, error)
	WriteFileContents(path string, contents []byte) error
}

type API struct {
	tmpl   *template.Template
	config *config.Config
	data   db.Queries
	fs     filesystem
	router *chi.Mux
}

func InitAPI(r *chi.Mux, tmpl *template.Template, config *config.Config, data db.Queries, fs filesystem) *API {
	api := &API{
		tmpl:   tmpl,
		config: config,
		data:   data,
		fs:     fs,
		router: r,
	}
	api.initPostRoutes()
	return api
}
