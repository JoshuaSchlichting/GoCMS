package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/auth/kratos"
	auth "github.com/joshuaschlichting/gocms/auth/oauth2"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/admin/components"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/data/db"
)

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
}
func InitRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.DBCache, middlewareMap map[string]func(http.Handler) http.Handler) {
	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		// if not logged in
		if r.Context().Value(User) == nil {
			// redirect to login
			// set Content-Type header
			w.Header().Set("Content-Type", "")
			logger.Debug("redirecting to login: " + config.Auth.SignInUrl)
			http.Redirect(w, r, config.Auth.SignInUrl, http.StatusFound)
			return
		}

		http.Redirect(w, r, "/secure", http.StatusFound)

	})

	r.Get("/loggedout", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "logged_out", nil)
	})

	r.Get("/getjwtandlogin", func(w http.ResponseWriter, r *http.Request) {

		var jwtTokenS string
		var err error

		var username string
		var email string
		// get access code from request
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Println("no access code found in request URL params")
		} else {
			// get JWT
			jwtTokenS, err = auth.GetAccessJWT(code)
			if err != nil {
				logger.Debug("", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			username, email, err = auth.GetUserInfo(jwtTokenS)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		oryCookie, err := r.Cookie("ory_kratos_session")
		if err != nil {
			log.Println("no Ory Kratos cookie found in request")
		} else {
			logger.Debug("Ory Kratos cookie found in request", "cookie", oryCookie.Value)
			username, email, err = kratos.GetUserInfo(oryCookie.Value)
			if err != nil {
				logger.Debug("", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(config.Auth.JWT.SecretKey), nil)
		claims := map[string]interface{}{
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(time.Duration(config.Auth.JWT.ExpirationTime) * time.Second).Unix(),
			"iat":      time.Now().Unix(),
			"iss":      config.Auth.JWT.Issuer,
			"aud":      config.Auth.JWT.Audience,
			"sub":      config.Auth.JWT.Subject,
			// guid for jti
			"jti":        uuid.New().String(),
			"authSource": "xxxxxx",
		}
		// jwtauth.SetExpiryIn(claims, time.Duration(config.Auth.JWT.ExpirationTime))
		jwtToken, tokenString, err := tokenAuth.Encode(claims)
		// check expiry
		if jwtToken.Expiration().Before(time.Now()) {
			log.Println("Token expired")
			http.Error(w, "Token expired", http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: tokenString,
			Path:  "/",
		})
		http.Redirect(w, r, "/secure", http.StatusFound)
	})

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		// execute the register template
		err := tmpl.ExecuteTemplate(w, "registration_form", map[string]interface{}{
			"RegisterURL": config.Auth.SignInUrl,
			"sign_in_url": config.Auth.SignInUrl,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Group(func(r chi.Router) {
		jwtAuth := jwtauth.New("HS256", []byte(config.Auth.JWT.SecretKey), nil)
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(AddClientJWTStringToCtx)
		r.Use(AuthenticateJWT)
		r.Use(middlewareMap["addUserToCtx"])
		r.Get("/edit_org_form", func(w http.ResponseWriter, r *http.Request) {
			orgs, err := queries.ListOrganizations(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var orgMap []map[string]interface{}
			for _, org := range orgs {
				humanReadableAttributes, err := json.MarshalIndent(org.Attributes, "", "  ")
				if err != nil {
					log.Println("error prettifying json:" + err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				orgMap = append(orgMap, map[string]interface{}{
					"ID":         org.ID,
					"Name":       org.Name,
					"Email":      org.Email,
					"Attributes": string(humanReadableAttributes),
				})

			}
			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			p := components.NewPresentor(tmpl, w)
			err = p.EditListItemHTML("edit_org_form", "Edit Organization Form", "/api/organization", "put", "/edit_org_form", formFields, orgMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/compose_msg", func(w http.ResponseWriter, r *http.Request) {
			// Define a template string for the message form

			presentor := components.NewPresentor(tmpl, w)
			presentor.CreateItemFormHTML("compose_msg_form", "Compose Message", "/message", "/inbox", []components.FormField{
				{
					Name:  "ToUsername",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Subject",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Message",
					Type:  "text",
					Value: "",
				},
			})
		})

		r.Get("/create_org_form", func(w http.ResponseWriter, r *http.Request) {
			formFields := []components.FormField{
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "{}",
				},
			}
			p := components.NewPresentor(tmpl, w)
			err := p.CreateItemFormHTML("create_org_form", "Create Organization Form", "/api/organization", "/edit_org_form", formFields)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/create_user_form", func(w http.ResponseWriter, r *http.Request) {
			formFields := []components.FormField{
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "{}",
				},
			}

			p := components.NewPresentor(tmpl, w)
			err := p.CreateItemFormHTML("create_user_form", "Create User Form", "/api/user", "/edit_user_form", formFields)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/create_usergroup_form", func(w http.ResponseWriter, r *http.Request) {
			formFields := []components.FormField{
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "{}",
				},
			}
			p := components.NewPresentor(tmpl, w)
			err := p.CreateItemFormHTML("create_usergroup_form", "Create User Group Form", "/api/usergroup", "/edit_usergroup_form", formFields)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/delete_org_form", func(w http.ResponseWriter, r *http.Request) {
			orgs, err := queries.ListOrganizations(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var orgMap []map[string]interface{}
			for _, org := range orgs {
				orgMap = append(orgMap, map[string]interface{}{
					"ID":         org.ID,
					"Name":       org.Name,
					"Email":      org.Email,
					"Attributes": string(org.Attributes),
				})
			}

			p := components.NewPresentor(tmpl, w)
			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			additionalJS := `
					document.getElementById("deleteItemButton").setAttribute('hx-delete', '/api/organization/' + formData["ID"]);
					htmx.process(document.getElementById('deleteItemButton'));			
				`
			err = p.DeleteItemFormHTML("delete_org_form", "Delete Organization Form", "/api/organization", "/edit_org_form", additionalJS, formFields, orgMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/delete_user_form", func(w http.ResponseWriter, r *http.Request) {
			users, err := queries.ListUsers(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var userMap []map[string]interface{}
			for _, user := range users {
				userMap = append(userMap, map[string]interface{}{
					"ID":         user.ID,
					"Name":       user.Name,
					"Email":      user.Email,
					"Attributes": string(user.Attributes),
				})
			}

			p := components.NewPresentor(tmpl, w)
			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			additionalJS := `
				document.getElementById("deleteItemButton").setAttribute('hx-delete', '/api/user/' + formData["ID"]);
				htmx.process(document.getElementById('deleteItemButton'));			
			`
			err = p.DeleteItemFormHTML("delete_user_form", "Delete User Form", "/api/user", "/edit_user_form", additionalJS, formFields, userMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/delete_usergroup_form", func(w http.ResponseWriter, r *http.Request) {
			groups, err := queries.ListUserGroups(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var groupMap []map[string]interface{}
			for _, group := range groups {
				groupMap = append(groupMap, map[string]interface{}{
					"ID":         group.ID,
					"Name":       group.Name,
					"Email":      group.Email,
					"Attributes": string(group.Attributes),
				})
			}

			p := components.NewPresentor(tmpl, w)
			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Description",
					Type:  "text",
					Value: "",
				},
			}

			additionalJS := `
				document.getElementById("deleteItemButton").setAttribute('hx-delete', '/api/usergroup/' + formData["ID"]);
				htmx.process(document.getElementById('deleteItemButton'));			
			`

			err = p.DeleteItemFormHTML("delete_usergroup_form", "Delete Usergroup Form", "/api/usergroup", "/edit_usergroup_form", additionalJS, formFields, groupMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/edit_user_form", func(w http.ResponseWriter, r *http.Request) {
			users, err := queries.ListUsers(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var userMap []map[string]interface{}
			for _, user := range users {
				humanReadableAttributes, err := json.MarshalIndent(user.Attributes, "", "  ")
				if err != nil {
					log.Println("error prettifying json:" + err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				userMap = append(userMap, map[string]interface{}{
					"ID":         user.ID,
					"Name":       user.Name,
					"Email":      user.Email,
					"Attributes": string(humanReadableAttributes),
				})

			}
			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			p := components.NewPresentor(tmpl, w)
			err = p.EditListItemHTML("edit_user_form", "Edit User Form", "/api/user", "put", "/edit_user_form", formFields, userMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/edit_usergroup_form", func(w http.ResponseWriter, r *http.Request) {
			usergroups, err := queries.ListUserGroups(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var usergroupMap []map[string]interface{}
			for _, usergroup := range usergroups {
				humanReadableAttributes, err := json.MarshalIndent(usergroup.Attributes, "", "  ")
				if err != nil {
					log.Println("error prettifying json:" + err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				usergroupMap = append(usergroupMap, map[string]interface{}{
					"ID":         usergroup.ID,
					"Name":       usergroup.Name,
					"Attributes": string(humanReadableAttributes),
				})
			}

			formFields := []components.FormField{
				{
					Name:  "ID",
					Type:  "",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			p := components.NewPresentor(tmpl, w)
			err = p.EditListItemHTML("edit_usergroup_form", "Edit User Group Form", "/api/usergroup", "put", "/edit_usergroup_form", formFields, usergroupMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/upload_form", func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("/upload_form")
			err := tmpl.ExecuteTemplate(w, "upload_form", map[string]interface{}{
				"PostURL": "/api/upload",
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/debug", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.ExecuteTemplate(w, "debug", map[string]interface{}{})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
			var username string

			if r.Context().Value(User) == nil {
				username = ""
			} else {
				username = r.Context().Value(User).(db.User).Name
			}
			err := tmpl.ExecuteTemplate(w, "index", map[string]interface{}{
				"SecureText":  username,
				"sign_in_url": config.Auth.SignInUrl,
				"username":    r.Context().Value(User).(db.User).Name,
				"user_id":     r.Context().Value(User).(db.User).ID.String(),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/sidebar", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.ExecuteTemplate(w, "sidebar", map[string]interface{}{})
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/inbox", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			messages, err := queries.ListMessagesTo(r.Context(), r.Context().Value(User).(db.User).Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.New("messages").Parse(`
			<!DOCTYPE html>
			<html>
			  <head>
				<meta charset="utf-8">
				<title>Messages</title>
				<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
			  </head>
			  <body>
				<div class="container">
				  <h1>Messages</h1>
				  <table class="table">
					<thead>
					  <tr>
						<th>ID</th>
						<th>FromID</th>
						<th>Subject</th>
						<th>Message</th>
						<th>Created Time</th>
						<th>Updated Time</th>
					  </tr>
					</thead>
					<tbody>
					  {{range .}}
					  <tr>
						<td>{{.ID}}</td>
						<td>{{.FromID}}</td>
						<td>{{.Subject}}</td>
						<td>{{.Message}}</td>
						<td>{{.CreatedTS.Format "2006-01-02 15:04:05"  }}</td>
						<td>{{.UpdatedTS.Format "2006-01-02 15:04:05"}}</td>
					  </tr>
					  {{end}}
					</tbody>
				  </table>
				</div>
			  </body>
			</html>
			`)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, messages)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		r.Get("/sent_messages", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(User).(db.User).ID
			messages, err := queries.ListMessagesFrom(r.Context(), userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.New("messages").Parse(`
				<div class="container">
				  <h1>Sent Messages</h1>
				  <table class="table">
					<thead>
					  <tr>
						<th>ID</th>
						<th>To</th>
						<th>Subject</th>
						<th>Message</th>
						<th>Created Time</th>
						<th>Updated Time</th>
					  </tr>
					</thead>
					<tbody>
					  {{range .}}
					  <tr>
						<td>{{.ID}}</td>
						<td>{{.ToUsername}}</td>
						<td>{{.Subject}}</td>
						<td>{{.Message}}</td>
						<td>{{.CreatedTS.Format "2006-01-02 15:04:05"}}</td>
						<td>{{.UpdatedTS.Format "2006-01-02 15:04:05"}}</td>
					  </tr>
					  {{end}}
					</tbody>
				  </table>
				</div>
			  </body>
			</html>
			`)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, messages)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		r.Get("/new_blog_post", func(w http.ResponseWriter, r *http.Request) {
			presentor := components.NewPresentor(tmpl, w)
			presentor.CreateBlogHTML()
		})

		r.Post("/message", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			// get user id
			logger.Debug("/message", "form", r.Form)
			userID := r.Context().Value(User).(db.User).ID
			if userID == uuid.Nil {
				logger.Error(fmt.Sprintf("error getting user id"))
			}
			params := db.CreateMessageParams{
				ID:         uuid.New(),
				FromID:     userID,
				Subject:    r.FormValue("Subject"),
				ToUsername: r.FormValue("ToUsername"),
				Message:    r.FormValue("Message"),
			}

			newMessage, err := queries.CreateMessage(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error creating message: %v", err))
			}
			log.Println(newMessage)
		})

		r.Post("/blog_post", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			// get user id
			log.Printf("form: %v", r.Form)
			userID := r.Context().Value(User).(db.User).ID
			if userID == uuid.Nil {
				logger.Error(fmt.Sprintf("error getting user id"))
			}
			params := db.CreateBlogPostParams{
				ID:               uuid.New(),
				Title:            r.FormValue("Title"),
				Subtitle:         r.FormValue("Subtitle"),
				FeaturedImageURI: r.FormValue("FeaturedImageURI"),
				Body:             r.FormValue("Body"),
				AuthorID:         userID,
			}

			newBlogPost, err := queries.CreateBlogPost(r.Context(), params)
			if err != nil {
				logger.Error(fmt.Sprintf("error creating blog post: %v", err))
			}

			templateText := `
				<div class="container mt-5">
					<div class="alert alert-success" role="alert">
						<h4 class="alert-heading">New Blog Post Created!</h4>
						<p>Blog Post Name: <strong>{{ .Title }}</strong></p>
						<p>Subtitle: <strong>{{ .Subtitle }}</strong></p>
						<p>ID: <strong>{{ .ID }}</strong></p>
						<hr>
						<p class="mb-0">Blog posts can be viewed in the published posts center or the public page.</p>
					</div>
				</div>
			`

			tmpl, err := template.New("").Parse(templateText)
			if err != nil {
				logger.Error(fmt.Sprintf("error parsing template: %v", err))
			}

			err = tmpl.Execute(w, newBlogPost)
			if err != nil {
				logger.Error(fmt.Sprintf("error executing template: %v", err))
			}
		})

		r.Get("/my_blog_posts", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(User).(db.User).ID
			blogPosts, err := queries.ListBlogPostsByUser(r.Context(), userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			templateText := `
				<div class="container">
				  <h1>Blog Posts</h1>
				  <table class="table">
					<thead>
					  <tr>
					  	<th>ID</th>
						<th>Title</th>
						<th>Subtitle</th>
						<th>Body</th>
						<th>Created Time</th>
						<th>Updated Time</th>
					  </tr>
					</thead>
					<tbody>
					  {{range .}}
					  <tr>
					  	<td>{{.ID}}</td>
						<td>{{.Title}}</td>
						<td>{{.Subtitle}}</td>
						<td>{{.Body}}</td>
						<td>{{.CreatedTS.Format "2006-01-02 15:04:05"}}</td>
						<td>{{.UpdatedTS.Format "2006-01-02 15:04:05"}}</td>
					  </tr>
					  {{end}}
					</tbody>
				  </table>
				</div>`
			tmpl, err := template.New("").Parse(templateText)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, blogPosts)
		})
	})
}
