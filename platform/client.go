package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func EnsureClient(ctx context.Context, client uuid.UUID, user_agent string) error {
	_, err := models.PublicreportClients.Query(
		models.SelectWhere.PublicreportClients.UUID.EQ(client),
	).One(ctx, db.PGInstance.BobDB)
	if err == nil {
		//log.Debug().Str("client", client.String()).Msg("already exists")
		return nil
	} else if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("failed existing client %s: %w", client.String(), err)
	}
	_, err = models.PublicreportClients.Insert(&models.PublicreportClientSetter{
		Created:   omit.From(time.Now()),
		UserAgent: omit.From(user_agent),
		UUID:      omit.From(client),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("insert client: %w", err)
	}
	log.Debug().Str("client", client.String()).Str("ua", user_agent).Msg("Created client")
	return nil
}
