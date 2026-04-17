package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/lob"
)

func main() {
	to := flag.String("to", "", "")
	from := flag.String("from", "", "")
	file := flag.String("file", "", "")
	color := flag.Bool("color", false, "")
	use_type := flag.String("use_type", "operational", "")

	// Parse the flags
	flag.Parse()

	key := os.Getenv("LOB_API_KEY")
	if key == "" {
		log.Println("LOB_API_KEY is empty")
		os.Exit(1)
	}

	client := lob.NewLob(key)
	ctx := context.TODO()
	req := lob.RequestLetterCreate{
		To:      *to,
		From:    *from,
		File:    *file,
		Color:   *color,
		UseType: *use_type,
	}
	letter, err := client.LetterCreate(ctx, req)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(2)
	}
	log.Printf("done. Letter: %s", letter.ID)
}
