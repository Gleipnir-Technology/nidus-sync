package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
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
	tile_path := fmt.Sprintf("%s/tile-cache/%s/%d/%d/%d.raw", config.FilesDirectory, map_service_id, z, y, x)
	file, err := os.Open(tile_path)
	if err == nil {
		defer file.Close()
		img, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("readall from %s: %w", tile_path, err)
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))
		_, err = io.Copy(w, bytes.NewBuffer(img))
		if err != nil {
			return fmt.Errorf("copy bytes from %s: %w", tile_path)
		}
		return nil
	}
	content, err := imagetile.ImageAtTile(ctx, org, uint(z), uint(y), uint(x))
	if err != nil {
		if errors.Is(err, imagetile.ErrNoTile) {
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			_, err = io.Copy(w, bytes.NewBuffer(content))
			if err != nil {
				return fmt.Errorf("write image file: %w", err)
			}
			return nil
		}
		return fmt.Errorf("image at tile: %w", err)
	}
	parent := filepath.Dir(tile_path)
	err = os.MkdirAll(parent, 0750)
	if err != nil {
		return fmt.Errorf("mkdirall: %w", err)
	}
	err = os.WriteFile(tile_path, content, 0644)
	if err != nil {
		return fmt.Errorf("write image file: %w", err)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	_, err = io.Copy(w, bytes.NewBuffer(content))
	if err != nil {
		return fmt.Errorf("write image file: %w", err)
	}
	return nil
}
