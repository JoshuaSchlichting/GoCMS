package blog

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
)

func InitRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.Queries, middlewareMap map[string]func(http.Handler) http.Handler) {

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		// get posts from db
		dbPosts, err := queries.ListBlogPosts(r.Context())
		if err != nil {
			log.Printf("error getting posts: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// convert db posts to blog.Post's
		posts := []Post{}
		for _, post := range dbPosts {
			posts = append(posts, Post{
				ID:                 post.ID,
				Title:              post.Title,
				Subtitle:           post.Subtitle,
				FeaturedImageUri:   post.FeaturedImageUri,
				Body:               post.Body,
				PublishedTimestamp: post.CreatedTS,
				UpdatedTimestamp:   post.UpdatedTS,
			})
		}
		// get templates.HTML from "blog/posts" template
		var body bytes.Buffer
		// create writer
		bodyWriter := io.Writer(&body)
		err = tmpl.ExecuteTemplate(bodyWriter, "blog/posts", map[string]interface{}{
			"Posts":        posts[1:],
			"FeaturedPost": posts[0],
		})
		if err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// render template
		tmpl.ExecuteTemplate(w, "blog", Page{
			Title:      "Go CMS | A CMS built in Golang!",
			Brand:      "GoCMS",
			SignInURL:  config.Auth.SignInUrl,
			Heading:    "This is the heading",
			Subheading: "This is the subheading",
			NavBarLinks: []NavBarLink{
				{URL: "/", Text: "Home"},
				{URL: "/blog", Text: "Blog"},
				{URL: "/about", Text: "About"},
				{URL: "/contact", Text: "Contact"},
			},
			Body:         template.HTML(body.String()),
			FeaturedPost: posts[0],
			SideWidget: SideWidget{
				Title: "Blog Notice",
				Body:  "Look, there's a lot of 'frameworks' for creating a CMS. What if we just had a dead simple pattern to follow while using as much of Go's standard library as possible?",
			},
		})
	})

	r.Get("/blog/{id}", func(w http.ResponseWriter, r *http.Request) {
		// get post from db
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Printf("/blog/{id} error parsing id: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		post, err := queries.GetBlogPost(r.Context(), id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		p := Post{
			ID:                 post.ID,
			Title:              post.Title,
			Subtitle:           post.Subtitle,
			FeaturedImageUri:   post.FeaturedImageUri,
			Body:               post.Body,
			PublishedTimestamp: post.CreatedTS,
			UpdatedTimestamp:   post.UpdatedTS,
		}
		// get templates.HTML from "blog/post" template
		var body bytes.Buffer
		// create writer
		bodyWriter := io.Writer(&body)
		err = tmpl.ExecuteTemplate(bodyWriter, "blog/post", p)
		if err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// render template
		tmpl.ExecuteTemplate(w, "blog", Page{
			Title:      "Go CMS | A CMS built in Golang!",
			Brand:      "GoCMS",
			SignInURL:  config.Auth.SignInUrl,
			Heading:    "This is the heading",
			Subheading: "This is the subheading",
			NavBarLinks: []NavBarLink{
				{URL: "/", Text: "Home"},
				{URL: "/blog", Text: "Blog"},
				{URL: "/about", Text: "About"},
				{URL: "/contact", Text: "Contact"},
			},
			Body: template.HTML(body.String()),
			SideWidget: SideWidget{
				Title: "Blog Notice",
				Body:  "Look, there's a lot of 'frameworks' for creating a CMS. What if we just had a dead simple pattern to follow while using as much of Go's standard library as possible?",
			},
		})
	})
}
