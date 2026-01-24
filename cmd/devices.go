package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// DevicesCommand returns the devices command with subcommands.
func DevicesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("devices", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "devices",
		ShortUsage: "asc devices <subcommand> [flags]",
		ShortHelp:  "Manage registered devices.",
		LongHelp: `Manage registered devices.

Examples:
  asc devices list
  asc devices list --platform IOS
  asc devices register --name "Device" --udid "UDID" --platform IOS
  asc devices update --id "DEVICE_ID" --status DISABLED`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			DevicesListCommand(),
			DevicesRegisterCommand(),
			DevicesUpdateCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// DevicesListCommand returns the devices list subcommand.
func DevicesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	platform := fs.String("platform", "", "Filter by platform(s), comma-separated")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc devices list [flags]",
		ShortHelp:  "List registered devices.",
		LongHelp: `List registered devices.

Examples:
  asc devices list
  asc devices list --platform IOS
  asc devices list --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("devices list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("devices list: %w", err)
			}

			platforms, err := normalizePlatforms(splitCSVUpper(*platform))
			if err != nil {
				return fmt.Errorf("devices list: %w", err)
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("devices list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.DevicesOption{
				asc.WithDevicesLimit(*limit),
				asc.WithDevicesNextURL(*next),
			}
			if len(platforms) > 0 {
				opts = append(opts, asc.WithDevicesPlatforms(platforms))
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithDevicesLimit(200))
				firstPage, err := client.GetDevices(requestCtx, paginateOpts...)
				if err != nil {
					return fmt.Errorf("devices list: failed to fetch: %w", err)
				}

				paginated, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetDevices(ctx, asc.WithDevicesNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("devices list: %w", err)
				}

				return printOutput(paginated, *output, *pretty)
			}

			resp, err := client.GetDevices(requestCtx, opts...)
			if err != nil {
				return fmt.Errorf("devices list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// DevicesRegisterCommand returns the devices register subcommand.
func DevicesRegisterCommand() *ffcli.Command {
	fs := flag.NewFlagSet("register", flag.ExitOnError)

	name := fs.String("name", "", "Device name")
	udid := fs.String("udid", "", "Device UDID")
	platform := fs.String("platform", "", "Device platform: "+strings.Join(signingPlatformList(), ", "))
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "register",
		ShortUsage: "asc devices register --name \"Device\" --udid \"UDID\" --platform IOS",
		ShortHelp:  "Register a device.",
		LongHelp: `Register a device.

Examples:
  asc devices register --name "Device" --udid "UDID" --platform IOS`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			nameValue := strings.TrimSpace(*name)
			if nameValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --name is required")
				return flag.ErrHelp
			}
			udidValue := strings.TrimSpace(*udid)
			if udidValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --udid is required")
				return flag.ErrHelp
			}
			if strings.TrimSpace(*platform) == "" {
				fmt.Fprintln(os.Stderr, "Error: --platform is required")
				return flag.ErrHelp
			}
			platformValue, err := normalizePlatform(*platform)
			if err != nil {
				return fmt.Errorf("devices register: %w", err)
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("devices register: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			attrs := asc.DeviceCreateAttributes{
				Name:     nameValue,
				UDID:     udidValue,
				Platform: platformValue,
			}
			resp, err := client.RegisterDevice(requestCtx, attrs)
			if err != nil {
				return fmt.Errorf("devices register: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// DevicesUpdateCommand returns the devices update subcommand.
func DevicesUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	id := fs.String("id", "", "Device ID")
	status := fs.String("status", "", "Device status: ENABLED or DISABLED")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc devices update --id \"DEVICE_ID\" --status ENABLED|DISABLED",
		ShortHelp:  "Update a device's status.",
		LongHelp: `Update a device's status.

Examples:
  asc devices update --id "DEVICE_ID" --status DISABLED`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			idValue := strings.TrimSpace(*id)
			if idValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			statusValue := strings.TrimSpace(*status)
			if statusValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --status is required")
				return flag.ErrHelp
			}
			normalizedStatus, err := normalizeDeviceStatus(statusValue)
			if err != nil {
				return fmt.Errorf("devices update: %w", err)
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("devices update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			attrs := asc.DeviceUpdateAttributes{Status: normalizedStatus}
			resp, err := client.UpdateDevice(requestCtx, idValue, attrs)
			if err != nil {
				return fmt.Errorf("devices update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}
