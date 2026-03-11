package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aarondl/opt/omit"
	//"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/imagetile"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func getTile(w http.ResponseWriter, r *http.Request, org *models.Organization, user *models.User) {
	x_str := chi.URLParam(r, "x")
	y_str := chi.URLParam(r, "y")
	z_str := chi.URLParam(r, "z")

	x, err := strconv.Atoi(x_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(y_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	z, err := strconv.Atoi(z_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	err = handleTile(r.Context(), w, org, uint(z), uint(y), uint(x))
	if err != nil {
		log.Error().Err(err).Msg("failed to do tile")
		http.Error(w, "failed to do tile", http.StatusInternalServerError)
		return
	}
}
func handleTile(ctx context.Context, w http.ResponseWriter, org *models.Organization, z, y, x uint) error {
	if org.ArcgisMapServiceID.IsNull() {
		return fmt.Errorf("no map service ID set")
	}
	map_service_id := org.ArcgisMapServiceID.MustGet()
	tile_path := tilePath(map_service_id, z, y, x)
	tile_row, err := models.TileCachedImages.Query(
		models.SelectWhere.TileCachedImages.ArcgisID.EQ(map_service_id),
		models.SelectWhere.TileCachedImages.X.EQ(int32(x)),
		models.SelectWhere.TileCachedImages.Y.EQ(int32(y)),
		models.SelectWhere.TileCachedImages.Z.EQ(int32(z)),
	).One(ctx, db.PGInstance.BobDB)
	if err == nil {
		var tile *imagetile.TileRaster
		if tile_row.IsEmpty {
			tile = imagetile.TileRasterPlaceholder()
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
	image, err := imagetile.ImageAtTile(ctx, org, uint(z), uint(y), uint(x))
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
func loadTileFromDisk(tile_path string) (*imagetile.TileRaster, error) {
	file, err := os.Open(tile_path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()
	img, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("readall from %s: %w", tile_path, err)
	}
	return &imagetile.TileRaster{
		Content:       img,
		IsPlaceholder: false,
	}, nil
}
func saveTileToDisk(image *imagetile.TileRaster, tile_path string) error {
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

func writeTile(w http.ResponseWriter, image *imagetile.TileRaster) error {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image.Content)))
	_, err := io.Copy(w, bytes.NewBuffer(image.Content))
	if err != nil {
		return fmt.Errorf("io.copy: %w", err)
	}
	return nil
}
