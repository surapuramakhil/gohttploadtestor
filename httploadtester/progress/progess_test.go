package progress

import (
	"testing"
)

func TestInitializeProgressBar(t *testing.T) {
	bar := InitializeProgressBar(10)
	if bar == nil {
		t.Errorf("expected progress bar to be initialized, got nil")
	}
}