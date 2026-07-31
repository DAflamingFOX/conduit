package nodes

// IOType represents data input/output types (file, string, number, boolean, error).
type IOType string

const (
	IOTypeFile    IOType = "file"
	IOTypeString  IOType = "string"
	IOTypeNumber  IOType = "number"
	IOTypeBoolean IOType = "boolean"
	IOTypeError   IOType = "error"
)

// InputPin defines an input port on a node.
type InputPin struct {
	ID       string `json:"id"`
	Type     IOType `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}

// OutputPin defines an output port/branch on a node.
type OutputPin struct {
	ID    string `json:"id"`
	Type  IOType `json:"type"`
	Label string `json:"label"`
}

// ParameterType defines parameter input UI types.
type ParameterType string

const (
	ParamTypeText     ParameterType = "text"
	ParamTypeNumber   ParameterType = "number"
	ParamTypeSelect   ParameterType = "select"
	ParamTypeCheckbox ParameterType = "checkbox"
)

// ParameterDef defines user-configurable UI parameters for a node.
type ParameterDef struct {
	ID          string        `json:"id"`
	Label       string        `json:"label"`
	Type        ParameterType `json:"type"`
	Description string        `json:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Options     []string      `json:"options,omitempty"` // For select inputs
	Min         *float64      `json:"min,omitempty"`     // For numeric inputs
	Max         *float64      `json:"max,omitempty"`     // For numeric inputs
}

// ExecutionType defines how the node executes (command, script, internal).
type ExecutionType string

const (
	ExecTypeCommand  ExecutionType = "command"
	ExecTypeScript   ExecutionType = "script"
	ExecTypeInternal ExecutionType = "internal"
)

// ExecutionSpec describes execution binaries and template arguments.
type ExecutionSpec struct {
	Type        ExecutionType     `json:"type"`
	Binary      string            `json:"binary,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Script      string            `json:"script,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// Manifest represents a complete declarative node definition.
type Manifest struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Inputs      []InputPin     `json:"inputs"`
	Outputs     []OutputPin    `json:"outputs"`
	Parameters  []ParameterDef `json:"parameters"`
	Execution   ExecutionSpec  `json:"execution"`
}
