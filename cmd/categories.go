package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	categoriescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/categories"
)

// CategoriesCommand returns the categories command group.
func CategoriesCommand() *ffcli.Command {
	return categoriescli.CategoriesCommand()
}
