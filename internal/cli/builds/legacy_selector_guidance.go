package builds

import (
	"flag"
	"strconv"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

type trackedStringFlag struct {
	value string
	set   bool
}

func (f *trackedStringFlag) String() string {
	if f == nil {
		return ""
	}
	return f.value
}

func (f *trackedStringFlag) Set(value string) error {
	f.value = value
	f.set = true
	return nil
}

func (f *trackedStringFlag) Used() bool {
	return f != nil && f.set
}

func (f *trackedStringFlag) Value() string {
	if f == nil {
		return ""
	}
	return strings.TrimSpace(f.value)
}

type trackedBoolFlag struct {
	value bool
	set   bool
}

func (f *trackedBoolFlag) String() string {
	if f == nil {
		return "false"
	}
	return strconv.FormatBool(f.value)
}

func (f *trackedBoolFlag) Set(value string) error {
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	f.value = parsed
	f.set = true
	return nil
}

func (f *trackedBoolFlag) IsBoolFlag() bool {
	return true
}

func (f *trackedBoolFlag) Used() bool {
	return f != nil && f.set
}

func bindHiddenStringFlag(fs *flag.FlagSet, name string) *trackedStringFlag {
	value := &trackedStringFlag{}
	fs.Var(value, name, "DEPRECATED: use --build-id")
	shared.HideFlagFromHelp(fs.Lookup(name))
	return value
}

func bindHiddenBoolFlag(fs *flag.FlagSet, name string) *trackedBoolFlag {
	value := &trackedBoolFlag{}
	fs.Var(value, name, "DEPRECATED: use --latest")
	shared.HideFlagFromHelp(fs.Lookup(name))
	return value
}

func removedBuildFlagError(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return shared.UsageError("--build was removed; use --build-id BUILD_ID")
	}
	return shared.UsageErrorf("--build was removed; use --build-id %q", value)
}

func removedIDFlagError(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return shared.UsageError("--id was removed as a build selector; use --build-id BUILD_ID")
	}
	return shared.UsageErrorf("--id was removed as a build selector; use --build-id %q", value)
}

func removedNewestFlagError() error {
	return shared.UsageError("--newest was removed; use --latest")
}
