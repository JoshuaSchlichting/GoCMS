package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	database "github.com/joshuaschlichting/gocms/data/db"
)

func deleteAllUsers(db *sql.DB) {
	result, error := db.Exec("drop table public.user;")
	if error != nil {
		fmt.Println(error)
		return
	}
	fmt.Println(result)
	fmt.Println("All users have been deleted! :(")
}

func createSuperUser(queries database.Queries, username, email string) {
	user := database.CreateUserParams{
		Name:       username,
		Email:      email,
		Attributes: json.RawMessage(`{"is_superuser": true}`),
	}
	userModel, err := queries.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Successfully created user: ", userModel.Name)
}

func getUsers(queries database.Queries) {
	users, err := queries.ListUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Users: ")
	for _, user := range users {
		fmt.Println(user.Name)
	}
}
