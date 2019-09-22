package cmd

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <ticket/all>",
	Short: "remove specific / all tickets db entries",
	Long:  `remove specific / all tickets db entries; it does not remove the ticket notes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("change command requires an argument - either: a ticket number or all")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		toBeRemoved := args[0]

		// Get Paths
		ticketsPath := getTicketsPath()          // ~/tickets/
		dbPath := getDBPath(ticketsPath, dbName) // ~/tickets/.tickets.db

		// Create DB connection and migrates schema
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

		switch toBeRemoved {
		case "all":
			removeAllTickets(db)
		default:
			removeTicket(db, toBeRemoved)

		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
