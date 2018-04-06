# Playground for Zeebe with go

## Setup

1) Install Go

2) Install Zeebe zbctl => Use latest release from https://github.com/zeebe-io/zbctl/releases

3) If you are not on MacOS replace 0.0.0.0 in docker-compose.yml and main.go with your docker ip.

## Run it

1) Start broker with: `docker-compose up`

2) To run the program: `go run src/main.go`
This will deploy and start an easy process.

3) Download latest simple monitor: https://github.com/zeebe-io/zeebe-simple-monitor/releases

4) Start Monitor `java -jar zeebe-simple-monitor-0.3.0.jar`

5) Check Monitor on: http://127.0.0.1:8080/

6) Add Broker with "[DOCKER_IP]:51015" (e.g. 0.0.0.0:51015 on MacOS)

More infos: https://docs.zeebe.io/go-client/get-started.html
