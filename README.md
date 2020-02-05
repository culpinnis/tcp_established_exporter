# tcp\_established\_exporter
![Go](https://github.com/culpinnis/tcp_established_exporter/workflows/Go/badge.svg?event=push)

## The exporter

The exporter measures the number of TCP connections established to the ports of the host machine. To be measured, a connection has to be established for a defined duration. For this reason the exported metric is called netstat\_tcp\_longterm\_connections\_total.

## The configuration

The exporter uses the package github.com/namsral/flag and can be configured via:

* Flags
* Environment variables
* Configuration file

### Possible settings
*simple:* [bool] If set to true the exported metric will be a gauge counting all longterm established connections. If false (default) it will create a gauge with labels for each observed port and TCP version.

*tcpv6:* [bool] If set to true (default) IPv6 connections will be measured, too.

*port:* [int] Sets a specific port to observe. The default is -1 which measures all ports (except the listening port of the exporter).

*listen:* [int] Sets the listening port of the exporter. The default value is 9669.

*duration:* [int] The minimal duration in seconds after a connection is concerned as a longterm connection. The default value is 6.

## Usage

### Build 
The exporter is written in Go and can be build and used with

```bash
#Get & build the exporter
go get github.com/ipbhalle/tcp_established_exporter
#Run the exporter
go run github.com/ipbhalle/tcp_established_exporter
```
### Application of settings
You can export the settings as environmental variables:
```bash
export SIMPLE=true
export TCPV6=true
```

flag:
```bash
tcp_established_exporter --simple=true --tcpv6=true
```

or config file:

```bash
tcp_established_exporter --conf sample.conf
```
with sample.conf

```
simple true
tcpv6 true
```

The metrics will be exposed at /metrics
### Possible Use Cases
The exporter can be used if you want to measure/observe the number of connections to one or multiple services. It can be useful for applications that use WebSockets because each user establishes a longterm connection.
It was designed to count the number of users of Shiny web applications because the Shiny open-source server has not metrics endpoint out of the box.
But it could be also used for other use cases, e.g. if you want to guarantee that there is a connection to a service (and trigger an alert otherwise).  
### Kubernetes
The docker image of this exporter can be used as a sidecar in Kubernetes to measure the longterm connections of a pod.
Different containers running on the same pod are sharing the resources of the pod. Because of that, it possible to measure the TCP connections of applications running in other containers.
A complete example will be shown in the future.
