package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	namespace = ""
	subsystem string

	reqLabels = []string{"endpoint", "method"}
)

var (
	reqCnt *prometheus.CounterVec
)

func init() {
	registerMetrics("secret_service")
}

func registerMetrics(service string) {
	subsystem = service

	reqCntOpts := prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_count",
		Help:      "Counter of requests received into the system.",
	}
	reqCnt = register(prometheus.NewCounterVec(reqCntOpts, reqLabels), "").(*prometheus.CounterVec)
}

func register(c prometheus.Collector, name string) prometheus.Collector {
	err := prometheus.Register(c)
	if err != nil {
		if regErr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return regErr.ExistingCollector
		}
		log.Panicf("metric '%s' registration failed with error, %v", name, err)
	}
	log.Debugf("metric '%s' registered", name)
	return c
}

func recordRequestCounter(method, endpoint string) {
	reqCnt.WithLabelValues(method, endpoint).Inc()
}
