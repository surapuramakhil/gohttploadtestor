package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"httploadtester/config"
	"httploadtester/httpclient"
	"httploadtester/metrics"
	"httploadtester/progress"
)

func main() {
	cfg := config.ParseFlags()
	if cfg == nil {
		return
	}

	metrics := &metrics.Metrics{}
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second / time.Duration(cfg.QPS))
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Duration)*time.Second)
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	bar := progress.InitializeProgressBar(cfg.Duration)

	// Start a goroutine to update the progress bar
	go func() {
		for {
			select {
			case <-ctx.Done():
				bar.Finish()
				return
			case <-time.After(1 * time.Second):
				bar.Add(1)
				metrics.PrintStatus()
			}
		}
	}()

	client := &http.Client{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			metrics.PrintFinalMetrics()
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				var req *http.Request
				var err error

				if cfg.FilePath != "" {
					req, err = httpclient.NewFileUploadRequest(cfg.URL, cfg.Method, cfg.FilePath)
				} else {
					req, err = http.NewRequest(cfg.Method, cfg.URL, bytes.NewBufferString(cfg.Payload))
				}

				if err != nil {
					metrics.RecordError()
					return
				}

				startTime := time.Now()
				resp, err := httpclient.SendRequest(client, req)
				endTime := time.Now()
				if err != nil {
					metrics.RecordError()
					return
				}
				defer resp.Body.Close()

				metrics.RecordLatency(endTime.Sub(startTime))
			}()
		}
	}
}