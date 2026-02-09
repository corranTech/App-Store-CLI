package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared/errfmt"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/update"
)

// Run executes the CLI using the provided args (not including argv[0]) and version string.
// It returns the intended process exit code.
func Run(args []string, versionInfo string) int {
	root := RootCommand(versionInfo)
	defer CleanupTempPrivateKeys()

	if err := root.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ExitSuccess
		}
		fmt.Fprint(os.Stderr, errfmt.FormatStderr(err))
		return ExitCodeFromError(err)
	}

	// Validate CI report flags after parsing
	if err := shared.ValidateReportFlags(); err != nil {
		fmt.Fprint(os.Stderr, errfmt.FormatStderr(err))
		return ExitUsage
	}

	if versionRequested {
		if err := root.Run(context.Background()); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return ExitUsage
			}
			fmt.Fprint(os.Stderr, errfmt.FormatStderr(err))
			return ExitCodeFromError(err)
		}
		return ExitSuccess
	}

	updateResult, err := update.CheckAndUpdate(context.Background(), update.Options{
		CurrentVersion: versionInfo,
		AutoUpdate:     true,
		NoUpdate:       shared.NoUpdate(),
		Output:         os.Stderr,
		ShowProgress:   shared.ProgressEnabled(),
		CheckInterval:  24 * time.Hour,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update check failed: %v\n", err)
	}
	if updateResult.Updated {
		exitCode, restartErr := update.Restart(updateResult.ExecutablePath, os.Args, os.Environ())
		if restartErr != nil {
			fmt.Fprintf(os.Stderr, "Restart failed after update: %v\n", restartErr)
		} else {
			return exitCode
		}
	}

	start := time.Now()
	runErr := root.Run(context.Background())
	elapsed := time.Since(start)

	// Get command name (full subcommand path)
	commandName := getCommandName(root, args)

	// Write JUnit report if requested
	if shared.ReportFormat() == shared.ReportFormatJUnit && shared.ReportFile() != "" {
		reportErr := writeJUnitReport(commandName, runErr, elapsed)
		if reportErr != nil {
			// Report write failure is a hard error - CI depends on it
			fmt.Fprintf(os.Stderr, "Error: failed to write JUnit report: %v\n", reportErr)
			if runErr == nil {
				return ExitError
			}
		}
	}

	if runErr != nil {
		var reported ReportedError
		if errors.As(runErr, &reported) {
			return ExitCodeFromError(runErr)
		}
		if errors.Is(runErr, flag.ErrHelp) {
			return ExitUsage
		}
		fmt.Fprint(os.Stderr, errfmt.FormatStderr(runErr))
		return ExitCodeFromError(runErr)
	}

	return ExitSuccess
}

// getCommandName extracts the full subcommand path from the parsed args.
// args is os.Args[1:] (without program name).
// It finds the first token matching a known subcommand name, then walks the tree.
func getCommandName(root *ffcli.Command, args []string) string {
	current := root
	path := []string{current.Name}

	for i := 0; i < len(args); i++ {
		found := false
		for _, sub := range current.Subcommands {
			if sub.Name == args[i] {
				path = append(path, sub.Name)
				current = sub
				found = true
				break
			}
		}
		if !found {
			// Check if it's a known subcommand at any level to skip flag values
			if !isKnownSubcommand(root, args[i]) {
				continue // Skip non-subcommand tokens (flag values)
			}
			break
		}
	}

	return strings.Join(path, " ")
}

// isKnownSubcommand checks if a token matches any subcommand in the tree.
func isKnownSubcommand(cmd *ffcli.Command, name string) bool {
	for _, sub := range cmd.Subcommands {
		if sub.Name == name {
			return true
		}
		if isKnownSubcommand(sub, name) {
			return true
		}
	}
	return false
}

// writeJUnitReport writes a JUnit XML report if --report junit --report-file is configured.
func writeJUnitReport(commandName string, runErr error, elapsed time.Duration) error {
	reportFile := shared.ReportFile()
	if reportFile == "" {
		return nil
	}

	testCase := shared.JUnitTestCase{
		Name:      commandName,
		Classname: commandName,
		Time:      elapsed,
	}

	if runErr != nil {
		testCase.Failure = "ERROR"
		testCase.Message = runErr.Error()
	}

	report := shared.JUnitReport{
		Tests:     []shared.JUnitTestCase{testCase},
		Timestamp: time.Now(),
		Name:      "asc",
	}

	return report.Write(reportFile)
}
