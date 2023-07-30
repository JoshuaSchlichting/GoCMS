package main

import (
	"flag"
)

var createSuperUserFlag bool
var listUsersFlag bool
var deleteAllUsersFlag bool
var initFlag bool
var executeRawSqlFlag string
var destroySchemaFlag bool

func init() {

	flag.BoolVar(&listUsersFlag, "list-users", false, "List users")

	flag.BoolVar(&createSuperUserFlag, "create-superuser", false, "Create super user")

	flag.BoolVar(&deleteAllUsersFlag, "delete-all-users", false, "Delete all users")

	flag.BoolVar(&initFlag, "init", false, "Initialize database schema")

	flag.StringVar(&executeRawSqlFlag, "exec-sql", "", "Execute raw sql")

	flag.BoolVar(&destroySchemaFlag, "destroy-schema", false, "Destroy database schema")
}
