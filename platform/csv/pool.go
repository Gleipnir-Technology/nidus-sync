package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/rs/zerolog/log"
)

func ProcessJob(ctx context.Context, file_id int32) error {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return fmt.Errorf("Failed to get file %d from DB: %w", file_id, err)
	}
	r, err := userfile.NewFileReader(userfile.CollectionCSV, file.FileUUID)
	if err != nil {
		return fmt.Errorf("Failed to get filereader for %d: %w", file_id, err)
	}
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to read all CSV records for file %d: %w", file_id, err)
	}
	for _, rec := range records {
		log.Debug().Strs("rec", rec).Msg("Line")
	}
	return nil
}
