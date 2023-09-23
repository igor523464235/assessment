package fs

import (
	"fmt"
	"io"
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

	go func() {
		defer close(content)
		defer close(errorChan)

		file, err := os.Open(fmt.Sprintf("%s/%s:%d", fs.directory, chunkID, part))
		if err != nil {
			errorChan <- err
			return
		}
		defer file.Close()

		// TODO: need to know optimal buffer size
		for {
			buffer := make([]byte, 1024)
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				errorChan <- err
				return
			}
			if n == 0 {
				break
			}

			content <- buffer[:n]
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
