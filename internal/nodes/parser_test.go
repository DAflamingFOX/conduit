package nodes_test

import (
	"testing"

	"github.com/daflamingfox/conduit/internal/nodes"
)

func TestParseManifestJSONC(t *testing.T) {
	jsoncContent := `
	{
		// Single-line comment describing node
		"id": "conduit.file.test",
		"name": "Test Node",
		"category": "Testing",
		"version": "1.0.0",
		"description": "Test node manifest with comments",
		/*
		 Multi-line comment block
		*/
		"inputs": [
			{ "id": "file", "type": "file", "label": "Input File", "required": true }
		],
		"outputs": [
			{ "id": "success", "type": "file", "label": "Success" }
		],
		"parameters": [
			{ "id": "target", "label": "Target Path", "type": "text", "default": "/tmp/out" }
		],
		"execution": {
			"type": "command",
			"binary": "cp",
			"args": ["{{ inputs.file.path }}", "{{ parameters.target }}"]
		}, // Trailing comma supported by JWCC
	}
	`

	manifest, err := nodes.ParseManifest([]byte(jsoncContent))
	if err != nil {
		t.Fatalf("expected successful parsing of JSONC node manifest, got error: %v", err)
	}

	if manifest.ID != "conduit.file.test" {
		t.Errorf("expected ID 'conduit.file.test', got '%s'", manifest.ID)
	}
	if manifest.Name != "Test Node" {
		t.Errorf("expected Name 'Test Node', got '%s'", manifest.Name)
	}
	if len(manifest.Inputs) != 1 {
		t.Errorf("expected 1 input pin, got %d", len(manifest.Inputs))
	}
}

func TestRenderTemplate(t *testing.T) {
	tmpl := "mv {{ inputs.file.path }} {{ parameters.dest }}/{{ inputs.file.name }}"
	tv := nodes.TemplateValues{
		FilePath:   "/storage/video.mp4",
		FileName:   "video.mp4",
		WorkingDir: "/tmp/scratch",
		Params: map[string]interface{}{
			"dest": "/output/movies",
		},
	}

	rendered := nodes.RenderTemplate(tmpl, tv)
	expected := "mv /storage/video.mp4 /output/movies/video.mp4"

	if rendered != expected {
		t.Errorf("expected '%s', got '%s'", expected, rendered)
	}
}
