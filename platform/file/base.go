package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

var collectionToExtension map[Collection]string = map[Collection]string{
	CollectionAudioNormalized: "ogg",
	CollectionAudioRaw:        "raw",
	CollectionAudioTranscoded: "ogg",
	CollectionAvatar:          "png",
	CollectionCSV:             "csv",
	CollectionLogo:            "png",
	CollectionMailerPDF:       "pdf",
	CollectionPublicImage:     "img",
	CollectionImageRaw:        "raw",
}
var collectionToSubdir map[Collection]string = map[Collection]string{
	CollectionAudioNormalized: "audio-normalized",
	CollectionAudioRaw:        "audio-raw",
	CollectionAudioTranscoded: "audio-transcoded",
	CollectionAvatar:          "avatar",
	CollectionCSV:             "csv",
	CollectionLogo:            "logo",
	CollectionMailerPDF:       "mailer",
	CollectionPublicImage:     "public-image",
	CollectionImageRaw:        "image-raw",
}

func ContentPath(collection Collection, id string) string {
	return fileContentPath(collection, id)
}
func ContentPathUUID(collection Collection, uid uuid.UUID) string {
	return fileContentPathUUID(collection, uid)
}
func collectionName(collection Collection) string {
	n, ok := collectionToSubdir[collection]
	if !ok {
		return "unknown"
	}
	return n
}
func fileContentPath(collection Collection, id string) string {
	subdir, ok := collectionToSubdir[collection]
	if !ok {
		panic(fmt.Sprintf("No subdir for collection %d", int(collection)))
	}
	extension, ok := collectionToExtension[collection]
	return fmt.Sprintf("%s/%s/%s.%s", config.FilesDirectory, subdir, id, extension)
}
func fileContentPathUUID(collection Collection, uid uuid.UUID) string {
	return fileContentPath(collection, uid.String())
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
	defer lint.LogOnErr(file.Close, "close file")

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
