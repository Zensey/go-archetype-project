# Word of wisdom service

## Run
```
docker-compose up
```

## Design notes

### Why cashhash is used:
- allows to make server "stateless" no need to remember challenge id, b/c extension field is 
  authenticated (HMAC) thus making server resilient to replay attack
- ease of implementation
- pow complexity is adjustable, however not smoothly

### Design considerations & assumptions.
- to be able to make e2e test, server has to be shutdown-able (e.g. by means of context)
- we assume client can do any number of requests in one connection
- each client request follows is preceded by PoW challenge-response
- network I/O is the most error prone part, that's why stress e2e test is a must have
