package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	kemp "github.com/giantswarm/kemp-client"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	WaitTime = time.Second * 10
)

var (
	username string
	password string
	endpoint string
	debug    bool

	connsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_totals_connections_per_second",
		Help: "The number of connections per second.",
	})
	bytesPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_totals_bytes_per_second",
		Help: "The number of bytes per second.",
	})
	packetsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_totals_packets_per_second",
		Help: "The number of packets per second.",
	})

	virtualServerTotalConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_total_connections",
		Help: "The number of total connections per virtual server.",
	}, []string{"address", "port"})
	virtualServerTotalPackets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_total_packets",
		Help: "The number of total packets per virtual server.",
	}, []string{"address", "port"})
	virtualServerTotalBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_total_bytes",
		Help: "The number of total bytes per virtual server.",
	}, []string{"address", "port"})
	virtualServerActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_active_connections",
		Help: "The number of active connections per virtual server.",
	}, []string{"address", "port"})
	virtualServerConnsPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_connections_per_second",
		Help: "The number of connections per second per virtual server.",
	}, []string{"address", "port"})
	virtualServerBytesRead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_bytes_read",
		Help: "The number of bytes read per virtual server.",
	}, []string{"address", "port"})
	virtualServerBytesWritten = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_server_bytes_written",
		Help: "The number of bytes written per virtual server",
	}, []string{"address", "port"})

	realServerTotalConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_total_connections",
		Help: "The number of total connections per real server.",
	}, []string{"address", "port"})
	realServerTotalPackets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_total_packets",
		Help: "The number of total packets per real server.",
	}, []string{"address", "port"})
	realServerTotalBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_total_bytes",
		Help: "The number of total bytes per real server.",
	}, []string{"address", "port"})
	realServerActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_active_connections",
		Help: "The number of active connections per real server.",
	}, []string{"address", "port"})
	realServerConnsPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_connections_per_second",
		Help: "The number of connections per second per real server.",
	}, []string{"address", "port"})
	realServerBytesRead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_bytes_read",
		Help: "The number of bytes read per real server.",
	}, []string{"address", "port"})
	realServerBytesWritten = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_real_server_bytes_written",
		Help: "The number of bytes written per real server",
	}, []string{"address", "port"})
)

func init() {
	flag.StringVar(&username, "username", "", "username to connect to the Kemp API")
	flag.StringVar(&password, "password", "", "password to connect to the Kemp API")
	flag.StringVar(&endpoint, "endpoint", "", "API endpoint to connect to")
	flag.BoolVar(&debug, "debug", false, "enable debug output")

	prometheus.MustRegister(connsPerSec)
	prometheus.MustRegister(bytesPerSec)
	prometheus.MustRegister(packetsPerSec)

	prometheus.MustRegister(virtualServerTotalConnections)
	prometheus.MustRegister(virtualServerTotalPackets)
	prometheus.MustRegister(virtualServerTotalBytes)
	prometheus.MustRegister(virtualServerActiveConnections)
	prometheus.MustRegister(virtualServerConnsPerSec)
	prometheus.MustRegister(virtualServerBytesRead)
	prometheus.MustRegister(virtualServerBytesWritten)

	prometheus.MustRegister(realServerTotalConnections)
	prometheus.MustRegister(realServerTotalPackets)
	prometheus.MustRegister(realServerTotalBytes)
	prometheus.MustRegister(realServerActiveConnections)
	prometheus.MustRegister(realServerConnsPerSec)
	prometheus.MustRegister(realServerBytesRead)
	prometheus.MustRegister(realServerBytesWritten)
}

func main() {
	flag.Parse()

	client := kemp.NewClient(kemp.Config{
		User:     username,
		Password: password,
		Endpoint: endpoint,
		Debug:    debug,
	})

	go func() {
		for {
			statistics, err := client.GetStatistics()
			if err != nil {
				fmt.Println(err)
			}

			connsPerSec.Set(float64(statistics.Totals.ConnectionsPerSec))
			bytesPerSec.Set(float64(statistics.Totals.BytesPerSec))
			packetsPerSec.Set(float64(statistics.Totals.PacketsPerSec))

			for _, vs := range statistics.VirtualServers {
				virtualServerTotalConnections.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalConnections))
				virtualServerTotalPackets.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalPackets))
				virtualServerTotalBytes.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalBytes))
				virtualServerActiveConnections.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.ActiveConnections))
				virtualServerConnsPerSec.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.ConnectionsPerSec))
				virtualServerBytesRead.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.BytesRead))
				virtualServerBytesWritten.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.BytesWritten))
			}

			for _, rs := range statistics.RealServers {
				realServerTotalConnections.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.TotalConnections))
				realServerTotalPackets.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.TotalPackets))
				realServerTotalBytes.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.TotalBytes))
				realServerActiveConnections.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.ActiveConnections))
				realServerConnsPerSec.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.ConnectionsPerSec))
				realServerBytesRead.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.BytesRead))
				realServerBytesWritten.WithLabelValues(rs.Address, strconv.Itoa(rs.Port)).Set(float64(rs.BytesWritten))
			}

			time.Sleep(WaitTime)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	http.Handle("/metrics", prometheus.Handler())

	http.ListenAndServe(":8000", nil)
}
