package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Printf("Downloading file from URL: %s\n", argsWithoutProg[0])

	filename := extractFilenameFromURL(argsWithoutProg[0])
	err := downloadFile(argsWithoutProg[0], filename)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
}

func downloadFile(url string, filepath string) error {
	state := DownloadState{}

	var isResuming bool
	// Can ignore error here - if file doesn't exist, we start fresh.
	if info, err := os.Stat(filepath); err == nil {
		state.CurrentBytes = info.Size()
		isResuming = true
		fmt.Printf("Resuming download for %s...\n", filepath)
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file for appending:", err)
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	if isResuming && state.CurrentBytes > 0 {
		req.Header.Set("Range", "bytes="+strconv.FormatInt(state.CurrentBytes, 10)+"-")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err
	}
	defer resp.Body.Close()

	// Maybe overly naive, but as we calculate these values,
	// if we get a 416 then we assume that we've downloaded the whole file.
	// This is because it should only be an invalid range if it's requesting
	// a range starting from the end of the file, meaning we've gotten everything.
	if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		fmt.Println("File already fully downloaded.")
		return nil
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent) {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	totalStr := resp.Header.Get("Content-Length")
	if totalStr != "" {
		state.TotalBytes, err = strconv.ParseInt(totalStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not parse Content-Length: %v\n", err)
		}
		if isResuming {
			state.TotalBytes += state.CurrentBytes
		}
	}

	// Check for Content-Range header to get total size if resuming
	// and Content-Length is not the full size.
	if rangeHeader := resp.Header.Get("Content-Range"); rangeHeader != "" {
		parts := strings.Split(rangeHeader, "/")
		if len(parts) == 2 {
			totalBytesStr := parts[1]
			if totalBytes, err := strconv.ParseInt(totalBytesStr, 10, 64); err == nil {
				state.TotalBytes = totalBytes
			}
		}
	}

	progressBar := &ProgressBar{
		total:   state.TotalBytes,
		current: state.CurrentBytes,
		start:   time.Now(),
	}

	// Use io.TeeReader to write to progress bar and file simultaneously.
	reader := io.TeeReader(resp.Body, progressBar)

	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	// Print newline after progress bar is done
	// to reset terminal line.
	fmt.Print("\n")

	return nil
}

// extractFilenameFromURL extracts the filename from a given URL - i.e,
// given a URL like "http://example.com/path/to/file.m3u8", it returns "file.m3u8".
func extractFilenameFromURL(url string) string {
	idx := strings.LastIndex(url, "/")
	filename := url[idx+1:]
	// This may include query string parts ("file.m3u8?token=..."), so strip if needed:
	if q := strings.Index(filename, "?"); q != -1 {
		filename = filename[:q]
	}
	return filename
}
