# go-archetype-project

Golang archetype project with following features:
 * Makefile
 * statical code analyzers & checkers,
 * local GOPATH and workplace
 ** dependecies got & stored locally and separately from sources
 * use of go dep to automatically find dependencies

 * stringer generator
 * logger helper with levels of logging, string formatting
 * `Dockerfile` and `docker-compose.yml` which allow to boot up application in a single `docker-compose up` command.

Makefile rules
* make get-deps
* make demo
* make lint
* make docker-build
