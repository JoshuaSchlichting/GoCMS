package manager

import (
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/joshuaschlichting/gocms/config"
	database "github.com/joshuaschlichting/gocms/db"
	_ "github.com/lib/pq"
)

var sqlDir fs.FS

func IsManagerProgramCall(configuration config.Config, sqlDirA fs.FS) bool {
	sqlDir = sqlDirA
	db, err := sql.Open("postgres", configuration.Database.ConnectionString)
	queries := database.New(db)

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
		return false
	}
	return true
}

func readFile(filename string) []byte {
	file, err := sqlDir.Open(filename)
	filePayload := make([]byte, 0)
	file.Read(filePayload)
	if err != nil {
		log.Fatal(err)
	}
	return filePayload
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
