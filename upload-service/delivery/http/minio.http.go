package http

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"time"
	"upload-service/utils"

	"github.com/google/uuid" // Import the UUID package
	"github.com/minio/minio-go/v7"
)

// UploadHandler handles file uploads to MinIO
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	utils.HttpSuccessResponse(w, "Upload Service ", http.StatusCreated, string(utils.SuccessMessage))
}

func UploadFile(bucketName string, minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.HttpSuccessResponse(w, "Method not allowed", http.StatusMethodNotAllowed, string(utils.ErrorMessage))
			return
		}
		err := r.ParseMultipartForm(10 << 20) // Max upload size: 10 MB
		if err != nil {

			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}
		defer file.Close()

		// Generate a new UUID for the filename and get the file extension
		extension := filepath.Ext(handler.Filename)    // Get the original file extension
		newFileName := uuid.New().String() + extension // Append the extension to the UUID
		contentType := handler.Header.Get("Content-Type")

		// Check if the bucket exists, and create it if it doesn't
		exists, err := minioClient.BucketExists(r.Context(), bucketName)
		if err != nil {
			log.Printf("Error checking if bucket exists: %v", err)
			http.Error(w, "Error checking bucket existence", http.StatusInternalServerError)
			return
		}
		if !exists {
			err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Printf("Error creating bucket: %v", err)
				http.Error(w, "Error creating bucket", http.StatusInternalServerError)
				return
			}
		}

		// Upload the file to MinIO with the new UUID filename
		_, err = minioClient.PutObject(r.Context(), bucketName, newFileName, file, handler.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			log.Printf("Failed to upload file: %v", err)
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
			return
		}

		// Create a response object
		response := map[string]string{
			//	"message":   "File uploaded successfully",
			//	"bucket":    bucketName,
			"file_name": newFileName, // Return only the new UUID file name with extension
		}
		utils.HttpSuccessResponse(w, response, http.StatusCreated, string(utils.SuccessMessage))
	}
}
func UploadWithOutChnageFilerNameHandler(bucketName string, minioClient *minio.Client, domain string, maxFileSize int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.HttpSuccessResponse(w, "Method not allowed", http.StatusMethodNotAllowed, string(utils.ErrorMessage))
			return
		}

		err := r.ParseMultipartForm(maxFileSize) // Use max file size from config
		if err != nil {
			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file")
		if err != nil {
			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}
		defer file.Close()

		// Use the original file name
		fileName := handler.Filename

		// Check if the bucket exists, and create it if it doesn't
		exists, err := minioClient.BucketExists(r.Context(), bucketName)
		if err != nil {
			utils.ErrorLog("Error checking if bucket exists: %v", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}
		if !exists {
			err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				utils.ErrorLog("Error creating bucket: %v", err.Error())
				utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
				return
			}
		}

		// Check if the file already exists in the bucket
		_, err = minioClient.StatObject(r.Context(), bucketName, fileName, minio.StatObjectOptions{})
		if err == nil {
			// If no error, the file exists, we can choose to overwrite it
			utils.InfoLog("File exists, updating: %s", fileName)
		} else if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			// If the error indicates that the object was not found, we can proceed to upload it
			utils.InfoLog("File does not exist, uploading as new: %s", fileName)
		} else {
			// If there's another error (not a 'file not found' error), log it
			utils.ErrorLog("Error checking file existence: %v", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}

		// Upload the file to MinIO (either new or updated)
		_, err = minioClient.PutObject(r.Context(), bucketName, fileName, file, handler.Size, minio.PutObjectOptions{
			ContentType: handler.Header.Get("Content-Type"),
		})
		if err != nil {
			utils.ErrorLog("Failed to upload file", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}

		fileURL := fmt.Sprintf("%s/%s/%s", domain, bucketName, fileName)
		response := map[string]string{
			"file_name": fileURL,
		}
		utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
	}
}

func UploadHandler(bucketName string, minioClient *minio.Client, domain string, maxFileSize int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.HttpSuccessResponse(w, "Method not allowed", http.StatusMethodNotAllowed, string(utils.ErrorMessage))
			return
		}

		err := r.ParseMultipartForm(maxFileSize) // Use max file size from config
		if err != nil {
			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file")
		if err != nil {
			utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}
		defer file.Close()

		// Generate a new file name using UUID
		newFileName := uuid.New().String() + path.Ext(handler.Filename)

		// Check if the bucket exists, and create it if it doesn't
		exists, err := minioClient.BucketExists(r.Context(), bucketName)
		if err != nil {
			utils.ErrorLog("Error checking if bucket exists: %v", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}
		if !exists {
			err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				utils.ErrorLog("Error creating bucket: %v", err.Error())
				utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
				return
			}
		}

		// Upload the file to MinIO
		_, err = minioClient.PutObject(r.Context(), bucketName, newFileName, file, handler.Size, minio.PutObjectOptions{
			ContentType: handler.Header.Get("Content-Type"),
		})
		if err != nil {

			utils.ErrorLog("Failed to upload file", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}
		fileURL := fmt.Sprintf("%s/%s/%s", domain, bucketName, newFileName)
		response := map[string]string{
			//"bucket":    bucketName,
			"file_name": fileURL,
		}
		utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
	}
}

// DeleteFileHandler deletes a file from the specified MinIO bucket
func DeleteFileHandler(bucketName string, minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.HttpSuccessResponse(w, "Method not allowed", http.StatusMethodNotAllowed, string(utils.ErrorMessage))
			return
		}

		fileName := r.URL.Query().Get("fileName")
		if fileName == "" {
			utils.HttpSuccessResponse(w, "File name is required", http.StatusBadRequest, string(utils.ErrorMessage))
			return
		}

		// Delete the file from MinIO
		err := minioClient.RemoveObject(r.Context(), bucketName, fileName, minio.RemoveObjectOptions{})
		if err != nil {

			utils.ErrorLog("Failed to  delete file", err.Error())
			utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, string(utils.ErrorMessage))
			return
		}

		// Respond with success
		utils.HttpSuccessResponse(w, "", http.StatusOK, string(utils.SuccessMessage))
	}
}

func UploadWithExpriredFilerNameHandler(bucketName string, minioClient *minio.Client, domain string, maxFileSize int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseMultipartForm(maxFileSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName := handler.Filename

		// Check if the bucket exists, and create it if necessary
		exists, err := minioClient.BucketExists(r.Context(), bucketName)
		if err != nil {
			log.Printf("Error checking bucket: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !exists {
			err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Printf("Error creating bucket: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Set the expiration time to 7 days from now
		expirationTime := time.Now().Add(0 * 0.2 * time.Hour).Format(time.RFC3339)

		// Upload file with expiration metadata
		_, err = minioClient.PutObject(r.Context(), bucketName, fileName, file, handler.Size, minio.PutObjectOptions{
			ContentType: handler.Header.Get("Content-Type"),
			UserMetadata: map[string]string{
				"expiration": expirationTime, // Set expiration metadata
			},
		})
		if err != nil {
			log.Printf("Failed to upload file: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the file URL
		fileURL := fmt.Sprintf("%s/%s/%s", domain, bucketName, fileName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"file_url": "%s"}`, fileURL)))
	}
}
