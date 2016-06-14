package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kemp "github.com/giantswarm/kemp-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server [endpoint] [username] [password]",
		Short: "Start the HTTP server",
		Run:   serverRun,
	}

	debug       bool
	waitSeconds int
	port        int

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

	virtualServiceTotalConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_connections",
		Help: "The number of total connections per virtual service.",
	}, []string{"address", "port"})
	virtualServiceTotalPackets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_packets",
		Help: "The number of total packets per virtual service.",
	}, []string{"address", "port"})
	virtualServiceTotalBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_bytes",
		Help: "The number of total bytes per virtual service.",
	}, []string{"address", "port"})
	virtualServiceActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_active_connections",
		Help: "The number of active connections per virtual service.",
	}, []string{"address", "port"})
	virtualServiceConnsPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_connections_per_second",
		Help: "The number of connections per second per virtual service.",
	}, []string{"address", "port"})
	virtualServiceBytesRead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_bytes_read",
		Help: "The number of bytes read per virtual service.",
	}, []string{"address", "port"})
	virtualServiceBytesWritten = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_bytes_written",
		Help: "The number of bytes written per virtual service",
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
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&port, "port", 8000, "port to listen on")
	serverCmd.Flags().IntVar(&waitSeconds, "wait", 10, "time (in seconds) between accessing the Kemp API")
	serverCmd.Flags().BoolVar(&debug, "debug", false, "enable debug output")

	prometheus.MustRegister(connsPerSec)
	prometheus.MustRegister(bytesPerSec)
	prometheus.MustRegister(packetsPerSec)

	prometheus.MustRegister(virtualServiceTotalConnections)
	prometheus.MustRegister(virtualServiceTotalPackets)
	prometheus.MustRegister(virtualServiceTotalBytes)
	prometheus.MustRegister(virtualServiceActiveConnections)
	prometheus.MustRegister(virtualServiceConnsPerSec)
	prometheus.MustRegister(virtualServiceBytesRead)
	prometheus.MustRegister(virtualServiceBytesWritten)

	prometheus.MustRegister(realServerTotalConnections)
	prometheus.MustRegister(realServerTotalPackets)
	prometheus.MustRegister(realServerTotalBytes)
	prometheus.MustRegister(realServerActiveConnections)
	prometheus.MustRegister(realServerConnsPerSec)
	prometheus.MustRegister(realServerBytesRead)
	prometheus.MustRegister(realServerBytesWritten)
}

func serverRun(cmd *cobra.Command, args []string) {
	flag.Parse()

	if len(cmd.Flags().Args()) != 3 {
		cmd.Help()
		os.Exit(1)
	}

	client := kemp.NewClient(kemp.Config{
		Endpoint: flag.Arg(1),
		User:     flag.Arg(2),
		Password: flag.Arg(3),
		Debug:    debug,
	})

	go func() {
		for {
			statistics, err := client.GetStatistics()
			if err != nil {
				log.Println("Error getting statistics ", err)
				os.Exit(1)
			}

			connsPerSec.Set(float64(statistics.Totals.ConnectionsPerSec))
			bytesPerSec.Set(float64(statistics.Totals.BytesPerSec))
			packetsPerSec.Set(float64(statistics.Totals.PacketsPerSec))

			for _, vs := range statistics.VirtualServices {
				virtualServiceTotalConnections.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalConnections))
				virtualServiceTotalPackets.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalPackets))
				virtualServiceTotalBytes.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.TotalBytes))
				virtualServiceActiveConnections.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.ActiveConnections))
				virtualServiceConnsPerSec.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.ConnectionsPerSec))
				virtualServiceBytesRead.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.BytesRead))
				virtualServiceBytesWritten.WithLabelValues(vs.Address, strconv.Itoa(vs.Port)).Set(float64(vs.BytesWritten))
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

			time.Sleep(time.Second * time.Duration(waitSeconds))
		}
	}()

	go func() {
		intChan := make(chan os.Signal)
		termChan := make(chan os.Signal)

		signal.Notify(intChan, syscall.SIGINT)
		signal.Notify(termChan, syscall.SIGTERM)

		select {
		case <-intChan:
			log.Print("Received SIGINT, exiting")
			os.Exit(0)
		case <-termChan:
			log.Print("Received SIGTERM, exiting")
			os.Exit(0)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	http.Handle("/metrics", prometheus.Handler())

	log.Print("Listening on port ", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
