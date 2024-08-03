package analytics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Analyze struct {
	RequestCounter *prometheus.CounterVec
}

func RunAnalyze() *Analyze {
	var (
		requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests",
		}, []string{"method", "path"})
	)

	return &Analyze{RequestCounter: requestCounter}

}

func (a *Analyze) Init() {
	prometheus.MustRegister(a.RequestCounter)
}

func (a *Analyze) Handler(w http.ResponseWriter, r *http.Request) {
	a.RequestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()
}
