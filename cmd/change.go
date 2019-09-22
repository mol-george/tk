/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

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
		if len(args) < 1 {
			return errors.New("change command requires a ticket argument")
		}
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		ticketNumber := args[0]
		newStatus, _ := cmd.Flags().GetString("newStatus")

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

		if (status == "workingOn" || status == "waitingFor") && (newStatus == "" || newStatus == "closed") {
			changeStatus(db, ticketNumber, "closed")
		} else if status == "closed" && (newStatus == "" || newStatus == "workingOn") {
			changeStatus(db, ticketNumber, "workingOn")
		} else if status == "closed" && newStatus == "waitingFor" {
			changeStatus(db, ticketNumber, "waitingFor")
		} else if (status == "workingOn" || status == "waitingFor") && (newStatus == "workingOn" || newStatus == "waitingFor") {
			changeStatus(db, ticketNumber, newStatus)
		} else {
			fmt.Println("no else/if was matche")
		}
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)
	changeCmd.Flags().StringP("newStatus", "n", "", "change ticket status to new status")

}
