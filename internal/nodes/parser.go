package nodes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tailscale/hujson"
)

// ParseManifest parses raw JSON or JSONC/JWCC (JSON with comments) bytes into a Manifest struct.
func ParseManifest(data []byte) (*Manifest, error) {
	// Parse with hujson to support comments (// and /* */) and trailing commas
	ast, err := hujson.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSONC node manifest: %w", err)
	}

	// Standardize strips comments and trailing commas to produce standard RFC-8259 JSON
	ast.Standardize()
	standardJSON := ast.Pack()

	var m Manifest
	if err := json.Unmarshal(standardJSON, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON node manifest into struct: %w", err)
	}

	if m.ID == "" {
		return nil, fmt.Errorf("node manifest is missing required 'id' field")
	}
	if m.Name == "" {
		return nil, fmt.Errorf("node manifest '%s' is missing required 'name' field", m.ID)
	}

	return &m, nil
}

// TemplateValues map used for string template substitutions.
type TemplateValues struct {
	FilePath   string
	FileName   string
	WorkingDir string
	Params     map[string]interface{}
}

// RenderTemplate substitutes {{ inputs.file.path }}, {{ working_dir }}, and {{ parameters.<key> }} placeholders.
func RenderTemplate(template string, tv TemplateValues) string {
	result := template
	result = strings.ReplaceAll(result, "{{ inputs.file.path }}", tv.FilePath)
	result = strings.ReplaceAll(result, "{{ inputs.file.name }}", tv.FileName)
	result = strings.ReplaceAll(result, "{{ working_dir }}", tv.WorkingDir)

	for k, v := range tv.Params {
		placeholder := fmt.Sprintf("{{ parameters.%s }}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
	}

	return result
}
