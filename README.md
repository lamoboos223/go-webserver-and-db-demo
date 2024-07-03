```sh
# install postgres driver
go get github.com/lib/pq
# create postgres docker container and seed it with data
docker run --name go-webserver-db -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 \
-v ./init.sql:/docker-entrypoint-initdb.d/init.sql \
-e POSTGRES_INITDB_ARGS="--username=postgres" \
postgres
# run go app
go run main.go

# make http request
curl --location 'http://localhost:8080/hello' \
--header 'apikey: apikey1' \
--header 'Content-Type: application/json' \
--data '{
    "text": "lama"
}'
curl --location 'http://localhost:8080/hello' \
--header 'apikey: apikey12' \
--header 'Content-Type: application/json' \
--data '{
    "text": "lama"
}'
curl --location 'http://localhost:8080/hello' \
--header 'Content-Type: application/json' \
--data '{
    "text": "lama"
}'
```