package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mindprince/gonvml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "nvidia_gpu"
)

var (
	addr = flag.String("web.listen-address", ":9445", "Address to listen on for web interface and telemetry.")

	labels = []string{"minor_number", "uuid", "name"}
)

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.numDevices.Desc()
	c.usedMemory.Describe(ch)
	c.totalMemory.Describe(ch)
	c.dutyCycle.Describe(ch)
	c.powerUsage.Describe(ch)
	c.temperature.Describe(ch)
	c.fanSpeed.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.Lock()
	defer c.Unlock()

	c.usedMemory.Reset()
	c.totalMemory.Reset()
	c.dutyCycle.Reset()
	c.powerUsage.Reset()
	c.temperature.Reset()
	c.fanSpeed.Reset()

	numDevices, err := gonvml.DeviceCount()
	if err != nil {
		log.Printf("DeviceCount() error: %v", err)
		return
	} else {
		c.numDevices.Set(float64(numDevices))
		ch <- c.numDevices
	}

	for i := 0; i < int(numDevices); i++ {
		dev, err := gonvml.DeviceHandleByIndex(uint(i))
		if err != nil {
			log.Printf("DeviceHandleByIndex(%d) error: %v", i, err)
			continue
		}

		minorNumber, err := dev.MinorNumber()
		if err != nil {
			log.Printf("MinorNumber() error: %v", err)
			continue
		}

	}
	c.usedMemory.Collect(ch)
	c.totalMemory.Collect(ch)
	c.dutyCycle.Collect(ch)
	c.powerUsage.Collect(ch)
	c.temperature.Collect(ch)
	c.fanSpeed.Collect(ch)
}

func main() {
	flag.Parse()

	if err := gonvml.Initialize(); err != nil {
		log.Fatalf("Couldn't initialize gonvml: %v. Make sure NVML is in the shared library search path.", err)
	}
	defer gonvml.Shutdown()

	if driverVersion, err := gonvml.SystemDriverVersion(); err != nil {
		log.Printf("SystemDriverVersion() error: %v", err)
	} else {
		log.Printf("SystemDriverVersion(): %v", driverVersion)
	}

	prometheus.MustRegister(NewCollector())

	log.Fatalf("ListenAndServe error: %v", http.ListenAndServe(*addr, promhttp.Handler()))
}
