package asc

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureXcodeCloudStdout(t *testing.T, fn func() error) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	os.Stdout = w

	err = fn()

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("close error: %v", closeErr)
	}
	os.Stdout = orig

	var buf bytes.Buffer
	if _, readErr := io.Copy(&buf, r); readErr != nil {
		t.Fatalf("read error: %v", readErr)
	}
	if err != nil {
		t.Fatalf("function error: %v", err)
	}

	return buf.String()
}

func TestPrintTable_XcodeCloudRunResult(t *testing.T) {
	result := &XcodeCloudRunResult{
		BuildRunID:        "run-123",
		BuildNumber:       42,
		WorkflowID:        "wf-456",
		WorkflowName:      "CI Build",
		GitReferenceID:    "ref-789",
		GitReferenceName:  "main",
		ExecutionProgress: "PENDING",
		CompletionStatus:  "",
		StartReason:       "MANUAL",
		CreatedDate:       "2026-01-22T10:00:00Z",
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintTable(result)
	})

	if !strings.Contains(output, "Build Run ID") {
		t.Fatalf("expected header in output, got: %s", output)
	}
	if !strings.Contains(output, "run-123") {
		t.Fatalf("expected build run ID in output, got: %s", output)
	}
	if !strings.Contains(output, "CI Build") {
		t.Fatalf("expected workflow name in output, got: %s", output)
	}
	if !strings.Contains(output, "PENDING") {
		t.Fatalf("expected execution progress in output, got: %s", output)
	}
}

func TestPrintMarkdown_XcodeCloudRunResult(t *testing.T) {
	result := &XcodeCloudRunResult{
		BuildRunID:        "run-123",
		BuildNumber:       42,
		WorkflowID:        "wf-456",
		WorkflowName:      "CI Build",
		GitReferenceID:    "ref-789",
		GitReferenceName:  "main",
		ExecutionProgress: "RUNNING",
		CompletionStatus:  "",
		StartReason:       "MANUAL",
		CreatedDate:       "2026-01-22T10:00:00Z",
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintMarkdown(result)
	})

	if !strings.Contains(output, "| Build Run ID |") {
		t.Fatalf("expected markdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "run-123") {
		t.Fatalf("expected build run ID in output, got: %s", output)
	}
	if !strings.Contains(output, "RUNNING") {
		t.Fatalf("expected execution progress in output, got: %s", output)
	}
}

func TestPrintTable_XcodeCloudStatusResult(t *testing.T) {
	result := &XcodeCloudStatusResult{
		BuildRunID:        "run-123",
		BuildNumber:       42,
		WorkflowID:        "wf-456",
		ExecutionProgress: "COMPLETE",
		CompletionStatus:  "SUCCEEDED",
		StartReason:       "MANUAL",
		CancelReason:      "",
		CreatedDate:       "2026-01-22T10:00:00Z",
		StartedDate:       "2026-01-22T10:01:00Z",
		FinishedDate:      "2026-01-22T10:05:00Z",
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintTable(result)
	})

	if !strings.Contains(output, "Build Run ID") {
		t.Fatalf("expected header in output, got: %s", output)
	}
	if !strings.Contains(output, "COMPLETE") {
		t.Fatalf("expected execution progress in output, got: %s", output)
	}
	if !strings.Contains(output, "SUCCEEDED") {
		t.Fatalf("expected completion status in output, got: %s", output)
	}
}

func TestPrintMarkdown_XcodeCloudStatusResult(t *testing.T) {
	result := &XcodeCloudStatusResult{
		BuildRunID:        "run-123",
		BuildNumber:       42,
		WorkflowID:        "wf-456",
		ExecutionProgress: "COMPLETE",
		CompletionStatus:  "FAILED",
		StartReason:       "MANUAL",
		CancelReason:      "",
		CreatedDate:       "2026-01-22T10:00:00Z",
		StartedDate:       "2026-01-22T10:01:00Z",
		FinishedDate:      "2026-01-22T10:05:00Z",
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintMarkdown(result)
	})

	if !strings.Contains(output, "| Build Run ID |") {
		t.Fatalf("expected markdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "FAILED") {
		t.Fatalf("expected completion status in output, got: %s", output)
	}
}

