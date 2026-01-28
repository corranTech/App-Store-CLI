package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	certificatescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/certificates"
)

// CertificatesCommand returns the certificates command group.
func CertificatesCommand() *ffcli.Command {
	return certificatescli.CertificatesCommand()
}
