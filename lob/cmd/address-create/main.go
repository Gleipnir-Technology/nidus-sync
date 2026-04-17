package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/lob"
)

func main() {
	name := flag.String("name", "", "The name of the address")
	line1 := flag.String("line1", "", "")
	city := flag.String("city", "", "")
	state := flag.String("state", "", "")
	zip := flag.String("zip", "", "")

	// Parse the flags
	flag.Parse()

	key := os.Getenv("LOB_API_KEY")
	if key == "" {
		log.Println("LOB_API_KEY is empty")
		os.Exit(1)
	}

	client := lob.NewLob(key)
	ctx := context.TODO()
	req := lob.RequestAddressCreate{
		AddressLine1: *line1,
		AddressCity:  *city,
		AddressState: *state,
		AddressZip:   *zip,
		Name:         *name,
	}
	addr, err := client.AddressCreate(ctx, req)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(2)
	}
	log.Printf("done. Address: %s", addr.ID)
}
