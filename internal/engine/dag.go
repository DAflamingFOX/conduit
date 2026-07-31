package engine

import (
	"fmt"

	"github.com/google/uuid"
)

// NodeInstance represents a node placed on the visual DAG canvas.
type NodeInstance struct {
	ID             string                 `json:"id"`               // Instance UUID v7
	NodeManifestID string                 `json:"node_manifest_id"` // e.g. "conduit.file.move"
	Name           string                 `json:"name"`
	Parameters     map[string]interface{} `json:"parameters"`
}

// Connection represents an edge between an output port of a node and an input port of a downstream node.
type Connection struct {
	ID         string `json:"id"` // Connection UUID v7
	FromNodeID string `json:"from_node_id"`
	FromOutput string `json:"from_output"` // e.g. "success", "error"
	ToNodeID   string `json:"to_node_id"`
	ToInput    string `json:"to_input"`
}

// FlowGraph represents a complete visual pipeline graph.
type FlowGraph struct {
	ID          string         `json:"id"` // Flow UUID v7
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Nodes       []NodeInstance `json:"nodes"`
	Connections []Connection   `json:"connections"`
}

// NewFlowGraph initializes a new FlowGraph with a UUID v7 ID.
func NewFlowGraph(name string, description string) (*FlowGraph, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID v7 for FlowGraph: %w", err)
	}
	return &FlowGraph{
		ID:          u.String(),
		Name:        name,
		Description: description,
		Nodes:       []NodeInstance{},
		Connections: []Connection{},
	}, nil
}

// Validate checks the FlowGraph for cycles and invalid connection references.
func (g *FlowGraph) Validate() error {
	nodeMap := make(map[string]bool)
	for _, n := range g.Nodes {
		if n.ID == "" {
			return fmt.Errorf("found node instance without an ID")
		}
		if nodeMap[n.ID] {
			return fmt.Errorf("duplicate node instance ID: %s", n.ID)
		}
		nodeMap[n.ID] = true
	}

	// Build adjacency list for cycle detection
	adj := make(map[string][]string)
	for _, conn := range g.Connections {
		if !nodeMap[conn.FromNodeID] {
			return fmt.Errorf("connection references non-existent source node: %s", conn.FromNodeID)
		}
		if !nodeMap[conn.ToNodeID] {
			return fmt.Errorf("connection references non-existent target node: %s", conn.ToNodeID)
		}
		adj[conn.FromNodeID] = append(adj[conn.FromNodeID], conn.ToNodeID)
	}

	// DFS Cycle Detection (0=Unvisited, 1=Visiting, 2=Visited)
	state := make(map[string]int)
	var hasCycle func(nodeID string) bool
	hasCycle = func(nodeID string) bool {
		state[nodeID] = 1 // Mark visiting
		for _, neighbor := range adj[nodeID] {
			if state[neighbor] == 1 {
				return true // Cycle detected!
			}
			if state[neighbor] == 0 {
				if hasCycle(neighbor) {
					return true
				}
			}
		}
		state[nodeID] = 2 // Mark visited
		return false
	}

	for _, n := range g.Nodes {
		if state[n.ID] == 0 {
			if hasCycle(n.ID) {
				return fmt.Errorf("cycle detected in flow graph '%s' containing node '%s'", g.Name, n.ID)
			}
		}
	}

	return nil
}

// GetNextNodes finds downstream node IDs connected to a source node's specific output branch.
func (g *FlowGraph) GetNextNodes(fromNodeID string, outputBranch string) []string {
	var next []string
	for _, conn := range g.Connections {
		if conn.FromNodeID == fromNodeID && conn.FromOutput == outputBranch {
			next = append(next, conn.ToNodeID)
		}
	}
	return next
}
