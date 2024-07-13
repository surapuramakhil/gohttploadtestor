package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Metrics struct {
	TotalRequests int
	TotalErrors   int
	Latencies     []time.Duration
	mu            sync.Mutex
}

func main() {
	url := flag.String("url", "", "HTTP address to load test")
	qps := flag.Int("qps", 1, "Queries per second")
	method := flag.String("method", "GET", "HTTP method to use")
	payload := flag.String("payload", "", "Payload to include in the request body")
	filePath := flag.String("file", "", "Path to the file to upload")
	duration := flag.Int("duration", 10, "Duration of the test in seconds")
	flag.Parse()

	if *url == "" {
		fmt.Println("URL is required")
		return
	}

	metrics := &Metrics{}
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second / time.Duration(*qps))
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*duration)*time.Second)
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Initialize the progress bar
	bar := progressbar.NewOptions(*duration,
		progressbar.OptionSetDescription("\rRunning load test..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Start a goroutine to update the progress bar
	go func() {
		for {
			select {
			case <-ctx.Done():
				bar.Finish()
				return
			case <-time.After(1 * time.Second):
				bar.Add(1)
				printStatus(metrics)
			}
		}
	}()

	// Start a goroutine to print the status periodically
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			printMetrics(metrics)
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				start := time.Now()
				var req *http.Request
				var err error

				if *filePath != "" {
					req, err = newFileUploadRequest(*url, *method, *filePath)
				} else {
					req, err = http.NewRequest(*method, *url, bytes.NewBufferString(*payload))
				}

				if err != nil {
					metrics.mu.Lock()
					metrics.TotalErrors++
					metrics.mu.Unlock()
					return
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode >= 400 {
					metrics.mu.Lock()
					metrics.TotalErrors++
					metrics.mu.Unlock()
					return
				}
				metrics.mu.Lock()
				metrics.Latencies = append(metrics.Latencies, time.Since(start))
				metrics.TotalRequests++
				metrics.mu.Unlock()
			}()
		}
	}
}

func newFileUploadRequest(uri, method, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func printStatus(metrics *Metrics) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	totalLatency := time.Duration(0)
	for _, latency := range metrics.Latencies {
		totalLatency += latency
	}
	averageLatency := time.Duration(0)
	if len(metrics.Latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(metrics.Latencies))
	}

	fmt.Printf("Total Requests: %d | Total Errors: %d | Average Latency: %s", metrics.TotalRequests, metrics.TotalErrors, averageLatency)
}

func printMetrics(metrics *Metrics) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	totalLatency := time.Duration(0)
	for _, latency := range metrics.Latencies {
		totalLatency += latency
	}
	averageLatency := time.Duration(0)
	if len(metrics.Latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(metrics.Latencies))
	}

	fmt.Printf("\nFinal Results:\n")
	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Total Errors: %d\n", metrics.TotalErrors)
	fmt.Printf("Average Latency: %s\n", averageLatency)
}