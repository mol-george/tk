/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "change ticket status",
	Long: `change ticket status to either:
		- closed if active (workingOn or waitingFor)
		- workingOn if closed
		- waitingFor if specified`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("change command requires at least one argument (ticket) but no more then two arguments (ticket, status)")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		newNote, _ := cmd.Flags().GetString("note")

		var ticketNumber, newStatus string
		switch {
		case len(args) == 2:
			*(&ticketNumber) = args[0]
			*(&newStatus) = args[1]
		case len(args) == 1:
			*(&ticketNumber) = args[0]
		}

		// Get Paths
		ticketsPath := getTicketsPath()          // ~/tickets/
		dbPath := getDBPath(ticketsPath, dbName) // ~/tickets/.tickets.db

		// Create DB connection and migrates schema
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

		status := getTicketStatus(db, ticketNumber)
		// fmt.Printf("ticketNumber: %s\noldStatus: %s\nnewStatus: %s\n", ticketNumber, status, newStatus)

		switch {
		case newStatus != "": //    if newStatus IS specified                  --> change to it
			changeStatus(db, ticketNumber, newStatus)
		case status == "closed": // if newStatus empty && status closed        --> workingON
			changeStatus(db, ticketNumber, "workingOn")
		default: //                 if newStatus empty && status is NOT closed --> closed
			changeStatus(db, ticketNumber, "closed")
		}

		if newNote != "" {
			changeNote(db, ticketNumber, newNote)
		}
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)
	changeCmd.Flags().StringP("note", "n", "", "info note")
}
