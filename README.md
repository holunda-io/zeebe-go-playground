# Playground for Zeebe with go

1) Start broker and monitor with: `docker-compose up`

2) Create topic: `zbctl create topic --name default-topic --partitions 1`

3) To run the program: `go run src/main.go`
This will deploy and start an easy process.

4) Check Monitor on: http://127.0.0.1:9000/

More infos: https://docs.zeebe.io/go-client/get-started.html
