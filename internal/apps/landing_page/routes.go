package landing_page

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/data/db"
)

func InitRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.DBCache, middlewareMap map[string]func(http.Handler) http.Handler) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		model := LandingPageModel{
			Title:     "GoCMS",
			SignInURL: config.Auth.SignInUrl,
			Body:      "Welcome to GoCMS",
			NavBarLinks: []NavBarLink{
				{
					URL:  "/",
					Text: "Home",
				},
				{
					URL:  config.Auth.SignInUrl,
					Text: "Sign In",
				},
			},
			FeaturedItems: []FeaturedItem{
				{
					Title:    "For those about to rock...",
					Body:     "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quod aliquid, mollitia odio veniam sit iste esse assumenda amet aperiam exercitationem, ea animi blanditiis recusandae! Ratione voluptatum molestiae adipisci, beatae obcaecati.",
					ImageURL: "static/onepagewonder/assets/img/01.jpg",
				},
				{
					Title:    "We salute you!",
					Body:     "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quod aliquid, mollitia odio veniam sit iste esse assumenda amet aperiam exercitationem, ea animi blanditiis recusandae! Ratione voluptatum molestiae adipisci, beatae obcaecati.",
					ImageURL: "static/onepagewonder/assets/img/02.jpg",
				},
				{
					Title:    "Let there be rock!",
					Body:     "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quod aliquid, mollitia odio veniam sit iste esse assumenda amet aperiam exercitationem, ea animi blanditiis recusandae! Ratione voluptatum molestiae adipisci, beatae obcaecati.",
					ImageURL: "static/onepagewonder/assets/img/01.jpg",
				},
			},
			Brand:      "GoCMS",
			Heading:    "GoCMS",
			Subheading: "A CMS built in your favorite language!",
		}

		err := tmpl.ExecuteTemplate(w, "onepagewonder/index", model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
