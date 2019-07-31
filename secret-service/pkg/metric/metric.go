package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	namespace = ""
	subsystem string

	reqLabels          = []string{"endpoint", "method"}
	respTimeObjectives = map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001}
)

var (
	reqCnt   *prometheus.CounterVec
	respTime *prometheus.SummaryVec
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
	reqCnt = register(prometheus.NewCounterVec(reqCntOpts, reqLabels), "request_count").(*prometheus.CounterVec)

	respTimeOpts := prometheus.SummaryOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "response_time",
		Help:       "Bucketed summary of response time.",
		Objectives: respTimeObjectives,
	}
	respTime = register(prometheus.NewSummaryVec(respTimeOpts, reqLabels), "response_time").(*prometheus.SummaryVec)
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

func recordRequestCounter(endpoint, method string) {
	reqCnt.WithLabelValues(endpoint, method).Inc()
}

func recordResponseTime(endpoint, method string, startTime time.Time) {
	respTime.WithLabelValues(endpoint, method).Observe(time.Since(startTime).Seconds())
}
