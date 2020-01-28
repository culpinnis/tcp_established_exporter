package main

import (
	"fmt"
	"github.com/cakturk/go-netstat/netstat"
	"net"
	"flag"
	"time"
	"strconv"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/namsral/flag"
)

var(
	//udp       = flag.Bool("udp", false, "display UDP sockets")
	//tcpv6       = flag.Bool("tcpv6", true, "display TCPv6 sockets")
	//port 	  = flag.Int("port", -1, "port that should be monitored")
	var tcpv6 bool;
	flag.bool(&tcpv6, "tcpv6", "Should TCPV6 sockets be monitored?")
	var port int;
	port := -1;
	flag.int(&port, "port", "The port that should be monitored. -1 monitors every port.")
)

const(
	proto = 0x01 | 0x02
	myport = 2112
)

var(
	netstat_tcp_connection_longterm_counts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "netstat_tcp_connection_longterm_counts",
		Help: "TCP connections that are established for a minimal duration"})
)
func init(){
	prometheus.MustRegister(netstat_tcp_connection_longterm_counts)
}
func displaySockInfo(s []netstat.SockTabEntry) {
	lookup := func(skaddr *netstat.SockAddr) string {
		addr := skaddr.IP.String()
			names, err := net.LookupAddr(addr)
			if err == nil && len(names) > 0 {
				addr = names[0]
			}
		return fmt.Sprintf("%s %d", addr, skaddr.Port)
	}

	for _, e := range s {
		saddr := lookup(e.LocalAddr)
		daddr := lookup(e.RemoteAddr)
		fmt.Printf("%s %s %s \n", saddr, daddr, e.State)
	}
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
		cur_c := string(daddr+":"+strconv.Itoa(int(dport)) + "|" + saddr + strconv.Itoa(int(sport)))
		//fmt.Printf("%s %s %s \n", saddr, daddr, e.State)
		if((*port==-1||*port==int(sport))&&sport!=myport) {
			connection_counts_new[cur_c] = connection_counts[cur_c] + 1
		}
	}
	return(connection_counts_new)
}
func main() {
	flag.Parse()

	/*if !*udp && !*tcp {
		flag.Usage()
		os.Exit(0)
	}*/

	//displaySockInfo(socks)
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

			if(*tcpv6) {
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
	go func() {
		for{
			var sum uint = 0
			for _, value := range(connections) {
				if(value > 6) {
					sum += 1
				}
			}
			if(*tcpv6) {
				for _, value := range(connections6) {
					if(value > 6) {
						sum += 1
					}
				}
			}
			time.Sleep(1 * time.Second)
			netstat_tcp_connection_longterm_counts.Set(float64(sum))
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":"+strconv.Itoa(myport) , nil)
}
