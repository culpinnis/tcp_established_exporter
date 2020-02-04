package main

import (
	"github.com/namsral/flag"
	"fmt"
	"strings"
	"github.com/cakturk/go-netstat/netstat"
	"net"
	"time"
	"strconv"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)


const(
	proto = 0x01 | 0x02
)
var(
	tcpv6 bool = true
	port int = -1
	myport uint = 9669
	duration int = 6
)
var(
	netstat_tcp_connection_longterm_counts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "netstat_tcp_longterm_connections_total",
		Help: "TCP connections that are established for a minimal duration"})
)
func init(){
	prometheus.MustRegister(netstat_tcp_connection_longterm_counts)
}
func countSockInfo(connection_counts map[string]uint, s []netstat.SockTabEntry) map[string]uint {
	connection_counts_new := make(map[string]uint)
	lookup := func(skaddr *netstat.SockAddr) (string,uint) {
		addr := skaddr.IP.String()
			names, err := net.LookupAddr(addr)
			if err == nil && len(names) > 0 {
				addr = names[0]
			}
		return addr,uint(skaddr.Port)
	}

	for _, e := range s {
		saddr,sport := lookup(e.LocalAddr)
		daddr,dport := lookup(e.RemoteAddr)
		cur_c := string(daddr+"_"+strconv.Itoa(int(dport)) + "|" + saddr +"_"+strconv.Itoa(int(sport)))
		//fmt.Printf("%s %s %s \n", saddr, daddr, e.State)
		if((port==-1||port==int(sport))&&sport!=myport) {
			connection_counts_new[cur_c] = connection_counts[cur_c] + 1
		}
	}
	return(connection_counts_new)
}
func main() {
//get flags
	flag.BoolVar(&tcpv6, "tcpv6", true, "Should TCPV6 sockets be monitored?")
	flag.IntVar(&port, "port", -1, "The port that should be monitored. -1 monitors every port.")
	flag.IntVar(&myport, "listen", 9669, "The port on that this exporter listens for requests.")
	flag.IntVar(&duration, "duration", 6, "The minimal duration in seconds after a connection is concerned as longterm.")
	flag.Parse()
	connections := make(map[string]uint)
	connections6 := make(map[string]uint)

	go func() { //create a second thread that counts the connections
		for{
			socks, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
				return s.State == netstat.Established
			})
			if err != nil {
				fmt.Printf("Error")
			}
			connections = countSockInfo(connections, socks)

			if(tcpv6) {
				socks6, err := netstat.TCP6Socks(func(s6 *netstat.SockTabEntry) bool {
					return s6.State == netstat.Established
				})
				if err != nil {
					fmt.Printf("Error")
				}
				connections6 = countSockInfo(connections6, socks6)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {//create a thread that counts and exports the metric
		if(port==-1) {
			for{
				var sum uint = 0
				for _, value := range(connections) {
					if(value >= duration) {
						sum += 1
					}
				}
				if(tcpv6) {
					for _, value := range(connections6) {
						if(value >= duration) {
							sum += 1
						}
					}
				}
				time.Sleep(1 * time.Second)
				netstat_tcp_connection_longterm_counts.Set(float64(sum))
			}
		// }	else{
		// 	for{
		// 		connection_counts_port := make(map[string]uint)
		// 		for connection, value :=range(connections){
		// 			var dport string = strings.SplitN(strings.SplitN(connection, "|", 2)[1], "_", 2)[2]
		// 			netstat_tcp_connection_longterm_counts.WithLabelValues(dport).Observe(value)
		// 		}
		// 		//Add tcpv6 here!
		// 	}
		// }
	}()
	http.Handle("/metrics", promhttp.Handler())
  http.ListenAndServe(":"+strconv.Itoa(myport) , nil)
	//Add welcome page here
}
