package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func connecToDB(c config.Config, appName string) *sql.DB {
	db, err := sql.Open("postgres", c.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Registering database", "database", parseConnectionString(c.Database.ConnectionString), "app", appName)
	// defer db.Close()
	// TODO: close db
	return db
}

func registerApps(r *chi.Mux, templ *template.Template, c map[string]config.Config, fs filesystem.LocalFilesystem) {
	/////////////////////////// CMS APP ///////////////////////////
	cmsAppName := "cms"
	// Register apps databases
	blogDB := connecToDB(c[cmsAppName], cmsAppName)
	queries := database.New(blogDB)
	mu := new(sync.RWMutex)
	memC := cache.New(mu)
	cmsCache := database.NewDBCache(queries, memC)
	cmsConfig := c[cmsAppName]

	// Register app middleware
	admin.InitMiddleware(&cmsConfig)
	middlewareMap := map[string]func(http.Handler) http.Handler{}
	dbMiddlware := admin.NewMiddlewareWithDB(*cmsCache, c[cmsAppName].Auth.JWT.SecretKey)
	middlewareMap["addUserToCtx"] = dbMiddlware.AddUserToCtx

	// Register apps routes
	admin.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
	api.InitAPI(r, templ, &cmsConfig, *cmsCache, &fs)
	blog.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
	landing_page.InitRoutes(r, templ, &cmsConfig, *cmsCache, middlewareMap)
}

func parseConnectionString(connStr string) string {
	u, err := url.Parse(connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("error when parsing connectiong string for a URL: %v", err))
	}

	// Get DB name from path
	dbName := strings.TrimPrefix(u.Path, "/")

	return dbName + "@" + u.Host
}
