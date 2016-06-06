package main

import (
	"flag"
	"fmt"
	"net/http"
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
		Name: "kemp_connections_per_second",
		Help: "The number of connections per second.",
	})
	bytesPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_bytes_per_second",
		Help: "The number of bytes per second.",
	})
	packetsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_packets_per_second",
		Help: "The number of packets per second.",
	})
)

func init() {
	flag.StringVar(&username, "username", "", "username to connect to the Kemp API")
	flag.StringVar(&password, "password", "", "password to connect to the Kemp API")
	flag.StringVar(&endpoint, "endpoint", "", "API endpoint to connect to")
	flag.BoolVar(&debug, "debug", false, "enable debug output")

	prometheus.MustRegister(connsPerSec)
	prometheus.MustRegister(bytesPerSec)
	prometheus.MustRegister(packetsPerSec)
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

			time.Sleep(WaitTime)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	http.Handle("/metrics", prometheus.Handler())

	http.ListenAndServe(":8000", nil)
}
