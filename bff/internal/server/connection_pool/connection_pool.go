package connectionpool

import (
	"container/heap"
	"str/bff/internal/server"
	csgrpc "str/storage/grpc"
	"sync"
)

// Connection is a grpc connection to file storage
// with additional data - id and free space on the storage
type Connection struct {
	ID        server.ConnectionID
	FreeSpace uint64
	csgrpc.StorageServiceClient
}

// ConnectionHeap is data structure for quickly obtaining a connection with the maximum value of FreeSpace.
type ConnectionHeap []*Connection

func (h ConnectionHeap) Len() int           { return len(h) }
func (h ConnectionHeap) Less(i, j int) bool { return h[i].FreeSpace > h[j].FreeSpace }
func (h ConnectionHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ConnectionHeap) Push(x interface{}) {
	*h = append(*h, x.(*Connection))
}

func (h *ConnectionHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

// ConnectionPool is a struct to get connections by their id (from map)
// or to get a connection from max-heap.
type ConnectionPool struct {
	mu          sync.Mutex
	connections map[server.ConnectionID]*Connection
	maxHeap     ConnectionHeap
}

// IterateConnections is just iteration through all stored connections.
func (cp *ConnectionPool) IterateConnections() <-chan func() (conn *Connection) {
	c := make(chan func() *Connection, len(cp.maxHeap))
	for _, conn := range cp.connections {
		connection := conn
		c <- func() *Connection {
			return connection
		}
	}
	close(c)

	return c
}

// New creates new empty connection pool.
func New() *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[server.ConnectionID]*Connection),
		maxHeap:     ConnectionHeap{},
	}
}

// AddConnection adds the connection to the connection pool.
func (cp *ConnectionPool) AddConnection(connection *Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.connections[connection.ID] = connection
	heap.Push(&cp.maxHeap, connection)
}

// UpdateFreeSpace updates free space value for the connection and rebalances the heap.
func (cp *ConnectionPool) UpdateFreeSpace(connectionID server.ConnectionID, size uint64) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conn, exists := cp.connections[connectionID]; exists {
		conn.FreeSpace = size
		heap.Init(&cp.maxHeap)
	}
}

// GetConnectionByFreeSpace returns a connection with the biggest free space value.
func (cp *ConnectionPool) GetConnectionByFreeSpace() *Connection {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if len(cp.maxHeap) == 0 {
		return nil
	}

	return cp.maxHeap[0]
}

// GetConnection returns connection by it's ID.
func (cp *ConnectionPool) GetConnection(connectionID server.ConnectionID) *Connection {
	return cp.connections[connectionID]
}
