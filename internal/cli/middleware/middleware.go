package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// CommandFunc is the function signature for a command that uses an ASC client.
type CommandFunc func(ctx context.Context, client *asc.Client) error

// CommandResultFunc is the function signature for a command that returns a result.
type CommandResultFunc[T any] func(ctx context.Context, client *asc.Client) (*T, error)

// Run executes a command with common middleware applied (client + timeout).
// This is the simplest way to run a command that needs authentication and timeout.
func Run(ctx context.Context, cmd CommandFunc) error {
	client, err := shared.GetASCClient()
	if err != nil {
		return fmt.Errorf("failed to get ASC client: %w", err)
	}

	requestCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	return cmd(requestCtx, client)
}

// RunWithResult executes a command that returns a result with common middleware,
// and prints the result using the provided output format.
func RunWithResult[T any](ctx context.Context, output string, pretty bool, cmd CommandResultFunc[T]) error {
	client, err := shared.GetASCClient()
	if err != nil {
		return fmt.Errorf("failed to get ASC client: %w", err)
	}

	requestCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	result, err := cmd(requestCtx, client)
	if err != nil {
		return err
	}

	return shared.PrintOutput(result, output, pretty)
}

// RunWithUploadTimeout executes a command with upload timeout (for long-running operations).
func RunWithUploadTimeout(ctx context.Context, cmd CommandFunc) error {
	client, err := shared.GetASCClient()
	if err != nil {
		return fmt.Errorf("failed to get ASC client: %w", err)
	}

	requestCtx, cancel := shared.ContextWithUploadTimeout(ctx)
	defer cancel()

	return cmd(requestCtx, client)
}

// RunWithResultAndUploadTimeout executes a command with upload timeout and prints the result.
func RunWithResultAndUploadTimeout[T any](ctx context.Context, output string, pretty bool, cmd CommandResultFunc[T]) error {
	client, err := shared.GetASCClient()
	if err != nil {
		return fmt.Errorf("failed to get ASC client: %w", err)
	}

	requestCtx, cancel := shared.ContextWithUploadTimeout(ctx)
	defer cancel()

	result, err := cmd(requestCtx, client)
	if err != nil {
		return err
	}

	return shared.PrintOutput(result, output, pretty)
}

// GetClient retrieves the ASC client.
// This is a convenience function for commands that need direct access to the client.
func GetClient() (*asc.Client, error) {
	return shared.GetASCClient()
}

// ContextWithTimeout creates a context with standard request timeout.
func ContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return shared.ContextWithTimeout(ctx)
}

// ContextWithUploadTimeout creates a context with upload timeout.
func ContextWithUploadTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return shared.ContextWithUploadTimeout(ctx)
}

// PrintOutput prints the result in the specified format.
func PrintOutput(data any, format string, pretty bool) error {
	return shared.PrintOutput(data, format, pretty)
}

// TimingFunc is a callback for reporting command execution timing.
type TimingFunc func(duration time.Duration, command string)

// NoOpTiming is a no-op timing function.
func NoOpTiming() TimingFunc {
	return func(duration time.Duration, command string) {
		_ = duration
		_ = command
	}
}
