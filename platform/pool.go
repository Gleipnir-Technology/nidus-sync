package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

type PoolUpload struct {
	ID int32
}

func NewPoolUpload(ctx context.Context, u *models.User, upload userfile.FileUpload) (PoolUpload, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to begin transaction: %w", err)
	}

	file, err := models.FileuploadFiles.Insert(&models.FileuploadFileSetter{
		ContentType:    omit.From(upload.ContentType),
		Created:        omit.From(time.Now()),
		CreatorID:      omit.From(u.ID),
		Deleted:        omitnull.FromPtr[time.Time](nil),
		Name:           omit.From(upload.Name),
		OrganizationID: omit.From(u.OrganizationID),
		Status:         omit.From(enums.FileuploadFilestatustypeUploaded),
		SizeBytes:      omit.From(int32(upload.SizeBytes)),
		FileUUID:       omit.From(upload.UUID),
	}).One(ctx, txn)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to create file upload: %w", err)
	}
	_, err = models.FileuploadCSVS.Insert(&models.FileuploadCSVSetter{
		FileID: omit.From(file.ID),
		Type:   omit.From(enums.FileuploadCsvtypePoollist),
	}).One(ctx, txn)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to create csv: %w", err)
	}
	txn.Commit(ctx)
	return PoolUpload{
		ID: file.ID,
	}, nil
}
