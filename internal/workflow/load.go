package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	// ErrWorkflowRead indicates workflow file read failure.
	ErrWorkflowRead = errors.New("read workflow")
	// ErrWorkflowParseJSON indicates workflow JSON decode failure.
	ErrWorkflowParseJSON = errors.New("parse workflow JSON")
)

// DefaultPath is the default location for the workflow definition file.
const DefaultPath = ".asc/workflow.json"

// Load reads and validates a workflow definition file.
func Load(path string) (*Definition, error) {
	def, err := LoadUnvalidated(path)
	if err != nil {
		return nil, err
	}
	if errs := Validate(def); len(errs) > 0 {
		return nil, errs[0]
	}
	return def, nil
}

// LoadUnvalidated reads and parses a workflow definition file without validation.
func LoadUnvalidated(path string) (*Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWorkflowRead, err)
	}

	var def Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWorkflowParseJSON, err)
	}

	return &def, nil
}
