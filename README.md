## Pre-requisites
* docker
* gnu make

### How to run

| Plugin | README |
| ------ | ------ |
| make docker-build | builds and runs services in docker |
| make all *        | builds program locally on uour host|


In second case you need a golang installed and sources placed in a proper place ($GOPATH/src/github.org/Zensey)

On start you will see a tmux session with this README and logs individually in separete panes.

## Api
To make a request you can run `sh review-probe.sh`

## Web page (report)
See a report by this link
http://admin:admin@localhost:8888/api/report

## Tests
To run unit-tests exit from nano and run `make test`

## Tips
* To select text region / link in tmux session press Shift + mouse button.
* To switch between screen windows use Ctrl-b + arrows.
* To scroll tty use Ctrl-b PgUp / PgDown
* Tmux cheatsheet https://gist.github.com/henrik/1967800

### Problems encountered with the AdventureWorks schema
* 67 tries of copy from inexistent files
* 1 violation of foreign key constraint "FK_ProductReview_Product_ProductID" (due to incomplete data)


## Logs
* /var/log/task/api-service.log
* /var/log/task/worker1.log
* /var/log/task/worker2.log

## Default run configuration
    /app/build/_workspace/bin/api -lb syslog -ll trace
    /app/build/_workspace/bin/worker1 -lb syslog -ll trace
    /app/build/_workspace/bin/worker2 -lb syslog -ll trace

## Arguments

| Plugin | README |
| ------ | ------ |
| -api string     | e.g.: :8888 (default ":8888") |
| -lb string      | log backend e.g.: console, syslog (default "console") |
| -ll string      | log level e.g.: error, warning, info, debug, trace (default "info") |
| -pg string      | default "postgres://docker:docker@localhost:5432/app" |
| -redis string   | default "foobared@localhost:6379") |
| -words string   | bad words. default "fee,nee,cruul,leent". Apllies only to `worker1`|
