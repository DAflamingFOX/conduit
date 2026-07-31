package engine

import (
	"fmt"
	"log/slog"

	"github.com/daflamingfox/conduit/internal/nodes"
)

// ExecutionStep represents a single node invocation step prepared for a worker.
type ExecutionStep struct {
	StepID         string          `json:"step_id"`
	NodeInstanceID string          `json:"node_instance_id"`
	Manifest       *nodes.Manifest `json:"manifest"`
	RenderedArgs   []string        `json:"rendered_args"`
	Context        *Context        `json:"context"`
}

// Evaluator evaluates and traverses a FlowGraph given an execution Context and Node Registry.
type Evaluator struct {
	graph    *FlowGraph
	registry *nodes.Registry
}

// NewEvaluator constructs an Evaluator.
func NewEvaluator(graph *FlowGraph, registry *nodes.Registry) (*Evaluator, error) {
	if err := graph.Validate(); err != nil {
		return nil, fmt.Errorf("invalid flow graph: %w", err)
	}
	return &Evaluator{
		graph:    graph,
		registry: registry,
	}, nil
}

// PrepareStep renders execution command arguments for a specific node instance in the flow.
func (e *Evaluator) PrepareStep(nodeInstanceID string, ctx *Context) (*ExecutionStep, error) {
	var targetInstance *NodeInstance
	for i := range e.graph.Nodes {
		if e.graph.Nodes[i].ID == nodeInstanceID {
			targetInstance = &e.graph.Nodes[i]
			break
		}
	}

	if targetInstance == nil {
		return nil, fmt.Errorf("node instance '%s' not found in flow graph", nodeInstanceID)
	}

	manifest, ok := e.registry.Get(targetInstance.NodeManifestID)
	if !ok {
		return nil, fmt.Errorf("node manifest '%s' not registered", targetInstance.NodeManifestID)
	}

	tv := nodes.TemplateValues{
		FilePath:   ctx.CurrentFile.Path,
		FileName:   ctx.CurrentFile.Name,
		WorkingDir: ctx.WorkingDir,
		Params:     targetInstance.Parameters,
	}

	renderedArgs := make([]string, len(manifest.Execution.Args))
	for i, arg := range manifest.Execution.Args {
		renderedArgs[i] = nodes.RenderTemplate(arg, tv)
	}

	slog.Debug("prepared execution step",
		"node_instance_id", nodeInstanceID,
		"manifest_id", manifest.ID,
		"binary", manifest.Execution.Binary,
		"args", renderedArgs,
	)

	return &ExecutionStep{
		StepID:         nodeInstanceID,
		NodeInstanceID: nodeInstanceID,
		Manifest:       manifest,
		RenderedArgs:   renderedArgs,
		Context:        ctx,
	}, nil
}
