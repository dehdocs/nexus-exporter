package main

import (
	"flag"
	//"os"
	//"strings"
	//"time"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/common/log"
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
	dsn			string
	availableProcessors 	prometheus.Gauge
	freeMemory         	 prometheus.Gauge
	totalMemory         	prometheus.Gauge
	maxMemory           	prometheus.Gauge
	threads             	prometheus.Gauge
}

func main() {
	flag.Parse()
	log.Infoln("Starting nexus_exporter " + version)
	dsn := os.Getenv("DATA_SOURCE_NAME")
	exporter := NewExporter(dsn)
	prometheus.MustRegister(exporter)
	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
