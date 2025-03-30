package daemon

import (
	"fmt"
	"net"
	"sync"
)

// TCPConnection holds both service and tunnel TCP connections
type TCPConnection struct {
	ServiceConn net.Conn // Connection to the local service
	TunnelConn  net.Conn // Connection to the remote tunnel endpoint
	Done        chan struct{}
}

// TunnelConnection represents all connections for a single tunnel
type TunnelConnection struct {
	TCP TCPConnection
	// TODO: UDP and HTTP connections will be added here
}

// ConnectionManager handles all active tunnel connections
type ConnectionManager struct {
	connections map[string]*TunnelConnection
	mu          sync.RWMutex // mutex for concurrent access to the map
}

// NewConnectionManager creates a new instance of ConnectionManager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*TunnelConnection),
	}
}

// AddConnection adds a new tunnel connection to the manager
func (cm *ConnectionManager) AddConnection(id string, tunnelConn, serviceConn net.Conn) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.connections[id]; exists {
		return fmt.Errorf("connection with id %s already exists", id)
	}

	cm.connections[id] = &TunnelConnection{
		TCP: TCPConnection{
			ServiceConn: serviceConn,
			TunnelConn:  tunnelConn,
			Done:        make(chan struct{}),
		},
	}

	return nil
}

// GetConnection retrieves a tunnel connection by its ID
func (cm *ConnectionManager) GetConnection(id string) (*TunnelConnection, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, exists := cm.connections[id]
	if !exists {
		return nil, fmt.Errorf("connection with id %s not found", id)
	}

	return conn, nil
}

// CloseConnection closes all connections for a specific tunnel
func (cm *ConnectionManager) CloseConnection(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[id]
	if !exists {
		return fmt.Errorf("connection with id %s not found", id)
	}

	// Close TCP connections
	if err := conn.TCP.TunnelConn.Close(); err != nil {
		return fmt.Errorf("error closing tunnel connection: %v", err)
	}
	if err := conn.TCP.ServiceConn.Close(); err != nil {
		return fmt.Errorf("error closing service connection: %v", err)
	}

	// Close the done channel
	close(conn.TCP.Done)

	// Remove the connection from the map
	delete(cm.connections, id)

	return nil
}

// CloseAll closes all active connections
func (cm *ConnectionManager) CloseAll() []error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var errors []error
	for id := range cm.connections {
		if err := cm.CloseConnection(id); err != nil {
			errors = append(errors, fmt.Errorf("error closing connection %s: %v", id, err))
		}
	}

	return errors
}

// ListConnections returns all active connection IDs
func (cm *ConnectionManager) ListConnections() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	connections := make([]string, 0, len(cm.connections))
	for id := range cm.connections {
		connections = append(connections, id)
	}

	return connections
}
