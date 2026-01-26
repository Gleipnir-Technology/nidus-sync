package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
)

func main() {
	var password string
	scanValue("Please enter your password : ", &password)

	hash, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

}

func scanValue(message string, result *string) {
	fmt.Printf(message)
	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); !ok {
		log.Fatal(errors.New("Failed to scan input"))
	}
	*result = scanner.Text()
}
