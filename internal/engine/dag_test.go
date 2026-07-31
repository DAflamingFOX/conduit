package engine_test

import (
	"testing"

	"github.com/daflamingfox/conduit/internal/engine"
	"github.com/daflamingfox/conduit/internal/nodes"
	"github.com/google/uuid"
)

func TestDAGValidationAndCycleDetection(t *testing.T) {
	node1UUID, _ := uuid.NewV7()
	node2UUID, _ := uuid.NewV7()
	node3UUID, _ := uuid.NewV7()

	// Valid Linear DAG: N1 -> N2 -> N3
	validGraph := &engine.FlowGraph{
		ID:   "flow-valid",
		Name: "Valid Flow",
		Nodes: []engine.NodeInstance{
			{ID: node1UUID.String(), NodeManifestID: "conduit.file.copy", Name: "Copy 1"},
			{ID: node2UUID.String(), NodeManifestID: "conduit.file.move", Name: "Move 1"},
			{ID: node3UUID.String(), NodeManifestID: "conduit.file.delete", Name: "Delete 1"},
		},
		Connections: []engine.Connection{
			{ID: "c1", FromNodeID: node1UUID.String(), FromOutput: "success", ToNodeID: node2UUID.String(), ToInput: "file"},
			{ID: "c2", FromNodeID: node2UUID.String(), FromOutput: "success", ToNodeID: node3UUID.String(), ToInput: "file"},
		},
	}

	if err := validGraph.Validate(); err != nil {
		t.Fatalf("expected valid DAG graph, got error: %v", err)
	}

	// Cyclic Graph: N1 -> N2 -> N3 -> N1
	cyclicGraph := &engine.FlowGraph{
		ID:   "flow-cyclic",
		Name: "Cyclic Flow",
		Nodes: []engine.NodeInstance{
			{ID: node1UUID.String(), NodeManifestID: "conduit.file.copy"},
			{ID: node2UUID.String(), NodeManifestID: "conduit.file.move"},
			{ID: node3UUID.String(), NodeManifestID: "conduit.file.delete"},
		},
		Connections: []engine.Connection{
			{ID: "c1", FromNodeID: node1UUID.String(), FromOutput: "success", ToNodeID: node2UUID.String(), ToInput: "file"},
			{ID: "c2", FromNodeID: node2UUID.String(), FromOutput: "success", ToNodeID: node3UUID.String(), ToInput: "file"},
			{ID: "c3", FromNodeID: node3UUID.String(), FromOutput: "success", ToNodeID: node1UUID.String(), ToInput: "file"}, // Loop!
		},
	}

	if err := cyclicGraph.Validate(); err == nil {
		t.Fatalf("expected cycle detection error for cyclic graph, but got nil")
	}
}

func TestEvaluatorPrepareStep(t *testing.T) {
	registry := nodes.NewRegistry()
	moveManifest := &nodes.Manifest{
		ID:   "conduit.file.move",
		Name: "Move File",
		Execution: nodes.ExecutionSpec{
			Type:   nodes.ExecTypeCommand,
			Binary: "mv",
			Args:   []string{"{{ inputs.file.path }}", "{{ parameters.destination }}"},
		},
	}
	registry.Register(moveManifest)

	node1UUID, _ := uuid.NewV7()
	graph := &engine.FlowGraph{
		ID:   "flow-test",
		Name: "Test Flow",
		Nodes: []engine.NodeInstance{
			{
				ID:             node1UUID.String(),
				NodeManifestID: "conduit.file.move",
				Parameters: map[string]interface{}{
					"destination": "/storage/movies",
				},
			},
		},
	}

	eval, err := engine.NewEvaluator(graph, registry)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}

	ctx, err := engine.NewContext("flow-test", engine.NewFileInfo("/tmp/input.mp4", 100), "/tmp/scratch")
	if err != nil {
		t.Fatalf("failed to create context: %v", err)
	}

	step, err := eval.PrepareStep(node1UUID.String(), ctx)
	if err != nil {
		t.Fatalf("failed to prepare step: %v", err)
	}

	if step.Manifest.Execution.Binary != "mv" {
		t.Errorf("expected binary 'mv', got '%s'", step.Manifest.Execution.Binary)
	}

	if len(step.RenderedArgs) != 2 || step.RenderedArgs[0] != "/tmp/input.mp4" || step.RenderedArgs[1] != "/storage/movies" {
		t.Errorf("rendered args mismatch: %v", step.RenderedArgs)
	}
}
