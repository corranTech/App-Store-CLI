// Package workflow is a standalone workflow runner for .asc/workflow.json files.
// It has ZERO imports from the rest of the codebase — only Go stdlib.
package workflow

import (
	"encoding/json"
	"fmt"
)

// Definition is the top-level .asc/workflow.json schema.
type Definition struct {
	Env       map[string]string   `json:"env,omitempty"`
	BeforeAll string              `json:"before_all,omitempty"`
	AfterAll  string              `json:"after_all,omitempty"`
	Error     string              `json:"error,omitempty"`
	Workflows map[string]Workflow `json:"workflows"`
}

// Workflow is a named automation sequence.
type Workflow struct {
	Description string            `json:"description,omitempty"`
	Private     bool              `json:"private,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Steps       []Step            `json:"steps"`
}

// Step is one executable action in a workflow.
// Bare JSON strings unmarshal to Step{Run: "..."} as shorthand.
type Step struct {
	Run      string            `json:"run,omitempty"`
	Workflow string            `json:"workflow,omitempty"`
	Name     string            `json:"name,omitempty"`
	If       string            `json:"if,omitempty"`
	With     map[string]string `json:"with,omitempty"`
}

// UnmarshalJSON handles the flexible step format:
//   - bare string → Step{Run: "..."}
//   - object → normal unmarshal
func (s *Step) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err == nil {
		s.Run = raw
		return nil
	}

	type stepAlias Step
	var alias stepAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("step must be a string or object: %w", err)
	}
	*s = Step(alias)
	return nil
}
