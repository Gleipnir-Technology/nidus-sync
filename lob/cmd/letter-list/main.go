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
	letters, err := client.LetterList(ctx)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(2)
	}

	for _, letter := range letters {
		log.Printf("%s %s %s", letter.ID, letter.To.ID, letter.From.ID)
	}
}
