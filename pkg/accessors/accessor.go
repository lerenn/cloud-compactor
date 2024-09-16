package accessors

type Accessor interface {
	List(path string) (chan string, error)
	Download(path string) (string, error)
	Upload(localPath, remotePath string) error
	Delete(path string) error
}
