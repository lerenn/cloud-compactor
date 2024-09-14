package ftp

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

type Accessor struct {
	config Config
}

func New(c Config) *Accessor {
	return &Accessor{
		config: c,
	}
}

func (a Accessor) getConnection() (*ftp.ServerConn, error) {
	// Create a TLS configuration.
	tlsConfig := &tls.Config{
		// Enable TLS 1.2.
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	// Connect to the FTP server.
	c, err := ftp.Dial(a.config.Address, ftp.DialWithTLS(tlsConfig), ftp.DialWithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}

	// Login to the FTP server.
	if err := c.Login(a.config.User, a.config.Password); err != nil {
		return nil, err
	}

	return c, nil
}

func quitConnection(c *ftp.ServerConn) {
	if err := c.Quit(); err != nil {
		log.Fatal(err.Error())
	}
}

func list(conn *ftp.ServerConn, path string) ([]string, error) {
	entries, err := conn.List(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		if entry.Type == ftp.EntryTypeFolder {
			log.Println("Exploring", path+"/"+entry.Name)
			subfiles, err := list(conn, path+"/"+entry.Name)
			if err != nil {
				return nil, err
			}

			files = append(files, subfiles...)
		} else {
			files = append(files, entry.Name)
		}
	}

	return files, nil
}

func (a Accessor) List(path string) ([]string, error) {
	conn, err := a.getConnection()
	if err != nil {
		return nil, err
	}
	defer quitConnection(conn)

	return list(conn, path)
}
