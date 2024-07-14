package config

import (
	"flag"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Simulate command line arguments
	args := []string{"cmd", "-url", "http://example.com", "-qps", "10", "-method", "POST", "-payload", "data", "-file", "test.txt", "-duration", "20"}
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ExitOnError)

	url := flag.CommandLine.String("url", "", "HTTP address to load test")
	qps := flag.CommandLine.Int("qps", 1, "Queries per second")
	method := flag.CommandLine.String("method", "GET", "HTTP method to use")
	payload := flag.CommandLine.String("payload", "", "Payload to include in the request body")
	filePath := flag.CommandLine.String("file", "", "Path to the file to upload")
	duration := flag.CommandLine.Int("duration", 10, "Duration of the test in seconds")

	flag.CommandLine.Parse(args[1:])

	cfg := &Config{
		URL:      *url,
		QPS:      *qps,
		Method:   *method,
		Payload:  *payload,
		FilePath: *filePath,
		Duration: *duration,
	}

	if cfg.URL != "http://example.com" {
		t.Errorf("expected URL to be 'http://example.com', got %s", cfg.URL)
	}
	if cfg.QPS != 10 {
		t.Errorf("expected QPS to be 10, got %d", cfg.QPS)
	}
	if cfg.Method != "POST" {
		t.Errorf("expected Method to be 'POST', got %s", cfg.Method)
	}
	if cfg.Payload != "data" {
		t.Errorf("expected Payload to be 'data', got %s", cfg.Payload)
	}
	if cfg.FilePath != "test.txt" {
		t.Errorf("expected FilePath to be 'test.txt', got %s", cfg.FilePath)
	}
	if cfg.Duration != 20 {
		t.Errorf("expected Duration to be 20, got %d", cfg.Duration)
	}
}