package xTool

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MetricHttpRequestTotal 请求数统计
	MetricHttpRequestTotal *prometheus.CounterVec
)

// InitPrometheusMetric 初始化 普罗米修斯 指标
func InitPrometheusMetric(namespace, subsystem string, collector ...prometheus.Collector) {
	MetricHttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_request_total",
			Help:      "http request total",
		},
		[]string{"from"},
	)
	collector = append(collector, MetricHttpRequestTotal)
	// 注册
	prometheus.MustRegister(collector...)
}

// PrometheusHttp 普罗米修斯 HTTP 入口注册到 GF
func PrometheusHttp(r *ghttp.Request) {
	promhttp.Handler().ServeHTTP(r.Response.Writer, r.Request)
}
