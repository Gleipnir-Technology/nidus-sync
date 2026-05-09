package main

import (
	"context"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
)

func main() {
	err := config.Parse()
	if err != nil {
		log.Printf("failed on config: %v", err)
		os.Exit(1)
	}
	ctx := context.TODO()
	err = db.InitializeDatabase(ctx, config.PGDSN)
	if err != nil {
		log.Printf("failed on db: %v", err)
		os.Exit(2)
	}

	txn, err := db.BeginTxn(ctx)
	if err != nil {
		log.Printf("failed on txn: %v", err)
		os.Exit(3)
	}
	defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")
	log.Printf("doing address")
	gid := "openaddresses:address:us/ca/tulare-addresses-county:0dc28458fd03e3fa"
	address, err := public.AddressFromGID(ctx, txn, gid)
	if err != nil {
		log.Printf("failed on query: %v", err)
		os.Exit(4)
	}
	//log.Printf("address %d lat %f lng %f", address.ID, *address.LocationLatitude, *address.LocationLongitude)
	log.Printf("Address id %d location %s", address.ID, address.Location)
	lint.LogOnErrCtx(txn.Commit, ctx, "commit")

	/*
		log.Printf("doing comm")
		id := int64(1)
		comm, err := public.CommunicationFromID(ctx, id)
		if err != nil {
			log.Printf("failed on query: %v", err)
			os.Exit(4)
		}
		log.Printf("communication %d", comm.ID)
	*/
}
