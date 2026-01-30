package productpages

import "github.com/peterbourgon/ff/v3/ffcli"

// Command returns the product pages command group.
func Command() *ffcli.Command {
	return ProductPagesCommand()
}
