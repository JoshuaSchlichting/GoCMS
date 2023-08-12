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
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/api"
	"github.com/joshuaschlichting/gocms/apps/admin"
	"github.com/joshuaschlichting/gocms/apps/public/blog"
	"github.com/joshuaschlichting/gocms/apps/public/landing_page"
	"github.com/joshuaschlichting/gocms/auth"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/data/cache"
	database "github.com/joshuaschlichting/gocms/data/db"
	"github.com/joshuaschlichting/gocms/filesystem"
	"github.com/joshuaschlichting/gocms/manager"
	"github.com/joshuaschlichting/gocms/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

//go:embed static
var fileSystem embed.FS

//go:embed apps
var templateFS embed.FS

//go:embed config.yml
var configYml []byte

//go:embed data/db/sql
var sqlFS embed.FS

var logger *slog.Logger

func init() {
	// Set up logging
	debugFlag := flag.Bool("debug", false, "debug logging mode")
	flag.Parse()
	var programLevel = new(slog.LevelVar)

	if *debugFlag {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	api.SetLogger(logger)
	admin.SetLogger(logger)
	blog.SetLogger(logger)
	auth.SetLogger(logger)
}

func main() {
	var (
		host = flag.String("host", "", "host http address to listen on")
		port = flag.String("port", "8000", "port number for http listener")
	)
	config := config.LoadConfig(readConfigFile())
	if manager.HandleIfManagerCall(*config, sqlFS) {
		logger.Info("Manager program call finished...")
		// This was a call to the manager program, not the web site executable
		return
	}

	flag.Parse()
	logger.Info(fmt.Sprintf("Starting server on port: %v", *port))
	db, err := sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create data layer
	queries := database.New(db)
	logger.Info("Connected to database: " + parseConnectionString(config.Database.ConnectionString))
	mu := new(sync.RWMutex)
	c := cache.New(mu)
	cache := database.NewDBCache(queries, c)

	// Add template functions
	funcMap := template.FuncMap{
		"mod": func(i, j int) int {
			return i % j
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"seq": func(start, end int) []int {
			var sequence []int
			for i := start; i <= end; i++ {
				sequence = append(sequence, i)
			}
			return sequence
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"eq": func(a, b int) bool {
			return a == b
		},
	}

	// Load templates
	logger.Info("Loading templates...")
	templ, err := parseTemplateDir("apps", templateFS, funcMap)
	if err != nil {
		errMsg := fmt.Sprintf("error parsing templates: %v", err)
		logger.Error(errMsg)
		log.Fatalf(errMsg)
	}

	// Middleware /////////////////////////////////////////////////////////////
	// Initialize middlware
	middleware.InitMiddleware(config)

	// Create the router
	r := chi.NewRouter()

	// Register common middleware with the router
	dbMiddlware := middleware.NewMiddlewareWithDB(*cache, config.Auth.JWT.SecretKey)
	r.Use(middleware.LogAllButStaticRequests)
	middlewareMap := map[string]func(http.Handler) http.Handler{}
	middlewareMap["addUserToCtx"] = dbMiddlware.AddUserToCtx

	// Register static file serve
	staticFS, _ := fs.Sub(fileSystem, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Create file system for content delivery
	homeDir, _ := os.UserHomeDir()
	gocmsPath := path.Join(homeDir, "gocms")
	logger.Info(fmt.Sprintf("Using the following gocmsPath for local filesystem: %s", gocmsPath))
	fs := filesystem.NewLocalFilesystem(gocmsPath)

	// Register apps routes
	admin.InitRoutes(r, templ, config, *cache, middlewareMap)
	api.InitAPI(r, templ, config, *cache, fs)
	blog.InitRoutes(r, templ, config, *cache, middlewareMap)
	landing_page.InitRoutes(r, templ, config, *cache, middlewareMap)

	// Start server
	addr := net.JoinHostPort(*host, *port)
	if err := listenServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func listenServe(listenAddr string, handler http.Handler) error {
	s := http.Server{
		Addr:    listenAddr,
		Handler: handler, // Our own instance of servemux
	}
	logger.Debug(fmt.Sprintf("Starting HTTP listener at %s", listenAddr))
	return s.ListenAndServe()
}

func readConfigFile() []byte {
	// read config.yml
	return configYml
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
