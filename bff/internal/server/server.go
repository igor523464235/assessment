package server

import (
	"sort"

	"github.com/google/uuid"
)

type FileMetadata struct {
	ID       uuid.UUID
	Name     string
	Size     uint64
	Storages Storages
}

type FilePart uint32
type ConnectionID uint8

type Storages map[FilePart]ConnectionID

type Storage struct {
	Part         FilePart
	ConnectionID ConnectionID
}

// SortByParts sorts connection ids according to file parts.
func (ss Storages) SortByParts() []*Storage {
	keys := make([]int, 0, len(ss))
	for k := range ss {
		keys = append(keys, int(k))
	}

	// asc
	sort.Ints(keys)

	storages := make([]*Storage, len(ss))

	for i, k := range keys {
		storages[i] = &Storage{
			Part:         FilePart(k),
			ConnectionID: ss[FilePart(k)],
		}
	}

	return storages
}
