package main

import (
	"encoding/csv"
	"fmt"
	"github.com/dimiro1/banner"
	"github.com/gofiber/fiber/v3"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Contact struct {
	Name        string
	PhoneNumber string
}

// Lista imion
var names = []string{
	"Pilar Sauer", "Julio Mathews", "Ahmad Kilgore", "Kolten Geer", "Barry Ostrander", "Lizet Daly", "Juliann Freeland", "Darryl Rivers", "Juana Arnett", "Cierra Lance",
	"Ulysses Gipson", "Rachelle Dix", "Karyme Locke", "Madeleine Ledford", "Stacy Isom", "Kalvin Cheek", "Julia Waller", "Brennen Reichert", "Jane Flint", "Jakob Gabriel",
	"Benjamin Hinkle", "Amirah Mchenry", "Jolie Braswell", "Jett Hays", "Tessa Hodgson", "Jazlynn Grubbs", "Tyler Wynn", "Domenic Penny", "Simran Romeo", "Bryant Foley",
	"Yusuf Pacheco", "Jackelyn Allred", "Shantel Begay", "Liam Carver", "Aryanna Vera", "Jamison Corley", "Alvin Cornejo", "Ashli Barnhart", "Tate Montanez", "Gordon Montanez",
	"Ezra Desantis", "Kenyon Meador", "Rylan Schultz", "Frederick Pickett", "Malia Bittner", "Tyree Whitt", "Kari Andres", "Maximilian Coyle", "Tierney Stoner", "Grant Hu",
	"Rhett Littleton", "Ciara Turley", "Josue Longoria", "Parker Nielsen", "Kai Mulligan", "Javen Workman", "Derick Hawks", "Jennie Anderson", "Hallie Santos", "Infant Moran",
	"Leighton Hamby", "Ansley Sommer", "Austin Noriega", "Aysha Hawk", "Leif Hirsch", "Maci Haag", "Sterling Bourne", "Delanie Cardwell", "Taniya Santoro", "Devonte Gallardo",
	"Briza Ellsworth", "Kailey Shipman", "Armani Christman", "Janaya Blackwell", "Arman Bigelow", "Maritza Salerno", "Bryant Moore", "Jaylin Bryant", "Jeff Pitt", "Caley Lemus",
	"Maddie Watt", "Stacie Krueger", "Belen Denson", "Tavian Crain", "Grecia Christopher", "Triniti Downs", "Maximiliano Eubanks", "Daquan Grimes", "Kia Satterfield", "Sarah Honeycutt",
	"Holly Chavez", "June Pearce", "Demond Hendricks", "Kenton Larkin", "Davin Valentin", "Leslie Hadley", "Gary Serna", "Jackie Day", "Katherine Mueller", "Jayleen Grover",
	"Christion Daigle", "Terence Carpenter", "Shane Allen", "Kacie Lentz", "Bethany Greer", "Carleigh Stark", "Dalila Stewart", "Efren Aldridge", "Leroy Steffen", "Lyric Stubbs",
	"Astrid Dunaway", "Kavon Parks", "Martha Pruett", "Annalise Forrester", "Brett Fox", "Mohamed Leal", "Mariam Bernstein", "Marisol Barrientos", "Fabiola Doan", "Shlomo Ko",
}

func GenerateRandomName() string {
	index := rand.Intn(len(names))
	name := names[index]
	names = append(names[:index], names[index+1:]...)
	return name
}

func GenerateRandomPhoneNumber() string {
	return fmt.Sprintf("+48 %03d %03d %03d", rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))
}

func GenerateContacts() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Create("./http/contacts.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	contacts := make(map[string]bool)
	for len(contacts) < 100 {
		contact := Contact{
			Name:        GenerateRandomName(),
			PhoneNumber: GenerateRandomPhoneNumber(),
		}

		contactKey := contact.Name + contact.PhoneNumber
		if _, exists := contacts[contactKey]; exists {
			continue
		}

		contacts[contactKey] = true

		err := writer.Write([]string{contact.Name, contact.PhoneNumber})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	writer.Flush()
}

// Główna funkcja
func main() {

	banner.Init(os.Stdout, true, true, strings.NewReader("    {{.AnsiColor.BrightBlue}}\n     _______     __   __     _____        _\n    |___  / |   /_ | / /    / ____|      | |\n        / /| |    | |/ /_   | |  __  ___  | |     __ _ _ __   __ _\n       / / | |    | | '_ \\  | | |_ |/ _ \\ | |    / _` | '_ \\ / _` |\n      / /__| |____| | (_) | | |__| | (_) || |___| (_| | | | | (_| |\n     /_____|______|_|\\___/   \\_____|\\___/ |______\\__,_|_| |_|\\__, |\n                                     ______                  __/ |\n                                    |______|                |___/{{.AnsiColor.Default}}\n\n    {{ .AnsiColor.BrightGreen }}\n    GoVersion: {{ .GoVersion }}\n    GOOS: {{ .GOOS }}\n    GOARCH: {{ .GOARCH }}\n    NumCPU: {{ .NumCPU }}\n    {{.AnsiColor.Default}}"))

	// Wygeneruj kontakty (100) (Wywołaj funkcję)
	GenerateContacts()

	fmt.Println("Starting server on port 1999")

	// Utwórz nowy definicję fiber app
	app := fiber.New()

	// Podawaj pliki z folderu http i wsakż na index.html
	app.Static("/", "./http", fiber.Static{
		Index: "index.html",
	})

	// Podawaj wygenerowane kontakty
	app.Get("/contacts", func(c fiber.Ctx) error {
		file, err := os.Open("./http/contacts.csv")
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		contactsList := make([]Contact, len(records))
		for i, record := range records {
			contactsList[i] = Contact{
				Name:        record[0],
				PhoneNumber: record[1],
			}
		}

		return c.JSON(contactsList)
	})

	// API do nowych kontaktów (dodawanie i pobieranie)
	app.Post("/addContact", func(c fiber.Ctx) error {
		name := c.FormValue("name")
		phoneNumber := c.FormValue("phoneNumber")

		file, err := os.OpenFile("./http/custom-contacts.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		err = writer.Write([]string{name, phoneNumber})
		if err != nil {
			return err
		}

		return c.SendString("Kontakt został dodany!")
	})

	// Podawaj plik dodanych kontaków
	app.Get("/custom-contacts", func(c fiber.Ctx) error {
		file, err := os.Open("./http/custom-contacts.csv")
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		customContactsList := make([]Contact, len(records))
		for i, record := range records {
			customContactsList[i] = Contact{
				Name:        record[0],
				PhoneNumber: record[1],
			}
		}

		return c.JSON(customContactsList)
	})

	// Rozpocznij web serwer na porcie 1999
	app.Listen(":1999")
}
