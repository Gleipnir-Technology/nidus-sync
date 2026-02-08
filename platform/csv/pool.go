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
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to read header of CSV for file %d: %w", file_id, err)
	}
	err = validateHeader(header)
	if err != nil {
		addImportError(file, err)
	}
	for {
		row, err := reader.Read()
		if err != nil {
			return fmt.Errorf("Failed to read all CSV records for file %d: %w", file_id, err)
		}
		log.Debug().Strs("row", row).Msg("Line")
	}
	return nil
}

func addImportError(file *models.FileuploadFile, err error) {
	log.Debug().Err(err).Int32("file_id", file.ID).Msg("Fake add import error")
}
func validateHeader(row []string) error {
	return nil
}
