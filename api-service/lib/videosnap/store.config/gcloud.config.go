package storeconfig

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"telegram-service/utils"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadFilesToGCS(bucketName, serviceAccountPath, prefix string, filePaths []string, localDelete bool) {
	// Check if the service account file exists
	if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
		utils.ErrorLog(fmt.Errorf("service account file does not exist at path: %s", serviceAccountPath), err.Error())
		return
	}

	// Create a context
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		utils.ErrorLog(fmt.Errorf("failed to create storage client: %v", err), err.Error())
		return
	}
	defer client.Close()

	// Upload each file
	for _, filePath := range filePaths {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			utils.ErrorLog(fmt.Errorf("failed to open file %s: %v", filePath, err), err.Error())
			continue
		}

		// Ensure the file is closed after upload
		func() {
			defer file.Close()

			// Determine object name
			objectName := filepath.Join(prefix, filepath.Base(filePath))
			bucket := client.Bucket(bucketName)
			object := bucket.Object(objectName)

			// Create a writer for the object
			writer := object.NewWriter(ctx)
			writer.Metadata = map[string]string{
				"uploaded-at": time.Now().Format(time.RFC3339),
			}

			// Copy file contents to the writer
			if _, err := io.Copy(writer, file); err != nil {
				utils.ErrorLog(fmt.Errorf("failed to write file %s to GCS: %v", filePath, err), err.Error())
				return
			}

			// Finalize the upload
			if err := writer.Close(); err != nil {
				utils.ErrorLog(fmt.Errorf("failed to finalize GCS upload for file %s: %v", filePath, err), err.Error())
				return
			}

			utils.InfoLog(fmt.Sprintf("File %s uploaded to bucket %s successfully", filePath, bucketName), string(utils.SuccessMessage))

			// Delete local file if needed
			if localDelete {
				if err := os.Remove(filePath); err != nil {
					utils.ErrorLog(fmt.Errorf("failed to delete local file %s: %v", filePath, err), err.Error())
				} else {
					utils.WarnLog(fmt.Sprintf("Local file %s deleted successfully", filePath), string(utils.SuccessMessage))
				}
			}
		}()
	}
}
