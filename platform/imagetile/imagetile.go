package imagetile

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
)

func ImageAtPoint(ctx context.Context, org *models.Organization, level uint, lat, lng float64) ([]byte, error) {
	oauth, err := background.GetOAuthForOrg(ctx, org)
	if err != nil {
		return []byte{}, fmt.Errorf("get oauth for org: %w", err)
	}
	fssync, err := background.NewFieldSeeker(
		ctx,
		oauth,
	)
	if err != nil {
		return []byte{}, fmt.Errorf("create fssync: %w", err)
	}
	map_service, err := aerialImageService(ctx, fssync.Arcgis)
	if err != nil {
		return []byte{}, fmt.Errorf("no map service: %w", err)
	}
	return map_service.TileGPS(ctx, level, lat, lng)
}

func aerialImageService(ctx context.Context, gis *arcgis.ArcGIS) (*arcgis.MapService, error) {
	map_services, err := gis.MapServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("aerial image service: %w", err)
	}
	for _, ms := range map_services {
		return &ms, nil
	}
	return nil, fmt.Errorf("non found")
}
