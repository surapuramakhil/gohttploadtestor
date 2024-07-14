package progress

import (
	"github.com/schollz/progressbar/v3"
)

func InitializeProgressBar(duration int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(duration,
		progressbar.OptionSetDescription("\rRunning load test..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}