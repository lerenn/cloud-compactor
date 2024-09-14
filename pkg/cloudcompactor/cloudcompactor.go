package cloudcompactor

import (
	"fmt"

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
	var a accessors.Accessor
	switch {
	case c.config.FTP.Address != "":
		a = ftp.New(c.config.FTP)
	default:
		return fmt.Errorf("no accessor found")
	}

	fmt.Println(a.List(c.config.Path))

	return nil
}
