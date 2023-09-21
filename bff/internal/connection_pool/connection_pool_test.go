package connectionpool

import (
	"testing"
)

func TestPool(t *testing.T) {
	// Creating pool
	connectionPool := New()

	// Creating connections
	conn1 := &Connection{ID: 1, FreeSpace: 100}
	conn2 := &Connection{ID: 2, FreeSpace: 250}
	conn3 := &Connection{ID: 3, FreeSpace: 150}
	conn4 := &Connection{ID: 4, FreeSpace: 150}
	conn5 := &Connection{ID: 5, FreeSpace: 160}
	conn6 := &Connection{ID: 6, FreeSpace: 150}
	conn7 := &Connection{ID: 7, FreeSpace: 150}
	conn8 := &Connection{ID: 8, FreeSpace: 190}
	conn9 := &Connection{ID: 9, FreeSpace: 150}
	conn10 := &Connection{ID: 10, FreeSpace: 50}

	// Adding connections
	connectionPool.AddConnection(conn1)
	connectionPool.AddConnection(conn2)
	connectionPool.AddConnection(conn3)
	connectionPool.AddConnection(conn4)
	connectionPool.AddConnection(conn5)
	connectionPool.AddConnection(conn6)
	connectionPool.AddConnection(conn7)
	connectionPool.AddConnection(conn8)
	connectionPool.AddConnection(conn9)
	connectionPool.AddConnection(conn10)

	// Getting connection with the biggest free space on storage
	maxConn := connectionPool.GetConnectionByFreeSpace()
	if maxConn.ID != conn2.ID {
		t.Errorf("Expected connection ID: 2, but got: %d", maxConn.ID)
	}

	// Updating size for a storage
	connectionPool.UpdateFreeSpace(2, 30)

	// Getting connection with the biggest free space on storage
	maxConn = connectionPool.GetConnectionByFreeSpace()
	if maxConn.ID != conn8.ID {
		t.Errorf("Expected connection ID: 8, but got: %d", maxConn.ID)
	}

	// Updating non-existent connection
	connectionPool.UpdateFreeSpace(5, 60)
	maxConn = connectionPool.GetConnectionByFreeSpace()
	if maxConn.ID != conn8.ID {
		t.Errorf("Expected connection ID: 8, but got: %d", maxConn.ID)
	}
}
