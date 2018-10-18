#--------------------------------------
# Stage: Building heartbirth
#--------------------------------------
FROM golang:alpine

RUN mkdir -p /usr/local/go/src/github.com/enixdark/faktory_docker_heartbeat
ADD . /usr/local/go/src/github.com/enixdark/faktory_docker_heartbeat
WORKDIR /usr/local/go/src/github.com/enixdark/faktory_docker_heartbeat
RUN go build -o heartbirt .

#--------------------------------------
# Stage: Packaging App
#--------------------------------------
FROM contribsys/faktory:latest

VOLUME /app
COPY --from=0 /usr/local/go/src/github.com/enixdark/faktory_docker_heartbeat/heartbirt /

EXPOSE 7419 
EXPOSE 7420

CMD ["/faktory", "-w", "0.0.0.0:7420", "-b", "0.0.0.0:7419", "-e", "development"]
