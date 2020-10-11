package store

// InMemoryStore is an implementation of Store that simply holds bytes in memory
// For testing purposes only
type InMemoryStore struct {
	bytes []byte
}

// Read returns the bytes in the store
func (ims *InMemoryStore) Read() ([]byte, error) {
	return ims.bytes, nil
}

// Write stores bytes in the store
func (ims *InMemoryStore) Write(bytes []byte) error {
	ims.bytes = bytes
	return nil
}

// GetPath is just here to fulfil the interface contract
func (ims *InMemoryStore) GetPath() string {
	return "This is an in-memory store not backed by a file on disk."
}