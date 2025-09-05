package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"upload-service/config"
	routers "upload-service/delivery/routers"
	"upload-service/utils"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

var whitelist = map[string]bool{
	"/api/upload/profile":                      true,
	"/api/upload/employee":                     true,
	"/api/upload/pdf":                          true,
	"/api/video/videosnap":                     true,
	"/api/video/videostop":                     true,
	"/api/v1/items/list":                       true,
	"/api/swagger/index.html":                  true,
	"/swagger/index.html":                      true,
	"/swagger/swagger-ui-bundle.js":            true,
	"/swagger/swagger-ui.css":                  true,
	"/swagger/swagger-ui-standalone-preset.js": true,
	"/swagger/doc.json":                        true,
	"/swagger/favicon-32x32.png":               true,
	"/swagger/favicon-16x16.png":               true,
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	// Initialize the logger

	cfg, _, err := initConfigAndDatabase("config/config.yml")
	if err != nil {
		logger.Errorf("Initialization error: %v", err)
	}
	err = config.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		utils.ErrorLog(err, "Failed to initialize Redis")
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	utils.InitializeLogger(cfg.Service.LogPtah)
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		logger.Errorf("Initialization error: %v", err)
	}

	r := mux.NewRouter()

	// Create a new CORS handler with desired options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(r)
	routers.SetupUploadRouter(r, "employee", minioClient, cfg.Minio.Domain, cfg.Minio.MaxFileSize)

	go func() {
		for {
			err := deleteExpiredFiles(minioClient, "pdf")
			if err != nil {
				log.Printf("Error deleting expired files: %v", err)
			}
			time.Sleep(1 * time.Minute) // Check every minute
		}
	}()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))

	log.Println("      \n                   *       \n.         *       *A*       *\n.        *A*     **=**     *A*\n        *\"\"\"*   *|\"\"\"|*   *\"\"\"*\n       *|***|*  *|*+*|*  *|***|*\n*********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\nLLLLLLLLLLLLLL//ព្រះសង្ឃ\\\\LLLLLLLLLLLLLLLLL\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n៚ ជោគជ័យ   //  ៚សិរីសួរស្តី \\\\   ៚សុវត្តិភាព \n___________//___៚(♨️)៚__\\\\____________\n៚Application Service is Running Port: 80 ")
	log.Fatal(http.ListenAndServe(":"+cfg.Service.Port, handler))
	select {}
}

func initConfigAndDatabase(configPath string) (config.Config, *gorm.DB, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cfg, nil, err
	}

	db := config.InitDatabase(configPath)
	return cfg, db, nil
}

func deleteExpiredFiles(minioClient *minio.Client, bucketName string) error {
	ctx := context.Background()

	// List objects in the bucket
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Error listing objects: %v", object.Err)
			continue
		}

		// Get object metadata
		objInfo, err := minioClient.StatObject(ctx, bucketName, object.Key, minio.StatObjectOptions{})
		if err != nil {
			log.Printf("Error getting object info: %v", err)
			continue
		}

		// Check the expiration time from metadata
		expirationStr := objInfo.UserMetadata["expiration"]
		if expirationStr != "" {
			expirationTime, err := time.Parse(time.RFC3339, expirationStr)
			if err != nil {
				log.Printf("Error parsing expiration time for object %s: %v", object.Key, err)
				continue
			}

			// Delete the object if it has expired
			if time.Now().After(expirationTime) {
				err := minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
				if err != nil {
					log.Printf("Error deleting object %s: %v", object.Key, err)
					continue
				}
				log.Printf("✅ Deleted expired file: %s", object.Key)
			}
		}
	}

	return nil
}
