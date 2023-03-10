package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/db"
)

func (a *API) initPostRoutes() {
	a.router.Route("/api", func(r chi.Router) {
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

		r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			params := db.CreateUserParams{
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUser, err := a.data.CreateUser(r.Context(), params)
			if err != nil {
				log.Printf("error creating user: %v", err)
			}
			log.Println(newUser)
		})

		r.Put("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				log.Printf("error parsing id: %v", err)
			}

			params := db.UpdateUserParams{
				ID:         id,
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUser, err := a.data.UpdateUser(r.Context(), params)
			if err != nil {
				log.Printf("error updating user: %v", err)
			}
			log.Println(newUser)
		})

		r.Delete("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			// log url params
			log.Printf("url params: %v", chi.URLParam(r, "id"))
			// get id from url param
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				log.Printf("error parsing id: %v", err)
			}

			err = a.data.DeleteUser(r.Context(), id)
			if err != nil {
				log.Printf("error deleting user: %v", err)
			}
		})

		r.Post("/upload_file", a.UploadFileHandler)
	})
}

// UploadFileHandler handles the upload of a file
func (a *API) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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
	err = a.fs.WriteFileContents(y.Filename, payload)
	if err != nil {
		log.Printf("error writing file: %v", err)
	}
}
