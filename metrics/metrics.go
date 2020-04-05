package metrics

import (
	"breeze/sensor"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Collector struct {
	temperature   prometheus.Gauge
	thermalSensor sensor.Thermal
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("temperature", "thermal sensor metrics", nil, nil)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	read, err := c.thermalSensor.Read()
	if err != nil {
		log.Warnf("collector: unable to read temperature: %s", err.Error())
		return
	}

	c.temperature.Set(read)
	ch <- c.temperature
}

func New(thermalSensor sensor.Thermal) *Collector {
	return &Collector{
		thermalSensor: thermalSensor,
		temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "temperature",
		}),
	}
}

func (c *Collector) Run(addr string, port int) error {
	if port == 0 {
		port = 9999
	}

	if addr == "" {
		addr = "0.0.0.0"
	}

	if err := prometheus.Register(c); err != nil {
		return err
	}

	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
}