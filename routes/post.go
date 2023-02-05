package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/data"
)

func InitPostRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, data data.Data) {
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
			data.UploadFile(payload, y.Filename, "userid")

			// payload := r.Context().Value("payload").(map[string]interface{})
		})
	})
}
