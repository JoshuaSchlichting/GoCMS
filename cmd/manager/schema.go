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
	filePayload := readFile(filename)

	queries := make([]string, 0)
	for _, query := range strings.Split(string(filePayload), ";") {
		queries = append(queries, strings.TrimSpace(query)+";")
	}
	return queries
}

func DestroySchema(db *sql.DB) {
	// drop all tables
	db.Exec("drop table public.messages;")
	db.Exec("drop table public.file;")
	db.Exec("drop table public.user;")
	db.Exec("drop table public.messages;")
	db.Exec("drop table public.invoice;")
}
