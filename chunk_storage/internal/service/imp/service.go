package imp

import (
	"os"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage}
}

func (s *Service) Save(content <-chan []byte, chunkID uuid.UUID, part uint8) error {
	return s.storage.Save(content, chunkID, part)
}

func (s *Service) Get(chunkID uuid.UUID, part uint8) (<-chan []byte, <-chan error) {
	return s.storage.Get(chunkID, part)
}

// TODO: implement
func (s *Service) Erase(chunkID uuid.UUID, part uint8) error {
	return nil
}

func (s *Service) GetFreeSpace() (uint64, error) {
	var stat unix.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	err = unix.Statfs(wd, &stat)
	if err != nil {
		return 0, err
	}

	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize), nil
}
