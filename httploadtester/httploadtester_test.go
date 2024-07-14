package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	cmd := exec.Command("go", "build", "-o", "http_load_tester")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	return "./http_load_tester"
}

func runBinary(t *testing.T, args []string) (string, error) {
	cmd := exec.Command(buildBinary(t), args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func getLastNLines(output string, n int) string {
	lines := strings.Split(output, "\n")
	if len(lines) < n {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

func extractTotalRequests(output string) (int, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Total Requests:") {
			var totalRequests int
			_, err := fmt.Sscanf(line, "Total Requests: %d", &totalRequests)
			return totalRequests, err
		}
	}
	return 0, fmt.Errorf("Total Requests not found in output")
}

func TestLoadTester(t *testing.T) {
	binary := buildBinary(t)
	defer func() {
		_ = exec.Command("rm", binary).Run()
	}()

	tests := []struct {
		name      string
		args      []string
		expected  int
		tolerance float64 // Tolerance in percentage
		wantErr   bool
	}{
		{
			name:      "valid GET request",
			args:      []string{"-url", "http://example.com", "-qps", "2", "-duration", "2"},
			expected:  4,
			tolerance: 0.0,
			wantErr:   false,
		},
		{
			name:      "surge requests",
			args:      []string{"-url", "http://example.com", "-qps", "1000", "-duration", "5"},
			expected:  5000,
			tolerance: 1.0, // 1% tolerance
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(t, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("runBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !strings.Contains(output, "URL is required") {
					t.Errorf("runBinary() output = %v, want %v", output, "URL is required")
				}
				return
			}

			lastLines := getLastNLines(output, 4)
			totalRequests, err := extractTotalRequests(lastLines)
			if err != nil {
				t.Errorf("extractTotalRequests() error = %v", err)
				return
			}

			lowerBound := float64(tt.expected) * (1 - tt.tolerance/100)
			upperBound := float64(tt.expected) * (1 + tt.tolerance/100)

			if float64(totalRequests) < lowerBound || float64(totalRequests) > upperBound {
				t.Errorf("total requests = %d, want between %.2f and %.2f", totalRequests, lowerBound, upperBound)
			}
		})
	}
}