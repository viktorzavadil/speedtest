package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(d)/float64(time.Millisecond))
}

func FormatBytes(b int64) string {
	switch {
	case b >= 1e12:
		return fmt.Sprintf("%.1fTB", float64(b)/1e12)
	case b >= 1e9:
		return fmt.Sprintf("%.1fGB", float64(b)/1e9)
	case b >= 1e6:
		return fmt.Sprintf("%.1fMB", float64(b)/1e6)
	case b >= 1e3:
		return fmt.Sprintf("%.1fkB", float64(b)/1e3)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func FormatBitrate(bytes int64, d time.Duration) string {
	var bps float64
	if d > 0 {
		bps = float64(bytes) * 8 / d.Seconds()
	} else {
		bps = 0
	}

	switch {
	case bps >= 1e12:
		return fmt.Sprintf("%.1fTbit/s", bps/1e12)
	case bps >= 1e9:
		return fmt.Sprintf("%.1fGbit/s", bps/1e9)
	case bps >= 1e6:
		return fmt.Sprintf("%.1fMbit/s", bps/1e6)
	case bps >= 1e3:
		return fmt.Sprintf("%.1fkbit/s", bps/1e3)
	default:
		return fmt.Sprintf("%.0fbit/s", bps)
	}
}
