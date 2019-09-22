package cmd

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <ticketNumber>",
	Short: "creates new ticket",
	Long:  `creates new ticket db entry, generates its templates and open them in a editor`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("new command requires a ticket argument")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Run: func(cmd *cobra.Command, args []string) {
		newTicketNumer := args[0]

		// Gets Paths
		ticketsPath := getTicketsPath()                                // ~/tickets/
		newTicketPath := getNewTicketPath(ticketsPath, newTicketNumer) // ~/tickets/myNewTicket
		dbPath := getDBPath(ticketsPath, dbName)

		// Sets Paths
		setTicketsPath(ticketsPath)     // ~/tickets/
		setNewTicketPath(newTicketPath) // ~/tickets/myNewTicket

		// Creates DB connection and migrates schema
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()
		db.AutoMigrate(&Ticket{})

		// create newTicket
		newTicketDB(db, newTicketNumer, "")
		newTicketTemplate(templateNames, newTicketPath, newTicketNumer)

		// open ticket in sublime
		// no windows implementation yet
		openTicket(newTicketPath)

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
