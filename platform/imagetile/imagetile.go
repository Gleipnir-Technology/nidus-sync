package imagetile

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/rs/zerolog/log"
)

//go:embed empty-tile.png
var emptyTileFS embed.FS

var ErrNoTile = errors.New("used placeholder tile")

var clientByOrgID = make(map[int32]*fieldseeker.FieldSeeker, 0)

func ImageAtPoint(ctx context.Context, org *models.Organization, level uint, lat, lng float64) ([]byte, error) {
	fssync, err := getFieldseeker(ctx, org)
	if err != nil {
		return []byte{}, fmt.Errorf("create fssync: %w", err)
	}
	map_service, err := aerialImageService(ctx, fssync.Arcgis)
	if err != nil {
		return []byte{}, fmt.Errorf("no map service: %w", err)
	}
	return map_service.TileGPS(ctx, level, lat, lng)
}
func ImageAtTile(ctx context.Context, org *models.Organization, level, y, x uint) ([]byte, error) {
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
	data, e := map_service.Tile(ctx, level, y, x)
	if e != nil {
		log.Error().Err(e).Msg("error getting tile")
		return []byte{}, fmt.Errorf("tile: %w", e)
	}
	// No data at this location, so supply the empty tile placeholder
	if len(data) == 0 {
		empty, err := emptyTileFS.ReadFile("empty-tile.png")
		if err != nil {
			return []byte{}, fmt.Errorf("read empty tile: %w", err)
		}
		return empty, ErrNoTile
	}
	return data, nil
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
func getFieldseeker(ctx context.Context, org *models.Organization) (*fieldseeker.FieldSeeker, error) {
	fssync, ok := clientByOrgID[org.ID]
	if ok {
		return fssync, nil
	}
	oauth, err := background.GetOAuthForOrg(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("get oauth for org: %w", err)
	}
	fssync, err = background.NewFieldSeeker(
		ctx,
		oauth,
	)
	clientByOrgID[org.ID] = fssync
	return fssync, nil
}
