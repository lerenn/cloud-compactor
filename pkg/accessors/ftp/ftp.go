package ftp

import (
	"crypto/tls"
	"io"
	"log"
	"os"
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
	if err := c.Login(a.config.Username, a.config.Password); err != nil {
		return nil, err
	}

	return c, nil
}

func quitConnection(c *ftp.ServerConn) {
	if err := c.Quit(); err != nil {
		log.Fatal(err.Error())
	}
}

func list(conn *ftp.ServerConn, ch chan string, path string) {
	entries, err := conn.List(path)
	if err != nil {
		log.Println("Error listing files:", err)
		return
	}

	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		if entry.Type == ftp.EntryTypeFolder {
			list(conn, ch, path+"/"+entry.Name)
		} else {
			ch <- path + "/" + entry.Name
		}
	}
}

func (a Accessor) List(path string) (chan string, error) {
	conn, err := a.getConnection()
	if err != nil {
		return nil, err
	}

	ch := make(chan string, 1024*1024)
	log.Println("Start of files listing from FTP...")
	go func() {
		// Close channel when done
		defer close(ch)

		// Close connection when done
		defer quitConnection(conn)

		// Log end of listing
		defer log.Println("End of files listing from FTP.")

		// List files
		list(conn, ch, path)
	}()

	return ch, nil
}

func (a Accessor) Download(path string) (string, error) {
	// Get connection
	conn, err := a.getConnection()
	if err != nil {
		return "", err
	}
	defer quitConnection(conn)

	// Create a temp file
	file, err := os.CreateTemp("/tmp", "download")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Download the file
	resp, err := conn.Retr(path)
	if err != nil {
		return "", err
	}

	// Copy the file to the temp file
	_, err = io.Copy(file, resp)
	if err != nil {
		log.Fatal(err)
	}

	return file.Name(), nil

}

func (a Accessor) Upload(localPath, remotePath string) error {
	// Get connection
	conn, err := a.getConnection()
	if err != nil {
		return err
	}
	defer quitConnection(conn)

	// Open the local file
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file
	return conn.Stor(remotePath, file)
}

func (a Accessor) Delete(path string) error {
	// Get connection
	conn, err := a.getConnection()
	if err != nil {
		return err
	}
	defer quitConnection(conn)

	// Delete the file
	return conn.Delete(path)
}
