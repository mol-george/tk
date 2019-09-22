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
	Use:   "list",
	Short: "list tickets",
	Long: `list tickets:
	* by default (without providing any flags) lists all active tickets:
		- tickets you are working on (workingON)
		- or tickets on which you wait for others to work on before you can resume your work (waitingFor)
	* you can also list only workignOn, waitingFor or closed passing the appropriate status flag value`,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		status, _ := cmd.Flags().GetString("status")

		// Get Paths
		ticketsPath := getTicketsPath()          // ~/tickets/
		dbPath := getDBPath(ticketsPath, dbName) // ~/tickets/.tickets.db

		// Create DB connection and migrates schema
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

		switch status {
		case "workingOn", "waitingFor", "closed":
			listAllTicketsOfStatus(db, status)
		default:
			listAllTicketsOfStatus(db, "workingOn")
			fmt.Println()
			listAllTicketsOfStatus(db, "waitingFor")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("status", "s", "", "possible values: workginOn, waitingFor, closed")
}
