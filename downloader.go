package main

import (
	"fmt"
	"os"
	"os/exec"
)

const downloadDir = "downloads"

// downloadVideo downloads the video at url via yt-dlp into outPath.
func downloadVideo(url, outPath string) error {
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return fmt.Errorf("could not create download dir: %w", err)
	}

	cmd := exec.Command("yt-dlp",
		"-o", outPath,
		"--merge-output-format", "mp4",
		"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
		url,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("yt-dlp failed: %w: %s", err, string(out))
	}

	return nil
}

// clearDownloadDir wipes any leftover files from a previous run. Only safe
// to call at startup, before any jobs are in flight.
func clearDownloadDir() error {
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		_ = os.Remove(downloadDir + string(os.PathSeparator) + e.Name())
	}
	return nil
}
