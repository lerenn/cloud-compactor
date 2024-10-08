package cloudcompactor

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	// List filesChan
	filesChan, err := c.listFiles(a)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	// Start daemons
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	delete := c.deleteFileDaemon(ctx, a, &wg)
	upload := c.uploadFileDaemon(ctx, a, &wg, delete)
	process := c.processFileDaemon(ctx, &wg, upload)
	download := c.downloadFileDaemeon(ctx, a, &wg, process)

	// Process files
	for f := range filesChan {
		wg.Add(1)
		download <- payload{
			remoteInputPath: f,
		}
	}

	wg.Wait()
	return nil
}

func (c CloudCompactor) listFiles(a accessors.Accessor) (chan string, error) {
	// List files
	log.Printf("Listing files in %s...", c.config.Path)
	rawFilesChan, err := a.List(c.config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	// Filter files
	filteredFilesChan := make(chan string, 1024*1024)
	log.Printf("Start of files filtering...")
	go func() {
		// Close channel when done
		defer close(filteredFilesChan)

		// Log end of files filtering
		defer log.Println("End of files filtering.")

		// Filter files
		for f := range rawFilesChan {
			if c.config.Formats.ProcessedSuffix != "" && strings.Contains(f, c.config.Formats.ProcessedSuffix) {
				continue
			}

			for _, i := range c.config.Formats.Inputs {
				if strings.HasSuffix(f, i) {
					filteredFilesChan <- f
					break
				}
			}
		}
	}()

	return filteredFilesChan, nil
}

type payload struct {
	remoteInputPath  string
	localInputPath   string
	localOutputPath  string
	remoteOutputPath string
}

func (c CloudCompactor) downloadFileDaemeon(
	ctx context.Context,
	a accessors.Accessor,
	wg *sync.WaitGroup,
	process chan payload,
) (download chan payload) {
	download = make(chan payload, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case payload := <-download:
				// Download file
				log.Printf("Downloading file %s...", payload.remoteInputPath)
				localPath, err := a.Download(payload.remoteInputPath)
				if err != nil {
					log.Printf("Failed to download file: %s", err)
					wg.Done()
					return
				}

				// Send to process
				payload.localInputPath = localPath
				process <- payload
			}
		}
	}()

	return download
}

func (c CloudCompactor) processFileDaemon(
	ctx context.Context,
	wg *sync.WaitGroup,
	upload chan payload,
) (process chan payload) {
	process = make(chan payload, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case payload := <-process:
				// Process file
				log.Printf("Process file %s...", payload.localInputPath)
				localOutputPath := payload.localInputPath + "." + c.config.Formats.ProcessedSuffix + "." + c.config.Formats.Output
				err := ffmpeg.Input(payload.localInputPath).
					Output(localOutputPath,
						ffmpeg.KwArgs{"c:v": "libx265"},
						ffmpeg.KwArgs{"c:a": "libfdk_aac"},
						ffmpeg.KwArgs{"b:a": "128k"},
						ffmpeg.KwArgs{"speed": c.config.Speed}).
					Run()
				if err != nil {
					log.Printf("Failed to process file: %s", err)
					wg.Done()
					return
				}

				// Send to upload
				payload.localOutputPath = localOutputPath
				upload <- payload

				// Deleting input local file
				log.Printf("Deleting input local file %s...", payload.localInputPath)
				if err := os.Remove(payload.localInputPath); err != nil {
					log.Printf("Failed to delete input local file %s: %s", payload.localInputPath, err)
					return
				}
			}
		}
	}()

	return process
}

func (c CloudCompactor) uploadFileDaemon(
	ctx context.Context,
	a accessors.Accessor,
	wg *sync.WaitGroup,
	delete chan payload,
) (upload chan payload) {
	upload = make(chan payload, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case payload := <-upload:
				// Upload file
				newPath := strings.TrimSuffix(payload.remoteInputPath, filepath.Ext(payload.remoteInputPath)) +
					"." + c.config.Formats.ProcessedSuffix +
					"." + c.config.Formats.Output
				log.Printf("Upload file %s...", newPath)
				if err := a.Upload(payload.localOutputPath, newPath); err != nil {
					log.Printf("Failed to upload file: %s", err)
					wg.Done()
					return
				}

				// Send to delete
				payload.remoteOutputPath = payload.localOutputPath
				delete <- payload

				// Deleting local file
				log.Printf("Deleting output local file %s...", payload.localOutputPath)
				if err := os.Remove(payload.localOutputPath); err != nil {
					log.Printf("Failed to delete output local file %s: %s", payload.localOutputPath, err)
					return
				}
			}
		}
	}()

	return upload
}

func (c CloudCompactor) deleteFileDaemon(
	ctx context.Context,
	a accessors.Accessor,
	wg *sync.WaitGroup,
) (delete chan payload) {
	delete = make(chan payload, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case payload := <-delete:
				// Delete file
				log.Printf("Delete input remote file %s...", payload.remoteInputPath)
				if err := a.Delete(payload.remoteInputPath); err != nil {
					log.Printf("Failed to delete input remote file: %s", err)
				}

				// Deplete wait group
				wg.Done()
			}
		}
	}()

	return delete
}
