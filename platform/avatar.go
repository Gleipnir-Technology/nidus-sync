package platform

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
)

func AvatarCreate(ctx context.Context, u User, upload file.Upload) error {
	f, err := file.NewFileReader(file.CollectionAvatar, upload.UUID)
	// Decode image (supports PNG, JPG, GIF)
	img, err := imaging.Decode(f)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	// Resize to 200x200, maintaining aspect ratio
	avatar := imaging.Fill(img, 200, 200, imaging.Center, imaging.Lanczos)

	// Save or encode
	//filename := fmt.Sprintf("avatar-%s.jpg", upload.UUID.String())
	//err = imaging.Save(avatar, filename)
	//log.Info().Str("filename", filename).Msg("wrote avatar file")
	// Or encode to buffer: imaging.Encode(writer, avatar, imaging.JPEG)
	writer := &bytes.Buffer{}
	err = imaging.Encode(writer, avatar, imaging.PNG)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	err = file.FileContentWrite(writer, file.CollectionAvatar, upload.UUID)
	if err != nil {
		return fmt.Errorf("write content: %w", err)
	}
	log.Info().Str("uuid", upload.UUID.String()).Msg("wrote avatar file")
	return nil
}
