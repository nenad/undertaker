package storage

// Gravedigger provides methods for inserting and fetching functions from a storage
type Gravedigger interface {
	// Bury stores the given functions
	Bury(functions []string) error
	// Dig returns the functions which were never invoked
	Dig() (functions []string, err error)
}
