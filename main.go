package main

import (
	"flag"
	"os"
	//"strings"
	"net/http"
	"time"
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)
func main() {
	nexusUrl := os.Getenv("NEXUS_URL")
	nexusPath := os.Getenv("NEXUS_PATH")
	nexusUser := os.Getenv("NEXUS_USER")
	nexusPass := os.Getenv("NEXUS_PASS")

	var (
		version		= "1.0.0.dev"
		listen		= flag.String("web.listen-address",":8080", "Addressto listen")
		metricPath	= flag.String("web.telemetry-path","/metrics","Path of the metrics")
		landingPage	= []byte("<html><head><title>Nexus-Exporter</title></head><h1>NEXUS EXPORTER "+version+"</h1>")
	)

	flag.Parse()

	nexusUrl, ok := os.LookupEnv("NEXUS_URL")
	if ok {
		*nexusUrl = nexusUrl
	}
	nexusPath, ok := os.LookupEnv("NEXUS_URL")
	if ok {
		*nexusPath = nexusPath
	}
	nexusUser, ok := os.LookupEnv("NEXUS_URL")
	if ok {
		*nexusUser = nexusUser
	}
	nexusPass, ok := os.LookupEnv("NEXUS_URL")
	if ok {
		*nexusPass = nexusPass
	}

	data = getMetrics();

	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func getMetrics(){
	log.Infoln(*nexusPass)
}