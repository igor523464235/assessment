package fs

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ChunkStorage struct {
	directory string
}

func New(directory string) *ChunkStorage {
	return &ChunkStorage{directory}
}

func (fs *ChunkStorage) Save(content <-chan []byte, chunkID uuid.UUID, part uint8) error {
	file, err := os.Create(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	if err != nil {
		return errors.Wrap(err, "create chunk file on disk")
	}

	defer func() error {
		if err := file.Close(); err != nil {
			return errors.Wrap(err, "close chunk file after saving on disk")
		}

		return nil
	}()

	for v := range content {
		if _, err := file.Write(v); err != nil {
			return errors.Wrap(err, "write chunk file content to disk")
		}
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}

func (fs *ChunkStorage) Get(chunkID uuid.UUID, part uint8) (<-chan []byte, <-chan error) {
	content := make(chan []byte)
	errorChan := make(chan error)

	// TODO: (important) make stream reading file, not full in memory
	go func() {
		defer close(content)
		defer close(errorChan)

		b, err := os.ReadFile(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
		if err != nil {
			errorChan <- errors.Wrapf(err, "read file from disk: %s", fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
			return
		}

		for i := 0; i < len(b); i += 1024 {
			end := i + 1024
			if end > len(b) {
				end = len(b)
			}
			chunk := b[i:end]
			content <- chunk
		}
	}()

	return content, errorChan
}

func (fs *ChunkStorage) Erase(chunkID uuid.UUID, part uint8) error {
	err := os.Remove(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	if err != nil {
		return errors.Wrapf(err, "remove file from disk: %s", fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	}

	return nil
}
