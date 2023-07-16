package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshuaschlichting/gocms/config"
	database "github.com/joshuaschlichting/gocms/db"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("GoCMS Manager")
	// read config.yml
	configYmlData := readConfigFile()
	configuration := config.LoadConfig(configYmlData)

	db, err := sql.Open("postgres", configuration.Database.ConnectionString)
	queries := database.New(db)

	print("Connected to db!\n")
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	switch {
	case createSuperUserFlag:
		executeCreateSuperUserViaTerminalInput(*queries)
	case initFlag:
		CreateSchema(db)
	case listUsersFlag:
		getUsers(*queries)
	case deleteAllUsersFlag:
		deleteAllUsers(db)
	case executeRawSqlFlag != "":
		if strings.HasPrefix(strings.ToLower(executeRawSqlFlag), "select") {
			rows, err := db.Query(executeRawSqlFlag)
			if err != nil {
				log.Fatal(err)
			}
			scanResult := interface{}(nil)
			for rows.Next() {
				fmt.Printf("%v", rows.Scan(&scanResult))
			}
		} else {
			result, err := db.Exec(executeRawSqlFlag)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}
	case destroySchemaFlag:
		log.Println("Destroying schema")
		DestroySchema(db)
	default:
		fmt.Println("No flags set")
	}

}

func readFile(filename string) []byte {
	filePayload, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return filePayload
}

func readConfigFile() []byte {
	configYml, err := os.ReadFile(filepath.Join(getProjectDir(), "config.yml"))
	if err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}
	return configYml
}

func executeCreateSuperUserViaTerminalInput(queries database.Queries) {
	var username string
	var email string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	createSuperUser(queries, username, email)
}

func getProjectDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	projectDir := strings.Split(wd, "gocms")[0] + "gocms"
	return projectDir
}
