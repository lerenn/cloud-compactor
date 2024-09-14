package accessors

type Accessor interface {
	List(path string) ([]string, error)
}
