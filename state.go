package main

// DownloadState holds information about the current state of a download.
type DownloadState struct {
	CurrentBytes int64 // Bytes downloaded so far
	TotalBytes   int64 // Total expected bytes, 0 if unknown
}
