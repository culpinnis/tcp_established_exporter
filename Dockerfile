FROM golang:alpine AS build

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/ipbhalle/tcp_established_exporter

# List project dependencies with Gopkg.toml and Gopkg.lock
# These layers are only re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml /go/src/github.com/ipbhalle/tcp_established_exporter/
WORKDIR /go/src/github.com/ipbhalle/tcp_established_exporter/
# Install library dependencies
RUN dep ensure --vendor-only

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
COPY prometheus_tcp_established_exporter.go /go/src/github.com/ipbhalle/tcp_established_exporter/
RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o /bin/prometheus_tcp_established_exporter
#the exports are necessary https://forums.docker.com/t/getting-panic-spanic-standard-init-linux-go-178-exec-user-process-caused-no-such-file-or-directory-red-while-running-the-docker-image/27318/14
#otherwise docker won't start the container: standard_init_linux.go:211: exec user process caused "no such file or directory"                                                                                                                                                       exit:1

# This results in a single layer image
FROM scratch
COPY --from=build /bin/prometheus_tcp_established_exporter /bin/prometheus_tcp_established_exporter

EXPOSE 9669
ENTRYPOINT ["/bin/prometheus_tcp_established_exporter"]
