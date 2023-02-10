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

	var listUsersFlag bool
	flag.BoolVar(&listUsersFlag, "list-users", false, "List users")

	var createSuperUserFlag bool
	flag.BoolVar(&createSuperUserFlag, "create-superuser", false, "Create super user")

	var deleteAllUsersFlag bool
	flag.BoolVar(&deleteAllUsersFlag, "delete-all-users", false, "Delete all users")

	var initFlag bool
	flag.BoolVar(&initFlag, "init", false, "Initialize database schema")

	var executeRawSqlFlag string
	flag.StringVar(&executeRawSqlFlag, "exec-sql", "", "Execute raw sql")
	flag.Parse()
	switch {
	case createSuperUserFlag:
		executeCreateSuperUserViaTerminalInput(queries)
	case initFlag:
		CreateSchema(db)
	case listUsersFlag:
		getUsers(queries)
	case deleteAllUsersFlag:
		deleteAllUsers(db)
	case executeRawSqlFlag != "":
		_, err := db.Exec(executeRawSqlFlag)
		if err != nil {
			log.Fatal(err)
		}
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

func executeCreateSuperUserViaTerminalInput(queries database.QueriesInterface) {
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
