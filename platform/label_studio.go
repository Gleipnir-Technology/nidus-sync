package platform

import (
	"context"
	"fmt"
)

func initializeLabelStudio() error {
	return nil
	/*
		// Initialize the minio client
		//minioBucket := os.Getenv("S3_BUCKET")

		var err error
		labelStudioClient, err = createLabelStudioClient()
		if err != nil {
			return fmt.Errorf("Failed to create label studio client: %w", err)
		}
		// Get the project we are going to upload to
		labelStudioProject, err = findLabelStudioProject(labelStudioClient, "Nidus Speech-to-Text Transcriptions")
		if err != nil {
			return fmt.Errorf("Failed to find the label studio project: %w", err)
		}
		minioClient, err = createMinioClient()
		if err != nil {
			return fmt.Errorf("Failed to create minio client: %w", err)
		}
		return nil
	*/
}
func jobLabelStudioAudioCreate(ctx context.Context, row_id int32) error {
	return fmt.Errorf("label studio integration has been disabled")
	/*
		customer := os.Getenv("CUSTOMER")
		if customer == "" {
			return errors.New("You must specify a CUSTOMER env var")
		}
		note, err := noteAudioGetLatest(ctx, job.UUID.String())
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to get note %s", note.UUID))
		}

		if note.Version != 1 {
			return errors.New(fmt.Sprintf("Got version %d of %s", note.Version, note.UUID))
		}
		task, err := findMatchingTask(labelStudioClient, project, customer, note)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to search for a task: %v", err))
		}
		// We already have a task, nothing to do.
		if task != nil {
			return nil
		}

		err = createTask(labelStudioClient, project, minioClient, minioBucket, customer, note)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create a task: %v", err))
		}
		return nil
	*/
}


