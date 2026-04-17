package main

import (
	"context"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/lob"
)

func main() {
	key := os.Getenv("LOB_API_KEY")
	if key == "" {
		log.Println("LOB_API_KEY is empty")
		os.Exit(1)
	}

	client := lob.NewLob(key)
	ctx := context.TODO()
	addresses, err := client.AddressList(ctx)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(2)
	}

	for _, addr := range addresses {
		log.Printf("%s %s %s: %s %s, %s, %s", addr.ID, addr.Name, addr.Company, addr.AddressLine1, addr.AddressCity, addr.AddressState, addr.AddressCountry, addr.AddressZip)
	}
}
