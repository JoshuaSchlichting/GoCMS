package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/cache"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/filesystem"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/admin"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/api"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/blog"
	database "github.com/joshuaschlichting/gocms/internal/apps/cms/data/db"
	"github.com/joshuaschlichting/gocms/internal/apps/landing_page"
)

func connecToDB(c config.Config) *sql.DB {
	db, err := sql.Open("postgres", c.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	// TODO: close db
	return db
}

func registerApps(r *chi.Mux, templ *template.Template, c map[string]config.Config, fs filesystem.LocalFilesystem) {
	/////////////////////////// CMS APP ///////////////////////////
	// Register apps databases
	blogDB := connecToDB(c["cms"])
	queries := database.New(blogDB)
	mu := new(sync.RWMutex)
	memC := cache.New(mu)
	cmsCache := database.NewDBCache(queries, memC)
	cmsConfig := c["cms"]

	// Register app middleware
	admin.InitMiddleware(&cmsConfig)
	middlewareMap := map[string]func(http.Handler) http.Handler{}
	dbMiddlware := admin.NewMiddlewareWithDB(*cmsCache, c["cms"].Auth.JWT.SecretKey)
	middlewareMap["addUserToCtx"] = dbMiddlware.AddUserToCtx

	// Register apps routes
	admin.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
	api.InitAPI(r, templ, &cmsConfig, *cmsCache, &fs)
	blog.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
	landing_page.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
}
