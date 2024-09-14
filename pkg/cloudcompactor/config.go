package cloudcompactor

import "github.com/lerenn/cloud-compactor/pkg/accessors/ftp"

type Config struct {
	FTP ftp.Config

	Path string
}
