package manager

import (
	"flag"
)

var createSuperUserFlag bool
var listUsersFlag bool
var deleteAllUsersFlag bool
var initFlag bool
var executeRawSqlFlag string
var destroySchemaFlag bool
var appName string

func init() {
	managerFlag := flag.NewFlagSet("manager", flag.ExitOnError)
	managerFlag.StringVar(&appName, "app", "", "The name of the app to manage")

	managerFlag.BoolVar(&listUsersFlag, "list-users", false, "List users")

	managerFlag.BoolVar(&createSuperUserFlag, "create-superuser", false, "Create super user")

	managerFlag.BoolVar(&deleteAllUsersFlag, "delete-all-users", false, "Delete all users")

	managerFlag.BoolVar(&initFlag, "init", false, "Initialize database schema")

	managerFlag.StringVar(&executeRawSqlFlag, "exec-sql", "", "Execute raw sql")

	managerFlag.BoolVar(&destroySchemaFlag, "destroy-schema", false, "Destroy database schema")
}
