package manager

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"

	"github.com/joshuaschlichting/gocms/internal/config"
	database "github.com/joshuaschlichting/gocms/internal/data/db"
	_ "github.com/lib/pq"
)

var sqlDir fs.FS

func HandleIfManagerCall(configuration config.Config, sqlDirA fs.FS) bool {
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
		log.Println("Initializing database schema...")
		createSchema(db)
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
	if err != nil {
		log.Fatal("error opening sql file:", err)
	}
	filePayload := new([]byte)
	*filePayload, err = io.ReadAll(file)
	if err != nil {
		log.Fatal("error reading sql file:", err)
	}
	return *filePayload
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
