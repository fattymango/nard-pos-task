package metrics

import (
	"context"
	"log"
	"multitenant/pkg/config"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	MONITORING_INTERVAL = 5 * time.Second
)

// Adding new metrics for CPU and Memory
type Metrics struct {
	config                *config.Config
	ctx                   context.Context
	apiRequestDuration    *prometheus.HistogramVec
	transactionsProcessed *prometheus.CounterVec
	cacheHits             prometheus.Counter
	cacheMisses           prometheus.Counter
	cpuUsage              *prometheus.GaugeVec
	memoryUsage           prometheus.Gauge
}

const (
	TX_PROCCESSED_STATUS_SUCCESS = "tx_success"
	TX_PROCCESSED_STATUS_FAILED  = "tx_failed"
)

func NewMetrics(ctx context.Context, config *config.Config) *Metrics {
	apiRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Histogram of response time for API requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)

	transactionsProcessed := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transactions_processed_total",
			Help: "Total number of transactions processed",
		},
		[]string{"status"},
	)

	cacheHits := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Adding CPU and Memory metrics
	cpuUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percentage",
			Help: "Percentage of CPU usage",
		},
		[]string{"cpu"},
	)

	memoryUsage := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Total memory usage in bytes",
		},
	)

	// Register all metrics
	prometheus.MustRegister(apiRequestDuration, transactionsProcessed, cacheHits, cacheMisses, cpuUsage, memoryUsage)

	return &Metrics{
		config:                config,
		ctx:                   ctx,
		apiRequestDuration:    apiRequestDuration,
		transactionsProcessed: transactionsProcessed,
		cacheHits:             cacheHits,
		cacheMisses:           cacheMisses,
		cpuUsage:              cpuUsage,
		memoryUsage:           memoryUsage,
	}
}

// API request duration histogram
var APIRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "api_request_duration_seconds",
		Help:    "Histogram of response time for API requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"route", "method"},
)

// Transactions processed per second
var TransactionsProcessed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "transactions_processed_total",
		Help: "Total number of transactions processed",
	},
	[]string{"status"},
)

// Cache hit/miss rates
var CacheHits = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total number of cache hits",
	},
)
var CacheMisses = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total number of cache misses",
	},
)

// Register all metrics
func RegisterMetrics() {
	prometheus.MustRegister(APIRequestDuration, TransactionsProcessed, CacheHits, CacheMisses)
}

func (m *Metrics) SetAPIRequestDuration(route, method string, duration float64) {
	m.apiRequestDuration.WithLabelValues(route, method).Observe(duration)
}

func (m *Metrics) IncrementTransactionsProcessed(status string) {
	m.transactionsProcessed.WithLabelValues(status).Inc()
}

func (m *Metrics) IncrementCacheHits() {
	m.cacheHits.Inc()
}

func (m *Metrics) IncrementCacheMisses() {
	m.cacheMisses.Inc()
}

func (m *Metrics) UpdateCPUUsage() {
	cpus, err := cpu.Percent(0, true)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return
	}
	for i, cpuPercent := range cpus {
		m.cpuUsage.WithLabelValues(strconv.Itoa(i)).Set(cpuPercent)
	}
}

// Function to observe memory usage
func (m *Metrics) UpdateMemoryUsage() {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return
	}
	m.memoryUsage.Set(float64(memStats.Used))
}

// Update CPU and memory usage at regular intervals
func (m *Metrics) StartMonitoring() {
	go func() {
		ticker := time.NewTicker(MONITORING_INTERVAL)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Update CPU and memory usage
				m.UpdateCPUUsage()
				m.UpdateMemoryUsage()
			case <-m.ctx.Done():
				return
			}
		}
	}()
}
