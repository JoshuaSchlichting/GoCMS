package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func init() {
	loadTemplates()
}

var templateFiles = []string{}

func loadTemplates() {
	filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
}

func parseTemplateDir(dir string, templateFS fs.FS) (*template.Template, error) {
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
	return template.ParseFS(templateFS, paths...)
}
