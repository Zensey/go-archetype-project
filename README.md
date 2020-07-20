#### To run 

* `make docker-build`
   
   Swagger-UI: http://localhost
   
   Swagger file: http://localhost:8080/files/swagger.yaml
   
   api_key (token): `key::1234`

#### To inspect the database (psql)

* `make docker-db-shell`

#### To reset the database

* `make docker-db-reset`


#### Addendum / P.S.

The most efficient way would be to make 1-st local persist of customer independent of
remote save operation, thus having local ID (pk) independent of and remote ID. 
But that requires use of _`integrationCode`_ attibute which is however not defined in
the Erply API wrapper. 