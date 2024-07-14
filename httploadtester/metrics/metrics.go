package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests int
	TotalErrors   int
	Latencies     []time.Duration
	mu            sync.Mutex
}

func (m *Metrics) RecordLatency(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Latencies = append(m.Latencies, latency)
	m.TotalRequests++
}

func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalErrors++
}

func (m *Metrics) PrintStatus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalLatency := time.Duration(0)
	for _, latency := range m.Latencies {
		totalLatency += latency
	}
	averageLatency := time.Duration(0)
	if len(m.Latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(m.Latencies))
	}

	fmt.Printf("Total Requests: %d | Total Errors: %d | Average Latency: %s\n", m.TotalRequests, m.TotalErrors, averageLatency)
}

func (m *Metrics) PrintFinalMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalLatency := time.Duration(0)
	for _, latency := range m.Latencies {
		totalLatency += latency
	}
	averageLatency := time.Duration(0)
	if len(m.Latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(m.Latencies))
	}

	fmt.Printf("\nFinal Results:\n")
	fmt.Printf("Total Requests: %d\n", m.TotalRequests)
	fmt.Printf("Total Errors: %d\n", m.TotalErrors)
	fmt.Printf("Average Latency: %s\n", averageLatency)
}