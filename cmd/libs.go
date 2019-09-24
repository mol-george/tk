package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gookit/color"
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
	start      = "=====================================================\n\n"
	middle     = "-----------------------------------------------------\n\n"
	end        = "+++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
)

// Ticket models ticket basic properties
// Further fields can be added: RAccount, AAccount
type Ticket struct {
	gorm.Model
	Number string
	Status string
	Note   string
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

func newTicketDB(db *gorm.DB, newTicketNumer, status, note string) {
	if status == "" {
		status = "workingon"
	}
	db.Create(&Ticket{Number: newTicketNumer, Status: status, Note: note})
}

func newTicketTemplate(templateNames []string, newTicketPath, newTicketNumer string) {
	body := "{{ .Number}}\n\n\n" + start + strings.Repeat(middle, 20) + end
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

// unused function
func listAllTickets(db *gorm.DB) {
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
}

func listAllTicketsOfStatus(db *gorm.DB, status string) {
	switch status {
	case "workingon":
		color.Red.Printf("=======\t %s \t=======\n", strings.ToUpper(status))
	case "sometimessoon":
		color.Danger.Printf("=======\t %s \t=======\n", strings.ToUpper(status))
	case "waitingfor":
		color.Yellow.Printf("=======\t %s \t=======\n", strings.ToUpper(status))
	case "closed":
		color.Secondary.Printf("=======\t %s \t=======\n", strings.ToUpper(status))
	}

	tickets := []Ticket{}
	db.Where(&Ticket{Status: status}).Find(&tickets)
	for _, ticket := range tickets {
		color.Notice.Printf("%s\t", ticket.Number)
		fmt.Printf("%s\n", ticket.Note)
	}
}

func getTicketStatus(db *gorm.DB, ticketNumber string) string {
	ticket := Ticket{}
	db.Where(&Ticket{Number: ticketNumber}).Find(&ticket)
	return ticket.Status
}

func changeStatus(db *gorm.DB, ticketNumer, newStatus string) {
	var ticket Ticket
	db.Model(&ticket).Where("number = ?", ticketNumer).Update("status", newStatus)
}

func changeNote(db *gorm.DB, ticketNumer, newNote string) {
	var ticket Ticket
	db.Model(&ticket).Where("number = ?", ticketNumer).Update("note", newNote)
}

func removeAllTickets(db *gorm.DB) {
	var ticket Ticket
	db.Delete(&ticket)
}

func removeTicket(db *gorm.DB, toBeRemoved string) {
	db.Where("number = ?", toBeRemoved).Delete(Ticket{})
}

func workTicket(ticketPath string) {
	cmd := exec.Command("subl", "-an", ticketPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func existsTicketDB(db *gorm.DB, newTicketNumer string) bool {
	ticket := Ticket{}
	db.Where(&Ticket{Number: newTicketNumer}).Find(&ticket)
	if ticket.Number != "" {
		return true
	}
	return false

}

func existsTicketPath(newTicketPath string) bool {
	if _, err := os.Stat(newTicketPath); !os.IsNotExist(err) {
		return true
	}
	return false

}

func existsTicket(existsTicketDB, existsTicketPath bool) bool {
	if existsTicketDB {
		fmt.Printf("Ticket Exists: Database\n")
	}
	if existsTicketPath {
		fmt.Printf("Ticket Exists: File System \n")
	}
	if existsTicketDB || existsTicketPath {
		return true
	}
	return false
}

func normalizeInput(rawInput string) string {
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	rawInput = reg.ReplaceAllString(*(&rawInput), "")
	rawInput = strings.ToLower(rawInput)

	var input string

	switch rawInput {
	case "workingon", "wo":
		*(&input) = "workingon"
	case "sometimessoon", "ss":
		*(&input) = "sometimessoon"
	case "waitingfor", "wf":
		*(&input) = "waitingfor"
	case "closed", "cl":
		*(&input) = "closed"
	case "all":
		*(&input) = "all"
	}
	return input
}
