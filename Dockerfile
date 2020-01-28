FROM alpine:latest

RUN apk add --no-cache git make musl-dev go

# Configure Go
RUN go get github.com/cakturk/go-netstat/netstat
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/namsral/flag

COPY prometheus_tcp_established_exporter.go /usr/local/bin/
EXPOSE 2112

ENTRYPOINT ["go", "run /usr/local/bin/prometheus_tcp_established_exporter.go"]
