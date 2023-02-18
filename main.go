package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/api"
	"github.com/joshuaschlichting/gocms/config"
	database "github.com/joshuaschlichting/gocms/db"
	"github.com/joshuaschlichting/gocms/middleware"
	"github.com/joshuaschlichting/gocms/routes"
	_ "github.com/lib/pq"
)

//go:embed static
var fileSystem embed.FS

//go:embed templates
var templateFS embed.FS

func main() {
	var (
		host = flag.String("host", "", "host http address to listen on")
		port = flag.String("port", "8000", "port number for http listener")
	)
	flag.Parse()

	config := config.LoadConfig(readConfigFile())
	log.Print("connection string: ", config.Database.ConnectionString)

	db, err := sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)
	defer db.Close()
	funcMap := template.FuncMap{}
	templ, err := parseTemplateDir("templates", templateFS, funcMap)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	addr := net.JoinHostPort(*host, *port)

	r := chi.NewRouter()
	// Middleware /////////////////////////////////////////////////////////////
	// Initialize middlware
	middleware.InitMiddleware(config)

	// Register common middleware
	dbMiddlware := middleware.NewMiddlewareWithDB(queries, config.Auth.JWT.SecretKey)
	r.Use(middleware.LogAllButStaticRequests)

	middlewareMap := map[string]func(http.Handler) http.Handler{}
	middlewareMap["addUserToCtx"] = dbMiddlware.AddUserToCtx
	// End Middleware /////////////////////////////////////////////////////////

	// Register static file serve
	// new file system made from fileSystem sub folder "static"
	staticFS, _ := fs.Sub(fileSystem, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Register routes
	routes.InitGetRoutes(r, templ, config, queries, middlewareMap)
	api.InitPostRoutes(r, templ, config, queries)

	if err := listenServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func listenServe(listenAddr string, handler http.Handler) error {
	s := http.Server{
		Addr:    listenAddr,
		Handler: handler, // Our own instance of servemux
	}
	fmt.Printf("Starting HTTP listener at %s\n", listenAddr)
	return s.ListenAndServe()
}

func readConfigFile() []byte {
	// read config.yml
	configYml, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}
	return configYml
}
