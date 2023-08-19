package main

import (
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
	auth "github.com/joshuaschlichting/gocms/auth/oauth2"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/filesystem"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/admin"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/api"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/blog"
	"github.com/joshuaschlichting/gocms/manager"
	"github.com/joshuaschlichting/gocms/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

//go:embed static
var fileSystem embed.FS

//go:embed internal/apps
var templateFS embed.FS

//go:embed config.yml
var configYml []byte

//go:embed internal/apps/cms/data/db/sql
var sqlFS embed.FS

var logger *slog.Logger

func init() {

}

func main() {
	var (
		host      = flag.String("host", "", "host http address to listen on")
		port      = flag.String("port", "8000", "port number for http listener")
		debugFlag = flag.Bool("debug", false, "debug logging mode")
	)
	// Set up logging
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
	config := config.LoadConfig(readConfigFile())
	if manager.HandleIfManagerCall(*config, sqlFS) {
		logger.Info("Manager program call finished...")
		// This was a call to the manager program, not the web site executable
		return
	}

	// Add template functions
	funcMap := commonFuncMap
	// Load templates
	logger.Info("Loading templates...")
	templ, err := parseTemplateDir("internal/apps", templateFS, funcMap)
	if err != nil {
		errMsg := fmt.Sprintf("error parsing templates: %v", err)
		logger.Error(errMsg)
		log.Fatalf(errMsg)
	}

	logger.Info(fmt.Sprintf("Starting server on port: %v", *port))

	// Create the router
	r := chi.NewRouter()

	// Create file system for content delivery
	homeDir, _ := os.UserHomeDir()
	gocmsPath := path.Join(homeDir, "gocms")
	logger.Info(fmt.Sprintf("Using the following gocmsPath for local filesystem: %s", gocmsPath))
	localFS := filesystem.NewLocalFilesystem(gocmsPath)

	// Register common middleware with the router

	r.Use(middleware.LogAllButStaticRequests)
	registerApps(r, templ, *config, *localFS)

	// Register static file serve
	staticFS, _ := fs.Sub(fileSystem, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

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

var commonFuncMap = template.FuncMap{
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
