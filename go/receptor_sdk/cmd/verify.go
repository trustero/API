package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

// Set up the 'verify' CLI subcommand.
var verifyCmd = &cobra.Command{
	Use:   "verify <trustero_access_token>|'dryrun'",
	Short: "Verify read-only access to a service provider account.",
	Long: `
Verify read-only access to a service provider account.  Verify command
decodes the base64 URL encoded credentials from the '--credentials' command
line flag and check it's validity.  If 'dryrun' is specified instead of a
Trustero access token, the verify command will not report the results to
Trustero and instead print the results to console.`,
	Args: cobra.MaximumNArgs(1),
	RunE: verify,
}

func init() {
	addReceptorFlags(verifyCmd)
}

// Cobra executes this function on verify command.
func verify(_ *cobra.Command, args []string) (err error) {
	// Run receptor's Verify function and report results to Trustero
	err = invokeWithContext(args[0],
		func(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

			// Call receptor's Verify method
			verifyResult := toVerifyResult(receptorImpl.Verify(credentials))

			// Notify behavior is different for the verify command.  When the '--notify' command line
			// flag is provided on a verify command, verify only notify Trustero of the command
			// status and does NOT invoke the Verified Trustero RPC method to save the credential
			// in the receptor record.
			if len(receptor_sdk.Notify) > 0 {
				_ = notify(rc, "verify", verifyResult.Message, err)
				return
			}

			// Let Trustero know if the service provider account credentials are valid.
			_, err = rc.Verified(context.Background(), verifyResult)
			return
		})
	return
}

func toVerifyResult(ok bool, err error) *receptor_v1.Credential {
	var message string
	if err != nil {
		message = "error"
	} else if ok {
		message = "successful"
	} else {
		message = "failed"
	}

	return &receptor_v1.Credential{ReceptorObjectId: receptor_sdk.ReceptorId, Message: message, IsCredentialValid: ok}
}