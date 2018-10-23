package main

import (
	"flag"
	"os"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
	//"strings"
	//"time"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)


func main() {


	var (
		listen			= flag.String("web.listen-address",":8080", "Addressto listen")
		metricPath		= flag.String("web.telemetry-path","/metrics","Path of the metrics")
		landingPage		= []byte("<html><head><title>Nexus-Exporter</title></head><h1>NEXUS EXPORTER </h1>")
		nUrl			= flag.String("nexus.uri", "http://10.129.176.139:8081", "HTTP API address of nexus.")
		nPath			= flag.String("nexus.path", "/service/siesta/atlas/system-information", "nexus api path.")
		nUser			= flag.String("nexus.user", "admin", "nexus password.")
		nPass			= flag.String("nexus.pass", "admin123", "nexus password.")
		data 			string
		jsonData 			map[string]interface{}

		threads = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexus_threads",
			Help: "Quantity of threads are available in nexus.",
		})
		availableProcessors = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexus_processors_available",
			Help: "Quantity of processors are available in nexus.",
		})
		freeMemory = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexus_freeMemory",
			Help: "Quantity of free memory are available in nexus.",
		})
		totalMemory = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexus_totalMemory",
			Help: "Quantity of total memory are available in nexus.",
		})
		maxMemory = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexus_maxMemory",
			Help: "Quantity of max memory in nexus.",
		})
	)

	flag.Parse()

	nexusUrl, ok := os.LookupEnv("NEXUS_URL")
	if ok {
		*nUrl = nexusUrl
	}
	nexusPath, ok := os.LookupEnv("NEXUS_PATH")
	if ok {
		*nPath = nexusPath
	}
	nexusUser, ok := os.LookupEnv("NEXUS_USER")
	if ok {
		*nUser = nexusUser
	}
	nexusPass, ok := os.LookupEnv("NEXUS_PASS")
	if ok {
		*nPass = nexusPass
	}

	
	data = getMetrics(nexusUrl, nexusPath, nexusUser, nexusPass);

	

	json.Unmarshal([]byte(data), &jsonData)
	runtime := jsonData["system-runtime"].(map[string]interface{})
	

	prometheus.MustRegister(threads)
	prometheus.MustRegister(availableProcessors)
	prometheus.MustRegister(freeMemory)
	prometheus.MustRegister(totalMemory)
	prometheus.MustRegister(maxMemory)
	
	threads.Set(runtime["threads"].(float64))
	availableProcessors.Set(runtime["availableProcessors"].(float64))
	freeMemory.Set(runtime["freeMemory"].(float64))
	totalMemory.Set(runtime["totalMemory"].(float64))
	maxMemory.Set(runtime["maxMemory"].(float64))


	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func getMetrics(url string, path string, user string, pass string) string {	
	client := &http.Client{}
	req, err:= http.NewRequest("GET", url+path, nil)
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)

	if err != nil{
        log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	return string(bodyText)
}