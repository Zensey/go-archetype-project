curl  --user "user:pass" -v -X POST http://0.0.0.0:8888/api/reviews \
-H 'Content-Type: application/json' \
-d '{
"name": "Elvis Presley",
"email": "theking@elvismansion.com",
"productid": "4",
"review": "I really love the product and will recommend! "
}'