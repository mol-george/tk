package cmd

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <ticket> [\"note\"]",
	Short: "creates new ticket",
	Long:  `creates new ticket db entry, generates its templates and open them in a editor`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("change command requires at least one argument (ticket) but no more then two arguments (ticket, note)")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		newTicketNumer := args[0]
		note := args[1]

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
			newTicketDB(db, newTicketNumer, "", note)                       // DB entry
			setNewTicketPath(newTicketPath)                                 // ~/tickets/myNewTicket
			newTicketTemplate(templateNames, newTicketPath, newTicketNumer) // ~/tickets/myNewTicket/00 ...

			// open newTicket in sublime
			// no windows implementation yet
			workTicket(newTicketPath)
		} else {
			fmt.Printf("Ticket Status: %s\n", getTicketStatus(db, newTicketNumer))
		}

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
