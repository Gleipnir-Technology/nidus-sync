package imagetile

import (
	"context"
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
)

//go:embed empty-tile.png
var emptyTileFS embed.FS

var clientByOrgID = make(map[int32]*fieldseeker.FieldSeeker, 0)
var tileRasterPlaceholder *TileRaster

type TileRaster struct {
	Content       []byte
	IsPlaceholder bool
}

func ImageAtPoint(ctx context.Context, org *models.Organization, level uint, lat, lng float64) (*TileRaster, error) {
	fssync, err := getFieldseeker(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("create fssync: %w", err)
	}
	map_service, err := aerialImageService(ctx, fssync.Arcgis)
	if err != nil {
		return nil, fmt.Errorf("no map service: %w", err)
	}
	data, e := map_service.TileGPS(ctx, level, lat, lng)
	if e != nil {
		return nil, fmt.Errorf("tilegps: %w", e)
	}
	if len(data) == 0 {
		return TileRasterPlaceholder(), nil
	}
	return &TileRaster{
		Content:       data,
		IsPlaceholder: false,
	}, nil
}
func ImageAtTile(ctx context.Context, org *models.Organization, level, y, x uint) (*TileRaster, error) {
	oauth, err := background.GetOAuthForOrg(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("get oauth for org: %w", err)
	}
	fssync, err := background.NewFieldSeeker(
		ctx,
		oauth,
	)
	if err != nil {
		return nil, fmt.Errorf("create fssync: %w", err)
	}
	map_service, err := aerialImageService(ctx, fssync.Arcgis)
	if err != nil {
		return nil, fmt.Errorf("no map service: %w", err)
	}
	data, e := map_service.Tile(ctx, level, y, x)
	if e != nil {
		return nil, fmt.Errorf("tile: %w", e)
	}
	// No data at this location, so supply the empty tile placeholder
	if len(data) == 0 {
		return TileRasterPlaceholder(), nil
	}
	return &TileRaster{
		Content:       data,
		IsPlaceholder: false,
	}, nil
}
func TileRasterPlaceholder() *TileRaster {
	if tileRasterPlaceholder != nil {
		return tileRasterPlaceholder
	}
	empty, err := emptyTileFS.ReadFile("empty-tile.png")
	if err != nil {
		panic(fmt.Sprintf("Failed to read empty-tile.png: %v", err))
	}
	tileRasterPlaceholder = &TileRaster{
		Content:       empty,
		IsPlaceholder: true,
	}
	return tileRasterPlaceholder
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
