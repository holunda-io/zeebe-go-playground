# Playground for Zeebe with go

## Setup

1) Install Go

2) Install Zeebe zbctl

    go get github.com/zeebe-io/zbc-go
    cd $GOPATH/src/github.com/zeebe-io/zbc-go
    make build
    sudo make install

3) If you are not on MacOS replace 0.0.0.0 in docker-compose.yml and main.go with your docker ip.

## Run it

1) Start broker with: `docker-compose up`

2) Create topic: `zbctl create topic --name default-topic --partitions 1`

3) To run the program: `go run src/main.go`
This will deploy and start an easy process.

4) Download latest simple monitor: https://github.com/zeebe-io/zeebe-simple-monitor/releases

5) Start Monitor `java -jar zeebe-simple-monitor-0.3.0.jar`

6) Check Monitor on: http://127.0.0.1:8080/

7) Add Broker with "[DOCKER_IP]:51015" (e.g. 0.0.0.0:51015 on MacOS)

More infos: https://docs.zeebe.io/go-client/get-started.html
