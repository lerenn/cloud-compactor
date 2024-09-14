package cloudcompactor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lerenn/cloud-compactor/pkg/accessors"
	"github.com/lerenn/cloud-compactor/pkg/accessors/ftp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type CloudCompactor struct {
	config Config
}

func New(config Config) *CloudCompactor {
	return &CloudCompactor{
		config: config,
	}
}

func (c *CloudCompactor) Run() error {
	// Get correct accessor based on config
	var a accessors.Accessor
	switch {
	case c.config.FTP.Address != "":
		log.Printf("Using FTP accessor")
		a = ftp.New(c.config.FTP)
	default:
		return fmt.Errorf("no accessor found")
	}

	// List files
	files, err := c.listFiles(a)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	// Process files
	for _, f := range files {
		if err := c.processFile(a, f); err != nil {
			return fmt.Errorf("failed to process file: %w", err)
		}
	}

	return nil
}

func (c CloudCompactor) listFiles(a accessors.Accessor) ([]string, error) {
	// List files
	log.Printf("Listing files in %s...", c.config.Path)
	rawList, err := a.List(c.config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	// Filter files
	var files []string
	log.Printf("Filtering files...")
	for _, f := range rawList {
		if c.config.Formats.ProcessedSuffix != "" && strings.Contains(f, c.config.Formats.ProcessedSuffix) {
			log.Printf("Skipping processed file: %s", f)
			continue
		}

		for _, i := range c.config.Formats.Inputs {
			if strings.HasSuffix(f, i) {
				files = append(files, f)
				log.Printf("Found file: %s", f)
				break
			}
		}
	}

	return files, nil
}

func (c CloudCompactor) processFile(a accessors.Accessor, path string) error {
	// Download file
	log.Printf("Downloading file %s...", path)
	localPath, err := a.Download(path)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		log.Printf("Deleting local file %s...", localPath)
		if err := os.Remove(localPath); err != nil {
			log.Printf("Failed to delete local file %s: %s", localPath, err)
		}
	}()

	// Process file
	log.Printf("Process file %s...", path)
	localOutputPath := localPath + "." + c.config.Formats.ProcessedSuffix + "." + c.config.Formats.Output
	if err := ffmpeg.Input(localPath).
		Output(localOutputPath,
			ffmpeg.KwArgs{"c:v": "libx265"},
			ffmpeg.KwArgs{"c:a": "libfdk_aac"},
			ffmpeg.KwArgs{"b:a": "128k"},
			ffmpeg.KwArgs{"speed": c.config.Speed}).Run(); err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	// Upload file
	newPath := strings.TrimSuffix(path, filepath.Ext(path)) + "." + c.config.Formats.ProcessedSuffix + "." + c.config.Formats.Output
	log.Printf("Upload file %s...", newPath)
	if err := a.Upload(localPath, path); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	// Delete file
	log.Printf("Delete file %s...", path)
	if err := a.Delete(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
