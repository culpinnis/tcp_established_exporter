FROM alpine:latest

RUN apk add --no-cache git make musl-dev go

# Configure Go
RUN go get github.com/cakturk/go-netstat/netstat
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/namsral/flag
RUN go build prometheus_tcp_established_exporter.go

COPY prometheus_tcp_established_exporter /usr/local/bin/
EXPOSE 9669

ENTRYPOINT /usr/local/bin/prometheus_tcp_established_exporter
