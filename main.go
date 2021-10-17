package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	speedtest "github.com/cloudziu/ookla-speedtest-prometheus/speedtest_ookla"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promversion "github.com/prometheus/common/version"
)

const (
	namespace = "speedtest"
	subsystem = "network"
)

var (

	//go:embed html/index.html
	defaultPage []byte

	// EXPORTER_PORT default prometheus exporter
	EXPORTER_PORT    = "9000" 

	// METRICS_ENDPOINT default prometheus metrics path
	METRICS_ENDPOINT = "/metrics"

	download = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "download"),
		"Download bandwidth (Mbps).",
		nil, nil,
	)
	upload = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "upload"),
		"Upload bandwidth (Mbps).",
		nil, nil,
	)
	latency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "latency"),
		"Latency (ms)",
		nil, nil,
	)
	jitter = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "jitter"),
		"Jitter (ms).",
		nil, nil,
	)
	packetLost = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "packetLost"),
		"Packet lost",
		nil, nil,
	)
)

//  Speedtest type that implements prometheus.Collector.
type Speedtest struct {
	*speedtest.Client
}

// NewSpeedtest return new speedtest client type
func NewSpeedtest() *Speedtest {
	return &Speedtest{
		&speedtest.Client{},
	}
}

// Describe describes all the metrics ever exported by the Speedtest exporter.
// It implements prometheus.Collector.
func (s *Speedtest) Describe(ch chan<- *prometheus.Desc) {
	ch <- download
	ch <- upload
	ch <- latency
	ch <- jitter
	ch <- packetLost
}

// Collect fetches the stats from configured Speedtest location
// and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (s *Speedtest) Collect(ch chan<- prometheus.Metric) {
	if s == (&Speedtest{}) {
		log.Fatal("Speedtest client not configured")
		return
	}
	// logic, execute speed test
	s.Client.RunSpeedtest()

	ch <- prometheus.MustNewConstMetric(download, prometheus.GaugeValue, s.Client.Download)
	ch <- prometheus.MustNewConstMetric(upload, prometheus.GaugeValue, s.Client.Upload)
	ch <- prometheus.MustNewConstMetric(latency, prometheus.GaugeValue, s.Client.Latency)
	ch <- prometheus.MustNewConstMetric(jitter, prometheus.GaugeValue, s.Client.Jitter)
	ch <- prometheus.MustNewConstMetric(packetLost, prometheus.GaugeValue, s.Client.PacketLoss)

}

func config() {
	if port, set := os.LookupEnv("SPEEDTEST_EXPORTER_PORT"); set {
		EXPORTER_PORT = port
	} else {
		log.Printf("Environment variable SPEEDTEST_EXPORTER_PORT not set, using default port: %s",
			EXPORTER_PORT)
	}

}

func init() {
	speedtest.InitialCheck()
	config()
	prometheus.MustRegister(promversion.NewCollector("speedtest_ookla"))
}

func main() {

	speedtestClient := NewSpeedtest()
	prometheus.MustRegister(speedtestClient)

	log.Printf("Start speedtest network collector %s", promversion.Info())
	log.Printf("Build context: %s", promversion.BuildContext())

	http.Handle(METRICS_ENDPOINT, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter,
		r *http.Request) {
		w.Write(defaultPage)
	})
	log.Printf("Start handler on port: %s", EXPORTER_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", EXPORTER_PORT), nil))
}
