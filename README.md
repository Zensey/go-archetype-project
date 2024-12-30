# Word of wisdom service

## Run
```
docker-compose up
```

## Design notes

### Why cashhash is used:
- allows to make server "stateless" no need to remember challenge id, b/c extension field is 
  authenticated (HMAC), that also making server resilient to replay attack
- ease of implementation
- pow complexity is adjustable, however not as smoothly as in other algorithms

### Design considerations & assumptions.
- to be able to make e2e test, server must be shutdown-able (by means of context)
- we assume client can do any number of requests during a session
- each client request should be preceded by PoW challenge & response
- network I/O is the most error prone part, that's why stress e2e stress test is a must have
- both server and client implementations of protocol should be tested by mocking the 
  counterpart; for that purpose not.Conn should be wrapped, the wrapper should be mocked
- the tests should cover a variety of possible cases to help safeguard the protocol implementation against potential breakages caused by future changes
