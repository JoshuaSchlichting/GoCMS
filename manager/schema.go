package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CreateSchema(db *sql.DB) {
	queries := LoadQueriesFromFile(filepath.Join(getProjectDir(), "db", "sql", "schema.sql"))
	for _, query := range queries {
		if query != "" {
			_, err := db.Exec(query)

			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					fmt.Println("object already exists, skipping... (original error: ", err, ")")
					continue
				}
				log.Fatalf("an error occurred while executing the query against db '%v': %v", db, err)
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
	dropTablesFromSQLFile(filepath.Join(getProjectDir(), "db", "sql", "schema.sql"), db)
}

// dropTablesFromSQLFile reads a SQL file and drops all tables found in "create table" statements.
func dropTablesFromSQLFile(sqlFilePath string, db *sql.DB) error {
	defer db.Close()
	// Read the SQL file
	sqlFile, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Parse the SQL file for "create table" statements
	re := regexp.MustCompile(`create table if not exists ([^\s\(]+)`)
	scanner := bufio.NewScanner(strings.NewReader(string(sqlFile)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			tableName := matches[1]
			// Execute the DROP TABLE statement
			_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName))
			if err != nil {
				return fmt.Errorf("failed to drop table %s: %v", tableName, err)
			}
			log.Printf("dropped table %s", tableName)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to parse SQL file: %v", err)
	}

	return nil
}
