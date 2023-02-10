package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func CreateSchema(db *sql.DB) {
	queries := LoadQueriesFromFile(filepath.Join(getProjectDir(), "db", "sql", "schema.sql"))
	for _, query := range queries {
		if query != "" {
			_, err := db.Exec(query)

			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println("Successfully executed query: ", query)
		}
	}
}

func LoadQueriesFromFile(filename string) []string {
	// load queries from file
	// read text from file named filename
	filePayload := readFile(filename)

	queries := make([]string, 0)
	// split file by ';'
	for _, query := range strings.Split(string(filePayload), ";") {
		queries = append(queries, strings.TrimSpace(query)+";")
	}
	return queries
}
