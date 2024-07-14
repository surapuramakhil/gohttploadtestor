package config

import (
	"flag"
	"fmt"
)

type Config struct {
	URL      string
	QPS      int
	Method   string
	Payload  string
	FilePath string
	Duration int
}

func ParseFlags() *Config {
	url := flag.String("url", "", "HTTP address to load test")
	qps := flag.Int("qps", 1, "Queries per second")
	method := flag.String("method", "GET", "HTTP method to use")
	payload := flag.String("payload", "", "Payload to include in the request body")
	filePath := flag.String("file", "", "Path to the file to upload")
	duration := flag.Int("duration", 10, "Duration of the test in seconds")
	flag.Parse()

	if *url == "" {
		fmt.Println("URL is required")
		return nil
	}

	return &Config{
		URL:      *url,
		QPS:      *qps,
		Method:   *method,
		Payload:  *payload,
		FilePath: *filePath,
		Duration: *duration,
	}
}