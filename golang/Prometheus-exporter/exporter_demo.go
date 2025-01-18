package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// 定义一个Prometheus指标来存储CPU利用率
var cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_usage_percent",
	Help: "CPU usage in percent",
})

// 定义一个Prometheus指标来存储内存使用率
var memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "memory_usage_percent",
	Help: "Memory usage in percent",
})

func init() {
	// 创建并注册指标到自定义注册表
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
}

func updateMetrics() {
	// 获取CPU利用率
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
	} else if len(percent) > 0 {
		cpuUsage.Set(percent[0])
	}

	// 获取内存使用率
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
	} else {
		memoryUsage.Set(v.UsedPercent)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	updateMetrics()

	// 使用自定义注册表创建HTTP处理器
	// handler := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
	// handler.ServeHTTP(w, r)
	promhttp.Handler().ServeHTTP(w, r)
}

func main() {

	// 设置HTTP处理器来暴露Prometheus指标
	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Starting server on :9100")
	log.Fatal(http.ListenAndServe(":9100", nil))
}
