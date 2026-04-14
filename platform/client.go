package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
)

func EnsureClient(ctx context.Context, client uuid.UUID, user_agent string) error {
	_, err := models.PublicreportClients.Query(
		models.SelectWhere.PublicreportClients.UUID.EQ(client),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return fmt.Errorf("failed existing client %s: %w", client.String(), err)
	}
	models.PublicreportClients.Insert(&models.PublicreportClientSetter{
		Created:   omit.From(time.Now()),
		UserAgent: omit.From(user_agent),
		UUID:      omit.From(client),
	}).One(ctx, db.PGInstance.BobDB)
	return nil
}
