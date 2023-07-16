# `cmd/manager/`

`cmd/manager/` is a [go workspace](https://go.dev/ref/mod#workspaces) that contains the source code for the manager binary, which is used for performing
various management tasks, such as first time startup and initialization of the content management system 
(and its dependencies), and for performing various maintenance tasks, such as database migrations.


You can run it from the repo's root directory using the command `go run ./cmd/manager <params>`.