package userfile

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

func audioFileContentWrite(audioUUID uuid.UUID, body io.Reader) error {
	return nil
}

var collectionToExtension map[Collection]string = map[Collection]string{
	CollectionAudioRaw:        "raw",
	CollectionAudioTranscoded: "ogg",
	CollectionCSV:             "csv",
	CollectionLogo:            "png",
	CollectionPublicImage:     "img",
	CollectionImageRaw:        "raw",
}
var collectionToSubdir map[Collection]string = map[Collection]string{
	CollectionAudioRaw:        "audio-raw",
	CollectionAudioTranscoded: "audio-transcoded",
	CollectionCSV:             "csv",
	CollectionLogo:            "logo",
	CollectionPublicImage:     "public-image",
	CollectionImageRaw:        "image-raw",
}

func fileContentPath(collection Collection, uid uuid.UUID) string {
	subdir, ok := collectionToSubdir[collection]
	if !ok {
		panic(fmt.Sprintf("No subdir for collection %d", int(collection)))
	}
	extension, ok := collectionToExtension[collection]
	return fmt.Sprintf("%s/%s/%s.%s", config.FilesDirectory, subdir, uid.String(), extension)
}

/*
	func fileContentWrite(body io.Reader, subdir string, uid uuid.UUID, extension string) error {
		// Create file in configured directory
		filepath := fileContentPath(subdir, uid, extension)
		dst, err := os.Create(filepath)
		if err != nil {
			log.Error().Err(err).Str("filepath", filepath).Msg("Failed to create file")
			return fmt.Errorf("Failed to create file at %s: %v", filepath, err)
		}
		defer dst.Close()

		// Copy rest of request body to file
		_, err = io.Copy(dst, body)
		if err != nil {
			return fmt.Errorf("Unable to save content of %s: %v", filepath, err)
		}
		return nil
	}
*/
func writeFileContent(w http.ResponseWriter, image_path string) {
	// Open the file
	file, err := os.Open(image_path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Image not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve image", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Get file info for Content-Length header
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get image information", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Copy file contents to response writer
	_, err = io.Copy(w, file)
	if err != nil {
		// Note: At this point, we've already started writing the response,
		// so we can't change the status code anymore. The best we can do
		// is log the error and abandon the connection.
		return
	}
}
