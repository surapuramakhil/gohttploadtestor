// +build integration

package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"net/http"
	"net/http/httptest"
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

func createTestServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the request
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	return server
}

func TestLoadTester(t *testing.T) {

	server := createTestServer()
	defer server.Close()

	binary := buildBinary(t)
	defer func() {
		_ = exec.Command("rm", binary).Run()
	}()

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "valid GET request",
			args:    []string{"-url", server.URL, "-qps", "45", "-duration", "10"},
			want:    "Total Requests: 450",
		},
		{
			name:    "missing URL",
			args:    []string{"-qps", "2", "-duration", "2"},
			want:    "URL is required",
		},
		{
			name:	"surge requests",
			args:	[]string{"-url", server.URL , "-qps", "1000", "-duration", "4"},
			want:	"Total Requests: 4000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := runBinary(t, tt.args)

			lastLines := getLastNLines(output, 4)
			if !strings.Contains(lastLines, tt.want) {
				t.Errorf("runBinary() output = %v, want %v", lastLines, tt.want)
			}
		})
	}
}