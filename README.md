# GoCMS
##### GoCMS is intended to be the boiler plate code you need to build a CMS.
___

### Current state of development
What is working:
- Auth via [Ory Kratos](https://github.com/ory/kratos)
- `filesystem` package provides a working local filesystem and s3 filesystem of the same interface.
- Dynamic HTML tables and forms are working.
- Basic CRUD functionality is working with GUI interfaces, currently impelmented for the `user`, `organization`, and `user_group` tables.

### Dependencies
- While the idea is to not have many dependencies, some things are still being imported. At the time of writing this, routing and JWT are handled by [Chi](https://github.com/go-chi/chi) (and they probably always will be).

- The `db` package is genereated using [sqlc](https://docs.sqlc.dev/en/latest/index.html)  (`db/sql/sqlc.yaml`).

### Using the manager
The manager is a set of CLI args that can be used to manage the application. The manager is located in the `manager` package. To use the manager, execute the `manager.sh` shell script, followed by flags outlined in [`manager/flags.go`](manager/flags.go).

### The Data Layer / SQL
The data layer is created by using [sqlc](https://docs.sqlc.dev/en/latest/index.html). The `db` package is generated using the `db/sql/sqlc.yaml` file. The `db/sql` folder contains the SQL files that are used to generate the `db` package. To generate new `sqlc` output, execute the following:
`cd db/sql && sqlc generate`

### Access Login Service
Currently under development: `http://127.0.0.1:4455/welcome`

## Getting Started
### Docker Compose
Run this with a dummy postgresql database using 

```sh
docker compose \
    -f docker-compose.yml \
    -f auth/kratos/quickstart.yml \
    -f auth/kratos/quickstart-standalone.yml \
    -f auth/kratos/quickstart-postgres.yml \
    up
```
> NOTE: include `--build` to rebuild as needed.

If this is the first time starting the database, initialize the database using `./manager.sh --init`.

First time accounts created in dev mode can easily be verified by email using the attached Mail Slurper server at `http://127.0.0.1:4436/`