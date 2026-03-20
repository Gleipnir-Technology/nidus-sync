package platform

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	"github.com/aarondl/opt/omit"
	//"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/oauth"
	"github.com/rs/zerolog/log"
)

//go:embed empty-tile.png
var emptyTileFS embed.FS

func GetTile(ctx context.Context, w http.ResponseWriter, org Organization, z, y, x uint) error {
	if org.model.ArcgisMapServiceID.IsNull() {
		return fmt.Errorf("no map service ID set")
	}
	map_service_id := org.model.ArcgisMapServiceID.MustGet()
	tile_path := tilePath(map_service_id, z, y, x)
	tile_row, err := models.TileCachedImages.Query(
		models.SelectWhere.TileCachedImages.ArcgisID.EQ(map_service_id),
		models.SelectWhere.TileCachedImages.X.EQ(int32(x)),
		models.SelectWhere.TileCachedImages.Y.EQ(int32(y)),
		models.SelectWhere.TileCachedImages.Z.EQ(int32(z)),
	).One(ctx, db.PGInstance.BobDB)
	if err == nil {
		var tile *TileRaster
		if tile_row.IsEmpty {
			tile = TileRasterPlaceholder()
		} else {
			tile, err = loadTileFromDisk(tile_path)
			if err != nil {
				return fmt.Errorf("load tile from disk: %w", err)
			}
		}
		log.Debug().Uint("z", z).Uint("y", y).Uint("x", x).Bool("is empty", tile_row.IsEmpty).Msg("tile from cache")
		return writeTile(w, tile)
	}
	if err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("query db: %w", err)
	}
	image, err := ImageAtTile(ctx, org.model, uint(z), uint(y), uint(x))
	if err != nil {
		return fmt.Errorf("image at tile: %w", err)
	}
	if !image.IsPlaceholder {
		err = saveTileToDisk(image, tile_path)
		if err != nil {
			return fmt.Errorf("save tile: %w", err)
		}
	}
	_, err = models.TileCachedImages.Insert(&models.TileCachedImageSetter{
		ArcgisID: omit.From(map_service_id),
		X:        omit.From(int32(x)),
		Y:        omit.From(int32(y)),
		Z:        omit.From(int32(z)),
		IsEmpty:  omit.From(image.IsPlaceholder),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("save to db: %w", err)
	}
	log.Debug().Uint("z", z).Uint("y", y).Uint("x", x).Bool("placeholder", image.IsPlaceholder).Msg("caching tile")
	return writeTile(w, image)
}
func ImageAtPoint(ctx context.Context, org Organization, level uint, lat, lng float64) (*TileRaster, error) {
	fssync, err := getFieldseeker(ctx, org.model)
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
func loadTileFromDisk(tile_path string) (*TileRaster, error) {
	file, err := os.Open(tile_path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()
	img, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("readall from %s: %w", tile_path, err)
	}
	return &TileRaster{
		Content:       img,
		IsPlaceholder: false,
	}, nil
}
func saveTileToDisk(image *TileRaster, tile_path string) error {
	parent := filepath.Dir(tile_path)
	err := os.MkdirAll(parent, 0750)
	if err != nil {
		return fmt.Errorf("mkdirall: %w", err)
	}
	err = os.WriteFile(tile_path, image.Content, 0644)
	if err != nil {
		return fmt.Errorf("write image file: %w", err)
	}
	return nil
}
func tilePath(map_service_id string, z, y, x uint) string {
	return fmt.Sprintf("%s/tile-cache/%s/%d/%d/%d.raw", config.FilesDirectory, map_service_id, z, y, x)
}

func writeTile(w http.ResponseWriter, image *TileRaster) error {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image.Content)))
	_, err := io.Copy(w, bytes.NewBuffer(image.Content))
	if err != nil {
		return fmt.Errorf("io.copy: %w", err)
	}
	return nil
}

var clientByOrgID = make(map[int32]*fieldseeker.FieldSeeker, 0)
var tileRasterPlaceholder *TileRaster

type TileRaster struct {
	Content       []byte
	IsPlaceholder bool
}

func ImageAtTile(ctx context.Context, org *models.Organization, level, y, x uint) (*TileRaster, error) {
	oauth, err := oauth.GetOAuthForOrg(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("get oauth for org: %w", err)
	}
	if oauth == nil {
		return nil, fmt.Errorf("get oauth for org nil oauth.")
	}
	fssync, err := newFieldSeeker(
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
	oauth, err := oauth.GetOAuthForOrg(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("get oauth for org: %w", err)
	}
	fssync, err = newFieldSeeker(
		ctx,
		oauth,
	)
	clientByOrgID[org.ID] = fssync
	return fssync, nil
}
