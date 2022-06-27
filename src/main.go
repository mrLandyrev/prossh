package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"table"
	"validators"
)

var profiles []Profile

func main() {
	loadProfiles()
	printProfiles()
	getUserCommandLoop()
}

func getUserCommandLoop() {
	for {
		fmt.Print("\n")
		userInput := getUserInputLoop("Enter profile index or command: ", []validators.IValidator{
			validators.MinLengthValidator{1},
		})
		switch userInput {
		case "exit":
			return
		case "new":
			executeNew()
		case "delete":
			executeDelete()
		default:
			index, err := strconv.Atoi(userInput)
			if err != nil {
				fmt.Println("Can not parse index")
				continue
			}
			if index <= 0 || index > len(profiles) {
				fmt.Println("Index out of range")
				continue
			}
			executeConnect(index - 1)
			return
		}
	}
}

func printProfiles() {
	table := table.Table{}
	table.AddColumn("Index")
	table.AddColumn("Name")

	for index, profile := range profiles {
		table.AddRow([]string{
			strconv.Itoa(index + 1),
			profile.Name,
		})
	}

	fmt.Print("Profiles list: \n\n")
	fmt.Print(table)
}

func executeConnect(index int) {
	profiles[index].Connect()
}

func executeNew() {
	newProfile := Profile{}

	names := []string{}
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	newProfile.Name = getUserInputLoop("Enter new profile name: ", []validators.IValidator{
		validators.MinLengthValidator{1},
		validators.CannotExistsValidator{names},
	})
	newProfile.User = getUserInputLoop("Enter user: ", []validators.IValidator{
		validators.MinLengthValidator{1},
	})
	newProfile.Adress = getUserInputLoop("Enter address: ", []validators.IValidator{
		validators.MinLengthValidator{1},
	})
	newProfile.Port = getUserInputLoop("Enter port (can be empty): ", []validators.IValidator{})
	newProfile.KeyPath = getUserInputLoop("Enter path to identity file (can be empty): ", []validators.IValidator{})

	profiles = append(profiles, newProfile)
	fmt.Print("New profile was created!\n\n")
	saveProfiles()
	printProfiles()
}

func executeDelete() {
	indexes := make([]string, len(profiles))
	for index := range profiles {
		indexes = append(indexes, strconv.Itoa(index+1))
	}
	indexStr := getUserInputLoop("Enter profile index to delete: ", []validators.IValidator{
		validators.ExistsValidator{indexes},
	})
	confirm := getUserInputLoop("Confirm delete (y or n): ", []validators.IValidator{
		validators.ExistsValidator{[]string{"y", "n"}},
	})
	if confirm == "n" {
		return
	}
	index, _ := strconv.Atoi(indexStr)
	profiles = append(profiles[:index-1], profiles[index:]...)
	fmt.Print("Profile was deleted!\n\n")
	saveProfiles()
	printProfiles()
}

func getUserInputLoop(message string, validators []validators.IValidator) (input string) {
	reader := bufio.NewReader(os.Stdin)
retry:
	fmt.Print(message)
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	for _, validator := range validators {
		if err := validator.Validate(input); len(err) > 0 {
			fmt.Println(err)
			goto retry
		}
	}
	return
}

func loadProfiles() {
	home, _ := os.UserHomeDir()
	dat, _ := os.ReadFile(home + "/.prossh/profiles")
	json.Unmarshal(dat, &profiles)
}

func saveProfiles() {
	data, _ := json.Marshal(&profiles)
	home, _ := os.UserHomeDir()
	os.WriteFile(home+"/.prossh/profiles", data, os.ModePerm)
}

type Profile struct {
	Name    string
	User    string
	Adress  string
	Port    string
	KeyPath string
}

func (profile Profile) Connect() {
	cmd := exec.Command("ssh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Args = append(cmd.Args, profile.User+"@"+profile.Adress)
	if len(profile.Port) > 0 {
		cmd.Args = append(cmd.Args, "-p", profile.Port)
	}
	if len(profile.KeyPath) > 0 {
		cmd.Args = append(cmd.Args, "-i", profile.KeyPath)
	}

	cmd.Run()
}
