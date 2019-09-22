package cmd

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <ticket>",
	Short: "creates new ticket",
	Long:  `creates new ticket db entry, generates its templates and open them in a editor`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("new command requires a ticket argument")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		newTicketNumer := args[0]

		// Get Paths
		ticketsPath := getTicketsPath()                                // ~/tickets/
		newTicketPath := getNewTicketPath(ticketsPath, newTicketNumer) // ~/tickets/myNewTicket
		dbPath := getDBPath(ticketsPath, dbName)                       // ~/tickets/.tickets.db

		// Set Paths
		setTicketsPath(ticketsPath) // ~/tickets/

		// Create DB connection and migrates schema
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()
		db.AutoMigrate(&Ticket{})

		// check if ticket exist
		if !existsTicket(existsTicketDB(db, newTicketNumer), existsTicketPath(newTicketPath)) {
			// create newTicket
			newTicketDB(db, newTicketNumer, "")                             // DB entry
			setNewTicketPath(newTicketPath)                                 // ~/tickets/myNewTicket
			newTicketTemplate(templateNames, newTicketPath, newTicketNumer) // ~/tickets/myNewTicket/00 ...

			// open newTicket in sublime
			// no windows implementation yet
			workTicket(newTicketPath)
		}

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
