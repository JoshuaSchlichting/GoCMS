package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

// parseTemplateDir parses all the templates in the given directory. Non *.html files are ignored.
func parseTemplateDir(dir string, templateFS fs.FS, funcMap template.FuncMap) (*template.Template, error) {
	var paths []string
	err := fs.WalkDir(templateFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking templates dir: %v", err)
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking templates dir: %v", err)
	}

	return template.New("").Funcs(funcMap).ParseFS(templateFS, paths...)
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
