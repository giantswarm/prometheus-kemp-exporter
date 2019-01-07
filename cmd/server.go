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

	kemp "github.com/majandres/kemp-client"
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

	kempUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_up",
		Help: "Whether the last kemp scrape was successful (1: up, 0: down)",
	})
	
	tpsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_total_tps",
		Help: "The total number of Transactions Per Second (TPS).",
	})
	tpsTotalSSL = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_total_tps_ssl",
		Help: "The total number of SSL Transactions Per Second (TPS).",
	})
	cpuUser = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_cpu_user",
		Help: "The percentage of the CPU spent processing in user mode.",
	})
	cpuSystem = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_cpu_system",
		Help: "The percentage of the CPU spent processing in system mode.",
	})
	cpuIdle = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_cpu_idle",
		Help: "The percentage of CPU which is idle.",
	})
	cpuIOWaiting = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_cpu_io_waiting",
		Help: "The percentage of the CPU spent waiting for I/O to complete.",
	})
	// Used - This will be 100 - Idle
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_cpu_percentage_used",
		Help: "The percentage of the CPU that has been utilized.",
	})
	memUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_memory_amount_used",
		Help: "The amount of memory in use.",
	})
	percentMemUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_percent_memory_used",
		Help: "The percentage of memory used.",
	})
	memFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_memory_amount_free",
		Help: "The amount of memory free.",
	})
	percentMemFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kemp_percent_memory_free",
		Help: "The percentage of free memory.",
	})	

        connsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_connections_per_second",
                Help: "The number of connections per second.",
        })
        totalConns = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_connections",
                Help: "The total number of connections made.",
        })
        bitsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_bits_per_second",
                Help: "The number of bits per second.",
        })  
        totalBits = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_bits",
                Help: "The total number of bits.",
        })  
        bytesPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_bytes_per_second",
                Help: "The number of bytes per second.",
        })  
        totalBytes = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_bytes",
                Help: "The total number of bytes.",
        })  
        packetsPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_packets_per_second",
                Help: "The number of packets per second.",
        })  
        totalPackets = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "kemp_totals_packets",
                Help: "The total number of packets.",
        })

	virtualServiceTotalConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_connections",
		Help: "The number of total connections per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceTotalPackets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_packets",
		Help: "The number of total packets per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceTotalBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_total_bytes",
		Help: "The number of total bytes per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_active_connections",
		Help: "The number of active connections per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceConnsPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_connections_per_second",
		Help: "The number of connections per second per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceBytesRead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_bytes_read",
		Help: "The number of bytes read per virtual service.",
	}, []string{"name", "address", "port"})
	virtualServiceBytesWritten = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kemp_virtual_service_bytes_written",
		Help: "The number of bytes written per virtual service",
	}, []string{"name", "address", "port"})

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

	prometheus.MustRegister(kempUp)

	prometheus.MustRegister(tpsTotal)
	prometheus.MustRegister(tpsTotalSSL)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuIdle)
	prometheus.MustRegister(cpuIOWaiting)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memUsed)
	prometheus.MustRegister(percentMemUsed)
	prometheus.MustRegister(memFree)
	prometheus.MustRegister(percentMemFree)	

        prometheus.MustRegister(connsPerSec)
        prometheus.MustRegister(totalConns)
        prometheus.MustRegister(bitsPerSec)
        prometheus.MustRegister(totalBits)
        prometheus.MustRegister(bytesPerSec)
        prometheus.MustRegister(totalBytes)
        prometheus.MustRegister(packetsPerSec)
        prometheus.MustRegister(totalPackets)

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

func sleep() {
	time.Sleep(time.Second * time.Duration(waitSeconds))
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
				kempUp.Set(0)

				sleep()
				continue
			}

			// The Kemp API stats endpoint doesn't return virtual service names.
			// So, let's construct a lookup table between addresses and names,
			// so we can use as the names as a label.
			addressNameLookup := map[string]string{}
			virtualServices, err := client.ListVirtualServices()
			if err != nil {
				log.Println("Error getting virtual services ", err)
				kempUp.Set(0)

				sleep()
				continue
			}
			kempUp.Set(1)

			tpsTotal.Set(float64(statistics.TPS.Total))
			tpsTotalSSL.Set(float64(statistics.TPS.SSL))
			cpuUser.Set(float64(statistics.CPU.User))
			cpuSystem.Set(float64(statistics.CPU.System))
			cpuIdle.Set(float64(statistics.CPU.Idle))
			cpuIOWaiting.Set(float64(statistics.CPU.IOWaiting))
			cpuUsage.Set(float64(statistics.CPU.Used))
			memUsed.Set(float64(statistics.Memory.MemUsed))
			percentMemUsed.Set(float64(statistics.Memory.PercentMemUsed))
			memFree.Set(float64(statistics.Memory.MemFree))
			percentMemFree.Set(float64(statistics.Memory.PercentMemFree))
			
			for _, vs := range virtualServices {
				addressNameLookup[vs.IPAddress] = vs.Name
			}

                        connsPerSec.Set(float64(statistics.Totals.ConnectionsPerSec))
                        totalConns.Set(float64(statistics.Totals.TotalConnections))
                        bitsPerSec.Set(float64(statistics.Totals.BitsPerSec))
                        totalBits.Set(float64(statistics.Totals.TotalBits))
                        bytesPerSec.Set(float64(statistics.Totals.BytesPerSec))
                        totalBytes.Set(float64(statistics.Totals.TotalBytes))
                        packetsPerSec.Set(float64(statistics.Totals.PacketsPerSec))
                        totalPackets.Set(float64(statistics.Totals.TotalPackets))

			for _, vs := range statistics.VirtualServices {
				name, _ := addressNameLookup[vs.Address]
				port := strconv.Itoa(vs.Port)

				virtualServiceTotalConnections.WithLabelValues(name, vs.Address, port).Set(float64(vs.TotalConnections))
				virtualServiceTotalPackets.WithLabelValues(name, vs.Address, port).Set(float64(vs.TotalPackets))
				virtualServiceTotalBytes.WithLabelValues(name, vs.Address, port).Set(float64(vs.TotalBytes))
				virtualServiceActiveConnections.WithLabelValues(name, vs.Address, port).Set(float64(vs.ActiveConnections))
				virtualServiceConnsPerSec.WithLabelValues(name, vs.Address, port).Set(float64(vs.ConnectionsPerSec))
				virtualServiceBytesRead.WithLabelValues(name, vs.Address, port).Set(float64(vs.BytesRead))
				virtualServiceBytesWritten.WithLabelValues(name, vs.Address, port).Set(float64(vs.BytesWritten))
			}

			for _, rs := range statistics.RealServers {
				port := strconv.Itoa(rs.Port)

				realServerTotalConnections.WithLabelValues(rs.Address, port).Set(float64(rs.TotalConnections))
				realServerTotalPackets.WithLabelValues(rs.Address, port).Set(float64(rs.TotalPackets))
				realServerTotalBytes.WithLabelValues(rs.Address, port).Set(float64(rs.TotalBytes))
				realServerActiveConnections.WithLabelValues(rs.Address, port).Set(float64(rs.ActiveConnections))
				realServerConnsPerSec.WithLabelValues(rs.Address, port).Set(float64(rs.ConnectionsPerSec))
				realServerBytesRead.WithLabelValues(rs.Address, port).Set(float64(rs.BytesRead))
				realServerBytesWritten.WithLabelValues(rs.Address, port).Set(float64(rs.BytesWritten))
			}

			sleep()
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
