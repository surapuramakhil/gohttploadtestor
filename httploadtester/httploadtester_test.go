package main

import (
	"bytes"
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

func TestLoadTester(t *testing.T) {
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
			args:    []string{"-url", "http://example.com", "-qps", "2", "-duration", "2"},
			want:    "Total Requests: 100",
			wantErr: false,
		},
		{
			name:    "missing URL",
			args:    []string{"-qps", "2", "-duration", "2"},
			want:    "URL is required",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(t, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("runBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			lastLines := getLastNLines(output, 4)
			if !strings.Contains(lastLines, tt.want) {
				t.Errorf("runBinary() output = %v, want %v", lastLines, tt.want)
			}
		})
	}
}