# GoCMS
##### GoCMS is intended to be the boiler plate code you need to build a CMS.
___

### Current state of development
What is working:
- Admin panel auth is working with AWS Cogntio.
- `filesystem` package provides a working local filesystem and s3 filesystem of the same interface.
- Dynamic HTML tables and forms are working.
- Basic CRUD functionality is working with GUI interfaces, currently impelmented for the `user`, `organization`, and `user_group` tables.

### Dependencies
- While the idea is to not have many dependencies, some things are still being imported. At the time of writing this, routing and JWT are handled by [Chi](https://github.com/go-chi/chi) (and they probably always will be).

- The `db` package is genereated using [sqlc](https://docs.sqlc.dev/en/latest/index.html)  (`db/sql/sqlc.yaml`).

### Using the manager
The manager is a CLI tool that is used to manage the CMS. It is embedded in the project as a [go workspace](https://go.dev/ref/mod#workspaces). You can run it using the command `go run ./cmd/manager <params>`.

### The Data Layer / SQL
The data layer is created by using [sqlc](https://docs.sqlc.dev/en/latest/index.html). The `db` package is generated using the `db/sql/sqlc.yaml` file. The `db/sql` folder contains the SQL files that are used to generate the `db` package. To generate new `sqlc` output, execute the following:
`cd db/sql && sqlc generate`

### Design
As of now, the CMS is relying on free templates from BootStrapMade.com.