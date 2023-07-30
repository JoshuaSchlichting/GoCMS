package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/api"
	"github.com/joshuaschlichting/gocms/config"
	database "github.com/joshuaschlichting/gocms/db"
	"github.com/joshuaschlichting/gocms/filesystem"
	"github.com/joshuaschlichting/gocms/manager"
	"github.com/joshuaschlichting/gocms/middleware"
	_ "github.com/lib/pq"
)

//go:embed static
var fileSystem embed.FS

//go:embed templates
var templateFS embed.FS

//go:embed config.yml
var configYml []byte

//go:embed db/sql
var sqlFS embed.FS

func main() {
	var (
		host = flag.String("host", "", "host http address to listen on")
		port = flag.String("port", "8000", "port number for http listener")
	)
	config := config.LoadConfig(readConfigFile())
	if manager.IsManagerProgramCall(*config, sqlFS) {
		log.Println("Manager program call finished...")
		// This was a call to the manager program, not the web site executable
		return
	}

	flag.Parse()
	log.Println("Starting server on port", *port)
	db, err := sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	queries := database.New(db)
	funcMap := template.FuncMap{}
	log.Println("Loading templates...")
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
	dbMiddlware := middleware.NewMiddlewareWithDB(*queries, config.Auth.JWT.SecretKey)
	r.Use(middleware.LogAllButStaticRequests)

	middlewareMap := map[string]func(http.Handler) http.Handler{}
	middlewareMap["addUserToCtx"] = dbMiddlware.AddUserToCtx
	// End Middleware /////////////////////////////////////////////////////////

	// Register static file serve
	// new file system made from fileSystem sub folder "static"
	staticFS, _ := fs.Sub(fileSystem, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Create file system for content delivery
	homeDir, _ := os.UserHomeDir()
	gocmsPath := path.Join(homeDir, "gocms")
	log.Println("Using the following gocmsPath for local filesystem: ", gocmsPath)
	fs := filesystem.NewLocalFilesystem(gocmsPath)

	// Register routes
	initRoutes(r, templ, config, *queries, middlewareMap)
	api.InitAPI(r, templ, config, *queries, fs)

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
	return configYml
}
