package main

import (
	"flag"
	"os"
	//"strings"
	"net/http"
	"time"
	"encoding/base64"

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
	nexusUser				string
	nexusPass				string
	auth					string
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
	ch <- e.availableProcessors
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	//var err error
	defer func(begun time.Time) {
		e.availableProcessors.Set(time.Since(begun).Seconds())
	}(time.Now())
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh

}

func main() {
	flag.Parse()
	log.Infoln("Starting nexus_exporter " + version)
	nexusUrl := os.Getenv("NEXUS_URL")
	nexusPath := os.Getenv("NEXUS_PATH")
	nexusUser := os.Getenv("NEXUS_USER")
	nexusPass := os.Getenv("NEXUS_PASS")


	client := &http.Client{}
	req, err:= http.NewRequest("GET", nexusUrl+nexusPath, nil)
	req.SetBasicAuth(nexusUser, nexusPass)
	resp, err := Client.Do(req)

	if err != nil{
        log.Fatal(err)
	}
	
	log.Infoln(resp)
	
	
	exporter := NewExporter(nexusUrl, nexusPath)
	
	prometheus.MustRegister(exporter)
	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}