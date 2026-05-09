package platform

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	"github.com/aarondl/opt/omit"
	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/oauth"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/rs/zerolog/log"
)

//go:embed empty-tile.png
var emptyTileFS embed.FS

func GetTile(ctx context.Context, w http.ResponseWriter, org Organization, use_placeholder bool, z, y, x uint) error {
	return getTileFlyover(ctx, w, org.model, use_placeholder, z, y, x)
}
func GetTileFlyoverLatLng(ctx context.Context, w http.ResponseWriter, org *models.Organization, use_placeholder bool, level uint, lat, lng float64) error {
	y, x := LatLngToTile(level, lat, lng)
	return getTileFlyover(ctx, w, org, use_placeholder, level, y, x)
}
func GetTileSatelliteLatLng(ctx context.Context, w http.ResponseWriter, level uint, lat, lng float64) error {
	y, x := LatLngToTile(level, lat, lng)
	return getTileSatellite(ctx, w, level, y, x)
}

func ImageAtPoint(ctx context.Context, org Organization, level uint, lat, lng float64) (*TileRaster, error) {
	return imageAtPoint(ctx, org.model, level, lat, lng)
}

// LatLngToTile converts GPS coordinates to ArcGIS tile coordinates
func LatLngToTile(level uint, lat, lng float64) (row, column uint) {
	// Get number of tiles per dimension at this zoom level
	numTiles := math.Pow(2, float64(level))

	// Convert longitude to tile column
	// Range: -180 to 180 degrees maps to 0 to numTiles
	column = uint(math.Floor((lng + 180.0) / 360.0 * numTiles))

	// Convert latitude to tile row using Mercator projection
	// First convert lat to radians
	latRad := lat * math.Pi / 180.0

	// Apply Mercator projection formula
	// This maps latitude from -85.0511 to 85.0511 degrees to 0 to numTiles
	mercatorY := 0.5 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/(2*math.Pi)
	row = uint(math.Floor(mercatorY * numTiles))

	// Ensure values are within valid range
	if column < 0 {
		column = 0
	} else if column >= uint(numTiles) {
		column = uint(numTiles) - 1
	}

	if row < 0 {
		row = 0
	} else if row >= uint(numTiles) {
		row = uint(numTiles) - 1
	}

	return row, column
}

