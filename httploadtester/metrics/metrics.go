package metrics

import (
	"fmt"
	"sync"
	"time"
	"sort"
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

func (m *Metrics) calculateAverageLatency() time.Duration {
	totalLatency := time.Duration(0)
	for _, latency := range m.Latencies {
		totalLatency += latency
	}
	averageLatency := time.Duration(0)
	if len(m.Latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(m.Latencies))
	}
	return averageLatency
}

func (m *Metrics) calculatePercentile(p float64) time.Duration {
	
	// Sort the latencies in ascending order
	sortedLatencies := make([]time.Duration, len(m.Latencies))
	copy(sortedLatencies, m.Latencies)
	sort.Slice(sortedLatencies, func(i, j int) bool {
		return sortedLatencies[i] < sortedLatencies[j]
	})

	// Calculate the index of the desired percentile
	index := int(p * float64(len(sortedLatencies)-1))

	return sortedLatencies[index]
}

func (m *Metrics) PrintStatus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	averageLatency := m.calculateAverageLatency()

	fmt.Printf("Total Requests: %d | Total Errors: %d | Average Latency: %s\n", m.TotalRequests, m.TotalErrors, averageLatency)
}

func (m *Metrics) PrintFinalMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	averageLatency := m.calculateAverageLatency()
	percentiles := []float64{0.5, 0.75, 0.9, 0.95, 0.99}

	fmt.Printf("\nFinal Results:\n")
	fmt.Printf("Total Requests: %d\n", m.TotalRequests)
	fmt.Printf("Total Errors: %d\n", m.TotalErrors)
	fmt.Printf("Average Latency: %s\n", averageLatency)

	fmt.Printf("Latency Percentiles:\n")
	for _, p := range percentiles {
		percentile := m.calculatePercentile(p)
		fmt.Printf("Percentile %.2f: %s\n", p*100, percentile)
	}
}