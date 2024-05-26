curl -X PUT -d 'Hello, key-value store!' http://localhost:8080/v1/key/key-a -H "Content-Type: application/json"
curl -X PUT -d 'Hello, key-value store 2!' http://localhost:8080/v1/key/key-b -H "Content-Type: application/json"
curl http://localhost:8080/v1/key/key-a
curl http://localhost:8080/v1/key/key-b
curl -X DELETE http://localhost:8080/v1/key/key-b
curl http://localhost:8080/v1/key/key-a