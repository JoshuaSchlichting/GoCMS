# GoCMS
##### GoCMS is intended to be the boiler plate code you need to build a CMS.
___

### Current state of development
GoCMS is in early stages of development and is not even close to being ready for production use. The code in this repository is not stable and will change rapidly over time.

### Dependencies
- While the idea is to not have many dependencies, some things are still being imported. At the time of writing this, routing and JWT are handled by [Chi](https://github.com/go-chi/chi) (and they probably always will be).

- The `db` package is genereated using [sqlc](https://docs.sqlc.dev/en/latest/index.html)  (`db/sql/sqlc.yaml`).

### Using the manager
The manager is a CLI tool that is used to manage the CMS. It is embedded in the project as a go workspace. You can run it using the command `go run ./cmd/manager <params>`.

### The Data Layer / SQL
The data layer is created by using [sqlc](https://docs.sqlc.dev/en/latest/index.html). The `db` package is generated using the `db/sql/sqlc.yaml` file. The `db/sql` folder contains the SQL files that are used to generate the `db` package. To generate new `sqlc` output, execute the following:
`cd db/sql && sqlc generate`

### Design
As of now, the CMS is relying on free templates from BootStrapMade.com.