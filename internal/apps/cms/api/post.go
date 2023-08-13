package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/internal/data/db"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
}

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
				logger.Error(fmt.Sprintf("error reading file: %v", err))
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
				logger.Error(fmt.Sprintf("error reading file: %v", err))
			}
			log.Printf("file: %v\n\tsize: %v", y.Header, size)
			// data.UploadFile(payload, y.Filename, "userid")
		})

		r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			params := db.CreateUserParams{
				ID:         uuid.New(),
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUser, err := a.data.CreateUser(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error creating user: %v", err))
			}
			log.Println(newUser)
		})

		r.Put("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			params := db.UpdateUserParams{
				ID:         id,
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUser, err := a.data.UpdateUser(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error updating user: %v", err))
			}
			log.Println(newUser)
		})

		r.Delete("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			// log url params
			log.Printf("url params: %v", chi.URLParam(r, "id"))
			// get id from url param
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			err = a.data.DeleteUser(r.Context(), id)
			if err != nil {
				logger.Error(fmt.Sprintf("error deleting user: %v", err))
			}
		})

		r.Post("/organization", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			params := db.CreateOrganizationParams{
				ID:         uuid.New(),
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newOrganization, err := a.data.CreateOrganization(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error creating organization: %v", err))
			}
			log.Println(newOrganization)
		})

		r.Put("/organization/{id}", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			params := db.UpdateOrganizationParams{
				ID:         id,
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newOrganization, err := a.data.UpdateOrganization(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error updating organization: %v", err))
			}
			log.Println(newOrganization)
		})

		r.Delete("/organization/{id}", func(w http.ResponseWriter, r *http.Request) {
			// log url params
			log.Printf("url params: %v", chi.URLParam(r, "id"))
			// get id from url param
			id, err := uuid.Parse((chi.URLParam(r, "id")))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			err = a.data.DeleteOrganization(r.Context(), id)
			if err != nil {
				logger.Error(fmt.Sprintf("error deleting organization: %v", err))
			}
		})
		r.Post("/usergroup", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			params := db.CreateUserGroupParams{
				ID:         uuid.New(),
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUserGroup, err := a.data.CreateUserGroup(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error creating user group: %v", err))
			}
			log.Println(newUserGroup)
		})

		r.Put("/usergroup/{id}", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			log.Printf("form: %v", r.Form)
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			params := db.UpdateUserGroupParams{
				ID:         id,
				Name:       r.FormValue("Name"),
				Email:      r.FormValue("Email"),
				Attributes: json.RawMessage([]byte(r.FormValue("Attributes"))),
			}

			newUserGroup, err := a.data.UpdateUserGroup(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error updating user group: %v", err))
			}
			log.Println(newUserGroup)
		})
		r.Delete("/usergroup/{id}", func(w http.ResponseWriter, r *http.Request) {
			// log url params
			log.Printf("url params: %v", chi.URLParam(r, "id"))
			// get id from url param
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing id: %v", err))
			}

			err = a.data.DeleteUserGroup(r.Context(), id)
			if err != nil {
				logger.Error(fmt.Sprintf("error deleting user group: %v", err))
			}
		})

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
		logger.Error(fmt.Sprintf("error reading file: %v", err))
	}
	log.Printf("file: %v\n\tsize: %v", y.Header, size)
	err = a.fs.WriteFileContents(y.Filename, payload)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing file: %v", err))
	}
}
