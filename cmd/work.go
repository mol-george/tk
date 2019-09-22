package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// workCmd represents the work command
var workCmd = &cobra.Command{
	Use:   "work <ticket>",
	Short: "work ticket",
	Long:  `work ticket opens ticket notes in sublimeText`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("change command requires a ticket argument")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"w"},
	Run: func(cmd *cobra.Command, args []string) {
		ticketNumber := args[0]

		// Get Paths
		ticketsPath := getTicketsPath()                           // ~/tickets/
		ticketPath := getNewTicketPath(ticketsPath, ticketNumber) // ~/tickets/myNewTicket

		// work
		workTicket(ticketPath)
	},
}

func init() {
	rootCmd.AddCommand(workCmd)
}
