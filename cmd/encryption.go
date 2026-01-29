package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	encryptioncli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/encryption"
)

// EncryptionCommand returns the encryption command group.
func EncryptionCommand() *ffcli.Command {
	return encryptioncli.EncryptionCommand()
}
