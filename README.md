# tx parser

## Run
```
make txparser; ./txparser
```

API (example)
* curl http://localhost:8181/current-block
* curl -d "address=0x0a05bc5728218e565cf16dae82b2fd4d439dacf7" -X POST http://localhost:8181/subscribe
* curl http://localhost:8181/transactions

Design considerations & assumptions.
* by default peer addresses set is empty, thus no need to scan the block chain
* to make observer to scan we have to supply at least one address
* previously scanned blocks won't be rescanned when we add new addresses