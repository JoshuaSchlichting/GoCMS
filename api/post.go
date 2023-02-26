package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
)

func InitPostRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, data db.Queries) {
	r.Group(func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			// get payload
			file, y, _ := r.FormFile("file")
			header := y.Header
			log.Printf("header: %v", header)
			// convert file to []byte
			payload := make([]byte, y.Size)
			size, err := file.Read(payload)
			if err != nil {
				log.Printf("error reading file: %v", err)
			}
			log.Printf("file: %v\n\tsize: %v", y.Header, size)
			// data.UploadFile(payload, y.Filename, "userid")
		})

		r.Post("/upload_media", func(w http.ResponseWriter, r *http.Request) {
			// get payload
			file, y, _ := r.FormFile("file")
			header := y.Header
			log.Printf("header: %v", header)
			payload := make([]byte, y.Size)
			size, err := file.Read(payload)
			if err != nil {
				log.Printf("error reading file: %v", err)
			}
			log.Printf("file: %v\n\tsize: %v", y.Header, size)
			// data.UploadFile(payload, y.Filename, "userid")
		})

		r.Put("/api/user", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
			if err != nil {
				log.Printf("error parsing id: %v", err)
			}

			params := db.UpdateUserParams{
				ID:         id,
				Name:       r.FormValue("name"),
				Email:      r.FormValue("email"),
				Attributes: json.RawMessage([]byte(r.FormValue("attributes"))),
			}

			newUser, err := data.UpdateUser(r.Context(), params)
			if err != nil {
				log.Printf("error updating user: %v", err)
			}
			log.Println(newUser)
		})
	})
}
