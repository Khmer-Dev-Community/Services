package routers

import (
	"net/http"

	uploader "upload-service/delivery/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7" // Make sure this is v7
)

// SetupUploadRouter initializes and configures the router for upload-related routes
func SetupUploadRouter(r *mux.Router, bucketName string, minioClient *minio.Client, domain string, maxFileSize int64) {
	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	units := api.PathPrefix("/upload").Subrouter()
	units.HandleFunc("/", uploader.HelloAPI).Methods(http.MethodPost, http.MethodGet)
	units.HandleFunc("/pdf", uploader.UploadWithExpriredFilerNameHandler("pdf", minioClient, domain, maxFileSize)).Methods(http.MethodPost)
	units.HandleFunc("/profile", uploader.UploadWithOutChnageFilerNameHandler(bucketName, minioClient, domain, maxFileSize)).Methods(http.MethodPost)
	units.HandleFunc("/employee", uploader.UploadHandler(bucketName, minioClient, domain, maxFileSize)).Methods(http.MethodPost)
	units.HandleFunc("/employee", uploader.DeleteFileHandler(bucketName, minioClient)).Methods(http.MethodDelete)
	// UploadWithOutChnageFilerNameHandler
}
