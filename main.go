package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/middleware"
	"github.com/joshuaschlichting/gocms/routes"
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

	config := config.LoadConfig()

	templ, err := parseTemplateDir("templates", templateFS)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	addr := net.JoinHostPort(*host, *port)

	r := chi.NewRouter()

	// Middleware /////////////////////////////////////////////////////////////
	// Initialize middlware
	middleware.InitMiddleware(config)

	// Register common middleware
	r.Use(middleware.LogAllButStaticRequests)
	// End Middleware /////////////////////////////////////////////////////////

	// Register static file serve
	// new file system made from fileSystem sub folder "static"
	staticFS, _ := fs.Sub(fileSystem, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Register routes
	routes.InitGetRoutes(r, templ, config)

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
