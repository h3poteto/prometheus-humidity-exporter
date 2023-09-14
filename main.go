package main

import (
	"log"
	"net/http"
	"time"

	"github.com/h3poteto/prometheus-humidity-exporter/dht20"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() error {
	sensor, err := dht20.New()
	if err != nil {
		return err
	}
	go func(sensor *dht20.DHT20) {
		defer sensor.Clean()
		for {
			hum, tmp, err := sensor.Get()
			if err != nil {
				log.Printf("[error] %v", err)
			}
			humidity.Set(hum)
			temperature.Set(tmp)
			time.Sleep(30 * time.Second)
		}
	}(sensor)

	return nil
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
