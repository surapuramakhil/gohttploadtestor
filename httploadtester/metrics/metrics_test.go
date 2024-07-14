package metrics

import (
	"sync"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	m := &Metrics{}

	m.RecordLatency(100 * time.Millisecond)
	m.RecordLatency(200 * time.Millisecond)
	m.RecordError()

	if m.TotalRequests != 2 {
		t.Errorf("expected TotalRequests to be 2, got %d", m.TotalRequests)
	}
	if m.TotalErrors != 1 {
		t.Errorf("expected TotalErrors to be 1, got %d", m.TotalErrors)
	}
	if len(m.Latencies) != 2 {
		t.Errorf("expected Latencies length to be 2, got %d", len(m.Latencies))
	}
}

func TestMetricsConcurrency(t *testing.T) {
	m := &Metrics{}
	var wg sync.WaitGroup
	numRequests := 100

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			m.RecordLatency(50 * time.Millisecond)
			m.RecordError()
		}()
	}

	wg.Wait()

	if m.TotalRequests != numRequests {
		t.Errorf("expected TotalRequests to be %d, got %d", numRequests, m.TotalRequests)
	}
	if m.TotalErrors != numRequests {
		t.Errorf("expected TotalErrors to be %d, got %d", numRequests, m.TotalErrors)
	}
	if len(m.Latencies) != numRequests {
		t.Errorf("expected Latencies length to be %d, got %d", numRequests, len(m.Latencies))
	}
}