// Writes a random tile from the cache. This is a very odd thing to do, it's for testing
func WriteTileRandom(ctx context.Context, w http.ResponseWriter) error {
	tile_rows, err := models.TileCachedImages.Query(
		sm.Where(psql.Quote("is_empty").EQ(psql.Arg(false))),
		sm.Limit(100),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("get tiles: %w", err)
	}
	tile_row := tile_rows[rand.Intn(len(tile_rows))]
	service, err := models.FindTileService(ctx, db.PGInstance.BobDB, tile_row.ServiceID)
	if err != nil {
		return fmt.Errorf("get service: %w", err)
	}
	tile_path := tilePath(service.Name, uint(tile_row.Z), uint(tile_row.Y), uint(tile_row.X))
	var tile *TileRaster
	if tile_row.IsEmpty {
		tile = TileRasterPlaceholder()
	} else {
		tile, err = loadTileFromDisk(tile_path)
		if err != nil {
			return fmt.Errorf("load tile from disk: %w", err)
		}
	}
	log.Debug().Int32("z", tile_row.Z).Int32("y", tile_row.Y).Int32("x", tile_row.X).Bool("is empty", tile_row.IsEmpty).Msg("random tile")
	return writeTile(w, tile)
}
func cacheImage(ctx context.Context, image *TileRaster, map_service *models.TileService, z, y, x uint) error {
	var err error
	if !image.IsPlaceholder {
		tile_path := tilePath(map_service.Name, z, y, x)
		err = saveTileToDisk(image, tile_path)
		if err != nil {
			return fmt.Errorf("save tile: %w", err)
		}
	}
	_, err = models.TileCachedImages.Insert(&models.TileCachedImageSetter{
		ServiceID: omit.From(map_service.ID),
		X:         omit.From(int32(x)),
		Y:         omit.From(int32(y)),
		Z:         omit.From(int32(z)),
		IsEmpty:   omit.From(image.IsPlaceholder),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("save to db: %w", err)
	}
	log.Debug().Str("service", map_service.Name).Uint("z", z).Uint("y", y).Uint("x", x).Bool("placeholder", image.IsPlaceholder).Msg("caching tile")
	return nil
}
func getTileCached(ctx context.Context, map_service *models.TileService, z, y, x uint) (*TileRaster, bool, error) {
	tile_path := tilePath(map_service.Name, z, y, x)
	tile_row, err := models.TileCachedImages.Query(
		models.SelectWhere.TileCachedImages.ServiceID.EQ(map_service.ID),
		models.SelectWhere.TileCachedImages.X.EQ(int32(x)),
		models.SelectWhere.TileCachedImages.Y.EQ(int32(y)),
		models.SelectWhere.TileCachedImages.Z.EQ(int32(z)),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("query db: %w", err)
	}
	if tile_row.IsEmpty {
		return TileRasterPlaceholder(), true, nil
	}
	tile, err := loadTileFromDisk(tile_path)
	if err != nil {
		return nil, false, fmt.Errorf("load tile from disk: %w", err)
	}
	//log.Debug().Uint("z", z).Uint("y", y).Uint("x", x).Bool("is empty", tile_row.IsEmpty).Msg("tile from cache")
	return tile, false, nil
}
func getTileFlyover(ctx context.Context, w http.ResponseWriter, org *models.Organization, use_placeholder bool, z, y, x uint) error {
	if org.ArcgisMapServiceID.IsNull() {
		return fmt.Errorf("no map service ID set")
	}
	map_service_id := org.ArcgisMapServiceID.MustGet()
	map_service, err := models.TileServices.Query(
		models.SelectWhere.TileServices.ArcgisID.EQ(map_service_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("get map service: %w", err)
	}
	cached_tile, is_placeholder, err := getTileCached(ctx, map_service, z, y, x)
	if err != nil {
		return fmt.Errorf("get cached tile: %w", err)
	}
	if is_placeholder && !use_placeholder {
		return fmt.Errorf("only a placeholder is available at %d %d %d", z, y, x)
	}
	if cached_tile != nil {
		return writeTile(w, cached_tile)
	}
	image, err := ImageAtTile(ctx, org, uint(z), uint(y), uint(x))
	if err != nil {
		return fmt.Errorf("image at tile: %w", err)
	}
	err = cacheImage(ctx, image, map_service, z, y, x)
	if err != nil {
		return fmt.Errorf("cache image: %w", err)
	}
	return writeTile(w, image)
}
func getTileSatellite(ctx context.Context, w http.ResponseWriter, z, y, x uint) error {
	map_service_id := "stadia"
	map_service, err := models.TileServices.Query(
		models.SelectWhere.TileServices.Name.EQ(map_service_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("get map service: %w", err)
	}
	cached_tile, is_placeholder, err := getTileCached(ctx, map_service, z, y, x)
	if err != nil {
		return fmt.Errorf("get cached tile: %w", err)
	}
	if is_placeholder {
		return fmt.Errorf("only a placeholder is available at %d %d %d", z, y, x)
	}
	if cached_tile != nil {
		return writeTile(w, cached_tile)
	}
	client := stadia.NewStadiaMaps(config.StadiaMapsAPIKey)
	data, err := client.TileRaster(ctx, z, y, x)
	if err != nil {
		return fmt.Errorf("stadia tile raster: %w", err)
	}
	tile := TileRaster{
		Content:       data,
		IsPlaceholder: false,
	}
	err = cacheImage(ctx, &tile, map_service, z, y, x)
	if err != nil {
		return fmt.Errorf("cache image: %w", err)
	}
	return writeTile(w, &tile)
}
func imageAtPoint(ctx context.Context, org *models.Organization, level uint, lat, lng float64) (*TileRaster, error) {
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
func loadTileFromDisk(tile_path string) (*TileRaster, error) {
	file, err := os.Open(tile_path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer lint.LogOnErr(file.Close, "close tile file")
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
	if oauth == nil {
		return nil, fmt.Errorf("no live oauth for org %d", org.ID)
	}
	fssync, err = newFieldSeeker(
		ctx,
		oauth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fieldseeker: %w", err)
	}
	clientByOrgID[org.ID] = fssync
	return fssync, nil
}
