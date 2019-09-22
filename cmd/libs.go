package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	ticketsPath   string
	templateNames = []string{"00", "01", "02", "03", "04", "05"}
)

const (
	dbName     = ".tickets.db"
	ticketsDir = "tickets"
	start      = "=====================================================\n"
	middle     = "-----------------------------------------------------\n\n"
	end        = "+++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
)

// Ticket ...
type Ticket struct {
	gorm.Model
	Number string
	Status string
}

func getTicketsPath() string {
	// constructs & returns ticketsPath if not set as per different OSes
	switch runtime.GOOS {
	case "linux", "darwin":
		ticketsPath = path.Join(os.Getenv("HOME"), ticketsDir)
	case "windows":
		ticketsPath = filepath.Join(
			os.Getenv("HOMEDRIVE"),
			os.Getenv("HOMEPATH"),
			ticketsDir)
	default:
		log.Fatal("Cannot identify the Operating System")
	}
	if _, err := os.Stat(ticketsPath); os.IsNotExist(err) {
		os.MkdirAll(ticketsPath, os.ModePerm)
	}
	return ticketsPath
}

func setTicketsPath(ticketsPath string) {
	if _, err := os.Stat(ticketsPath); os.IsNotExist(err) {
		os.MkdirAll(ticketsPath, os.ModePerm)
	}
}

func getNewTicketPath(ticketsPath, newTicket string) string {
	newTicketPath := path.Join(ticketsPath, newTicket)
	return newTicketPath
}

func setNewTicketPath(newTicketPath string) {
	if _, err := os.Stat(newTicketPath); os.IsNotExist(err) {
		os.MkdirAll(newTicketPath, os.ModePerm)
	}
}

func getDBPath(ticketsPath, dbName string) string {
	return path.Join(ticketsPath, dbName)
}

func newTicketDB(db *gorm.DB, newTicketNumer string, status string) {
	if status == "" {
		status = "workingOn"
	}
	db.Create(&Ticket{Number: newTicketNumer, Status: status})
}

func newTicketTemplate(templateNames []string, newTicketPath, newTicketNumer string) {
	body := start + "{{ .Number}}\n" + strings.Repeat(middle, 20) + end
	t, err := template.New("ticketTemplate").Parse(body)
	if err != nil {
		panic(err)
	}
	for _, templateName := range templateNames {
		templatePath := path.Join(newTicketPath, templateName)
		f, err := os.Create(templatePath)
		if err != nil {
			log.Println("create file: ", err)
			return
		}
		err = t.Execute(f, &Ticket{Number: newTicketNumer})
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func listAllTickets(db *gorm.DB) {
	fmt.Printf("ALL Tickets\n")
	var ticket Ticket

	rows, err := db.Model(&Ticket{}).Rows()
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("TicketNumber\tTicketStatus\n")
	for rows.Next() {
		db.ScanRows(rows, &ticket)
		fmt.Printf("%s\t\t%s\n", ticket.Number, ticket.Status)
	}
	fmt.Printf("-------------------------\n\n")
}

func listAllTicketsOfStatus(db *gorm.DB, status string) {
	fmt.Printf("All %s Tickets:\n", status)
	tickets := []Ticket{}

	db.Where(&Ticket{Status: status}).Find(&tickets)
	for _, ticket := range tickets {
		fmt.Println(ticket.Number)
	}
	fmt.Printf("-------------------------\n\n")

}

func closeTicket(db *gorm.DB, newTicketNumer string) {
	var ticket Ticket
	db.Model(&ticket).Where("number = ?", newTicketNumer).Update("status", "closed")
}

func deleteAllTickets(db *gorm.DB) {
	var ticket Ticket
	db.Delete(&ticket)
}

func openTicket(newTicketPath string) {
	cmd := exec.Command("subl", "-an", newTicketPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

/*
func to launch editor
newTicketTemplate("123")
newTicket(db, "345", "")
newTicket(db, "789", "waitingFor")
listAllTickets(db)
closeTicket(db, "345")
listAllTicketsOfStatus(db, "workingOn")
listAllTicketsOfStatus(db, "waitingFor")
listAllTicketsOfStatus(db, "closed")
deleteAllTickets(db)
*/
