# Task


Using only the standard library, create a Go HTTP server that on each request responds
with a counter of the total number of requests that it has received during the previous 60 seconds (moving window).

The server should continue to the return the correct numbers after restarting it, by persisting data to a file.


Makefile rules
* make get-deps
* make demo
* make lint
* make docker-build
