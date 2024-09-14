package cloudcompactor

import (
	"fmt"
	"log"
	"strings"

	"github.com/lerenn/cloud-compactor/pkg/accessors"
	"github.com/lerenn/cloud-compactor/pkg/accessors/ftp"
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
	log.Printf("Listing files in %s...", c.config.Path)
	rawList, err := a.List(c.config.Path)
	if err != nil {
		return fmt.Errorf("failed to list: %w", err)
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

	return nil
}
