package blog

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
)

func InitRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.Queries, middlewareMap map[string]func(http.Handler) http.Handler) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// get posts from db
		posts, err := queries.ListBlogPosts(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// render template
		tmpl.ExecuteTemplate(w, "blog", Page{
			Title:      "Go CMS | A CMS built in Golang!",
			Brand:      "GoCMS",
			SignInURL:  config.Auth.SignInUrl,
			Body:       "This is the blog page!",
			Heading:    "This is the heading",
			Subheading: "This is the subheading",
			NavBarLinks: []NavBarLink{
				{URL: "/", Text: "Home"},
				{URL: "/blog", Text: "Blog"},
				{URL: "/about", Text: "About"},
			},
			Posts: posts,
			SideWidget: SideWidget{
				Title: "Blog Notice",
				Body:  "Look, there's a lot of 'frameworks' for creating a CMS. What if we just had a dead simple pattern to follow while using as much of Go's standard library as possible?",
			},
		})
	})

	// r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	// get post from db
	// 	post, err := queries.GetPost(r.Context(), chi.URLParam(r, "id"))
	// 	if err != nil {
	// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// get page from db
	// 	page, err := queries.GetPage(r.Context(), "blog")
	// 	if err != nil {
	// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// get side widget from db
	// 	sideWidget, err := queries.GetSideWidget(r.Context())
	// 	if err != nil {
	// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// render template
	// 	tmpl.ExecuteTemplate(w, "post", Page{
	// 		Title:       page.Title,
	// 		Brand:       config.Brand,
	// 		SignInURL:   config.SignInURL,
	// 		Body:        page.Body,
	// 		Heading:     page.Heading,
	// 		Subheading:  page.Subheading,
	// 		NavBarLinks: page.NavBarLinks,
	// 		Posts: []db.BlogPost{
	// 			{
	// 				ID: post.ID,

	// 				Title: post.Title,
	// 				Body:  post.Body,
	// 			},
	// 		},
	// 		SideWidget: SideWidget{
	// 			Title: sideWidget.Title,
	// 			Body:  sideWidget.Body,
	// 		},
	// 	})
	// })
}
