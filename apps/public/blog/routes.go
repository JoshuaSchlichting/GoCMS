package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/data/db"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
}
func InitRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.DBCache, middlewareMap map[string]func(http.Handler) http.Handler) {

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		// get posts from db
		dbPosts, err := queries.ListBlogPosts(r.Context())
		if err != nil {
			logger.Error(fmt.Sprintf("error getting posts: %v", err))
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
				FeaturedImageURI:   post.FeaturedImageURI,
				Body:               post.Body,
				PublishedTimestamp: post.CreatedTS,
				UpdatedTimestamp:   post.UpdatedTS,
			})
		}
		// ... (The beginning part of your code remains the same)

		postsPerPage := 6
		totalPages := len(posts) / postsPerPage
		if len(posts)%postsPerPage != 0 {
			totalPages++
		}
		pageNumber := 1
		// Check if 'page' query parameter exists
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			if parsedPage, err := strconv.Atoi(pageParam); err == nil {
				pageNumber = parsedPage
			}
		}

		// Calculate start and end indices for slicing the posts based on the current page number.
		startIdx := (pageNumber - 1) * postsPerPage
		endIdx := startIdx + postsPerPage
		if endIdx > len(posts) {
			endIdx = len(posts)
		}

		var featuredPost *Post
		if pageNumber == 1 {
			featuredPost = &posts[0]
		} else {
			featuredPost = nil
		}
		currentPagePosts := posts[startIdx:endIdx] // This will contain the posts for the current page.

		// get templates.HTML from "blog/posts" template

		var body bytes.Buffer
		// create writer
		bodyWriter := io.Writer(&body)

		err = tmpl.ExecuteTemplate(bodyWriter, "blog/posts", PostsBody{
			Posts:        currentPagePosts,
			FeaturedPost: featuredPost,
			CurrentPage:  pageNumber,
			TotalPages:   totalPages,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("error executing template: %v", err))
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
				{URL: "/admin", Text: "Login"},
			},
			Body: template.HTML(body.String()),
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
			FeaturedImageURI:   post.FeaturedImageURI,
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
			logger.Error(fmt.Sprintf("error executing template: %v", err))
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
