/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "list tickets",
	Long: `list tickets:
	* by default lists all active tickets:
		- tickets you are working on (workingON)
		- or tickets on which you wait for others to work on before you can resume your work (waitingFor)
	* you can also list only workignOn, waitingFor or closed passing the appropriate status parameter`,
	Args: func(cmd *cobra.Command, args []string) error {
		// further validation of the parameter can be implemented
		return nil // functions returns nil only if argument passes all checks
	},
	Aliases: []string{"l", "ls"},
	Run: func(cmd *cobra.Command, args []string) {
		var status string
		switch {
		case len(args) > 0:
			*(&status) = args[0]
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

		switch {
		case db.HasTable("tickets"):
			switch status {
			case "all":
				listAllTicketsOfStatus(db, "workingOn")
				fmt.Println()
				listAllTicketsOfStatus(db, "waitingFor")
				fmt.Println()
				listAllTicketsOfStatus(db, "closed")
			case "workingOn", "waitingFor", "closed":
				listAllTicketsOfStatus(db, status)
			default:
				listAllTicketsOfStatus(db, "workingOn")
				fmt.Println()
				listAllTicketsOfStatus(db, "waitingFor")
			}
		default:
			fmt.Printf("Table tickets does not exist. Need to add tickets first\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
