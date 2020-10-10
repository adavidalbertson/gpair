package store

// InMemoryStore is an implementation of Store that simply holds bytes in memory
// For testing purposes only
type InMemoryStore struct {
	bytes []byte
}

func (ims *InMemoryStore) Read() ([]byte, error) {
	return ims.bytes, nil
}

func (ims *InMemoryStore) Write(bytes []byte) error {
	ims.bytes = bytes
	return nil
}

func (ims *InMemoryStore) GetPath() string {
	return "This is an in-memory store not backed by a file on disk."
}