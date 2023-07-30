package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func init() {
	loadTemplates(templateFS)
}

var templateFiles = []string{}

func loadTemplates(templateFS fs.FS) {
	err := fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking templates dir: %v", err)
	}
	// Print all HTML template files
	for _, file := range templateFiles {
		if strings.Contains(file, "templates") {
			log.Println(file)
		}
	}
}

func parseTemplateDir(dir string, templateFS fs.FS, funcMap template.FuncMap) (*template.Template, error) {
	var paths []string
	err := fs.WalkDir(templateFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("path=%q, isDir=%v\n", path, d.IsDir())
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Found templates %v", paths)

	return template.New("").Funcs(funcMap).ParseFS(templateFS, paths...)
}