func TestPrintTable_CiProducts(t *testing.T) {
	resp := &CiProductsResponse{
		Data: []CiProductResource{
			{
				ID: "prod-1",
				Attributes: CiProductAttributes{
					Name:        "MyApp",
					BundleID:    "com.example.myapp",
					ProductType: "APP",
					CreatedDate: "2026-01-22T10:00:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintTable(resp)
	})

	if !strings.Contains(output, "Bundle ID") {
		t.Fatalf("expected header in output, got: %s", output)
	}
	if !strings.Contains(output, "com.example.myapp") {
		t.Fatalf("expected bundle ID in output, got: %s", output)
	}
}

func TestPrintMarkdown_CiProducts(t *testing.T) {
	resp := &CiProductsResponse{
		Data: []CiProductResource{
			{
				ID: "prod-1",
				Attributes: CiProductAttributes{
					Name:        "MyApp",
					BundleID:    "com.example.myapp",
					ProductType: "APP",
					CreatedDate: "2026-01-22T10:00:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintMarkdown(resp)
	})

	if !strings.Contains(output, "| ID | Name | Bundle ID |") {
		t.Fatalf("expected markdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "MyApp") {
		t.Fatalf("expected app name in output, got: %s", output)
	}
}

func TestPrintTable_CiWorkflows(t *testing.T) {
	resp := &CiWorkflowsResponse{
		Data: []CiWorkflowResource{
			{
				ID: "wf-1",
				Attributes: CiWorkflowAttributes{
					Name:             "CI Build",
					IsEnabled:        true,
					LastModifiedDate: "2026-01-22T10:00:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintTable(resp)
	})

	if !strings.Contains(output, "Enabled") {
		t.Fatalf("expected header in output, got: %s", output)
	}
	if !strings.Contains(output, "CI Build") {
		t.Fatalf("expected workflow name in output, got: %s", output)
	}
}

func TestPrintMarkdown_CiWorkflows(t *testing.T) {
	resp := &CiWorkflowsResponse{
		Data: []CiWorkflowResource{
			{
				ID: "wf-1",
				Attributes: CiWorkflowAttributes{
					Name:             "Deploy",
					IsEnabled:        false,
					LastModifiedDate: "2026-01-22T10:00:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintMarkdown(resp)
	})

	if !strings.Contains(output, "| ID | Name | Enabled |") {
		t.Fatalf("expected markdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "Deploy") {
		t.Fatalf("expected workflow name in output, got: %s", output)
	}
}

func TestPrintTable_CiBuildRuns(t *testing.T) {
	resp := &CiBuildRunsResponse{
		Data: []CiBuildRunResource{
			{
				ID: "run-1",
				Attributes: CiBuildRunAttributes{
					Number:            1,
					ExecutionProgress: CiBuildRunExecutionProgressComplete,
					CompletionStatus:  CiBuildRunCompletionStatusSucceeded,
					StartReason:       "MANUAL",
					CreatedDate:       "2026-01-22T10:00:00Z",
					StartedDate:       "2026-01-22T10:01:00Z",
					FinishedDate:      "2026-01-22T10:05:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintTable(resp)
	})

	if !strings.Contains(output, "Progress") {
		t.Fatalf("expected header in output, got: %s", output)
	}
	if !strings.Contains(output, "COMPLETE") {
		t.Fatalf("expected execution progress in output, got: %s", output)
	}
}

func TestPrintMarkdown_CiBuildRuns(t *testing.T) {
	resp := &CiBuildRunsResponse{
		Data: []CiBuildRunResource{
			{
				ID: "run-1",
				Attributes: CiBuildRunAttributes{
					Number:            1,
					ExecutionProgress: CiBuildRunExecutionProgressRunning,
					CompletionStatus:  "",
					StartReason:       "MANUAL",
					CreatedDate:       "2026-01-22T10:00:00Z",
					StartedDate:       "2026-01-22T10:01:00Z",
				},
			},
		},
	}

	output := captureXcodeCloudStdout(t, func() error {
		return PrintMarkdown(resp)
	})

	if !strings.Contains(output, "| ID | Build # |") {
		t.Fatalf("expected markdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "RUNNING") {
		t.Fatalf("expected execution progress in output, got: %s", output)
	}
}

func TestIsBuildRunComplete(t *testing.T) {
	tests := []struct {
		progress CiBuildRunExecutionProgress
		want     bool
	}{
		{CiBuildRunExecutionProgressPending, false},
		{CiBuildRunExecutionProgressRunning, false},
		{CiBuildRunExecutionProgressComplete, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.progress), func(t *testing.T) {
			got := IsBuildRunComplete(tt.progress)
			if got != tt.want {
				t.Errorf("IsBuildRunComplete(%s) = %v, want %v", tt.progress, got, tt.want)
			}
		})
	}
}

func TestIsBuildRunSuccessful(t *testing.T) {
	tests := []struct {
		status CiBuildRunCompletionStatus
		want   bool
	}{
		{CiBuildRunCompletionStatusSucceeded, true},
		{CiBuildRunCompletionStatusFailed, false},
		{CiBuildRunCompletionStatusErrored, false},
		{CiBuildRunCompletionStatusCanceled, false},
		{CiBuildRunCompletionStatusSkipped, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := IsBuildRunSuccessful(tt.status)
			if got != tt.want {
				t.Errorf("IsBuildRunSuccessful(%s) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
