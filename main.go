package main

import (
	"flag"
	"os"
	//"strings"
	"net/http"
	//"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	version		= "1.0.0.dev"
	listen		= flag.String("web.listen-address",":8080", "Addressto listen")
	metricPath	= flag.String("web.telemetry-path","/metrics","Path of the metrics")
	landingPage	= []byte("<html><head><title>Nexus-Exporter</title></head><h1>NEXUS EXPORTER "+version+"</h1>")
)

const (
	namespace	= "nexus"
	exporter	= "exporter"
)

type Exporter struct {
	nexusUrl				string
	nexusPath				string
	availableProcessors		prometheus.Gauge
	//freeMemory				prometheus.Gauge
	//totalMemory				prometheus.Gauge
	//maxMemory				prometheus.Gauge
	//threads					prometheus.Gauge
}
func NewExporter(nexusUrl string, nexusPath string) *Exporter {
	return &Exporter{
		nexusUrl: nexusUrl,
		nexusPath: nexusPath,
		availableProcessors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "nexus_processors_available",
			Help:      "Quantity of processors are available in nexus.",
		}),
	}
}
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.duration
	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
	ch <- e.up
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	var err error
	defer func(begun time.Time) {
		e.duration.Set(time.Since(begun).Seconds())
		if err == nil {
			e.error.Set(0)
		} else {
			e.error.Set(1)
		}
	}(time.Now())

	db, err := sql.Open("oci8", e.dsn)
	if err != nil {
		log.Errorln("Error opening connection to database:", err)
		return
	}
	defer db.Close()

	isUpRows, err := db.Query("SELECT 1 FROM DUAL")
	if err != nil {
		log.Errorln("Error pinging oracle:", err)
		e.up.Set(0)
		return
	}
	isUpRows.Close()
	e.up.Set(1)

	if err = ScrapeActivity(db, ch); err != nil {
		log.Errorln("Error scraping for activity:", err)
		e.scrapeErrors.WithLabelValues("activity").Inc()
	}

	if err = ScrapeTablespace(db, ch); err != nil {
		log.Errorln("Error scraping for tablespace:", err)
		e.scrapeErrors.WithLabelValues("tablespace").Inc()
	}

	if err = ScrapeWaitTime(db, ch); err != nil {
		log.Errorln("Error scraping for wait_time:", err)
		e.scrapeErrors.WithLabelValues("wait_time").Inc()
	}

	if err = ScrapeSessions(db, ch); err != nil {
		log.Errorln("Error scraping for sessions:", err)
		e.scrapeErrors.WithLabelValues("sessions").Inc()
	}

	if err = ScrapeProcesses(db, ch); err != nil {
		log.Errorln("Error scraping for process:", err)
		e.scrapeErrors.WithLabelValues("process").Inc()
	}

}

func main() {
	flag.Parse()
	log.Infoln("Starting nexus_exporter " + version)
	nexusUrl := os.Getenv("NEXUS_URL")
	nexusPath := os.Getenv("NEXUS_PATH")
	exporter := NewExporter(nexusUrl, nexusPath)
	prometheus.MustRegister(exporter)
	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}