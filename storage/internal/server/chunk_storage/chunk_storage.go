package chunkstorage

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

func (fs *ChunkStorage) Get(chunkID uuid.UUID, part uint8) ([]byte, error) {
	b, err := os.ReadFile(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	if err != nil {
		return nil, errors.Wrapf(err, "read file from disk: %s", fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	}

	return b, nil
}

func (fs *ChunkStorage) Erase(chunkID uuid.UUID, part uint8) error {
	err := os.Remove(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	if err != nil {
		return errors.Wrapf(err, "remove file from disk: %s", fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
	}

	return nil
}
