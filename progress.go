package main

import (
	"fmt"
	"time"
)

// ProgressBar struct to track and display download progress.
type ProgressBar struct {
	total   int64
	current int64
	start   time.Time
}

// Write implements the io.Writer interface for ProgressBar.
// This allows it to slot into io.TeeReader to track progress.
func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.current += int64(n)
	pb.printProgress()

	// Return the number of bytes "written", as we are just tracking progress.
	// This is needed to satisfy the io.Writer interface, but doesn't serve
	// a real purpose otherwise.
	return n, nil
}

// Uses the current state to print the progress to stdout.
func (pb *ProgressBar) printProgress() {
	if pb.total == 0 {
		fmt.Printf("\rDownloading... %s", humanReadableByteCount(pb.current))
		return
	}

	percent := float64(pb.current) / float64(pb.total) * 100

	fmt.Printf("\rDownloading: %.2f%% %s/%s",
		percent,
		humanReadableByteCount(pb.current),
		humanReadableByteCount(pb.total),
	)
}

// humanReadableByteCount converts a raw number of bytes to a human-readable format.
// e.g "1.5 MB" instead of "1500000 B"
func humanReadableByteCount(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, index := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		index++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[index])
}
