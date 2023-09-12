package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			humidity.Set(62.6)
			temperature.Set(35.0)
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	humidity = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "humidity_actual",
		Help: "The value of humidity",
	})
	temperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "temperature_actual",
		Help: "The value of temperature",
	})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
