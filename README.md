# faktory_docker_heartbeat
helm chart and docker build for faktory  

how to setup:

- before run , you need to run faktory service
`docker run --rm -it -p 7419:7419 -p 7420:7420 contribsys/faktory:latest`

- Then, prepare `dep` tool to install dependency packages
- `dep ensure`
- `go build -o heartbeat`
- `./heartbeat`


