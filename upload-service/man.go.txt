package main

import (
	"fmt"
	"io"
	"log"
	"os"

	// Added for completeness, might be used by mknote
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote" // For Nikon/Canon/etc. maker notes
)

func main() {
	// Register all available maker note parsers
	exif.RegisterParsers(mknote.All...)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image-path>")
		return
	}

	imagePath := os.Args[1]
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	lat, lon, err := getGPSFromImage(file)
	if err != nil {
		fmt.Printf("Error getting GPS from image: %v\n", err)
		return
	}

	fmt.Printf("GPS Coordinates: Latitude = %.6f, Longitude = %.6f\n", lat, lon)

	locationInfo := reverseGeocode(lat, lon)
	fmt.Println(locationInfo)
}

// getGPSFromImage extracts GPS latitude and longitude from an image's EXIF data.
func getGPSFromImage(reader io.Reader) (float64, float64, error) {
	// exif.RegisterParsers(mknote.All...) // This should ideally be done once at application startup

	x, err := exif.Decode(reader)
	if err != nil {
		// Specific error handling for the most common decode failure
		if err == io.EOF {
			return 0, 0, fmt.Errorf("failed to decode EXIF data: EOF (empty or truncated file)")
		}
		// Goexif often returns "exif: error reading 4 byte header, got 0, EOF"
		// or similar for fundamentally invalid/empty files. We'll wrap it.
		return 0, 0, fmt.Errorf("failed to decode EXIF data: %w", err)
	}

	lat, long, err := x.LatLong()
	if err != nil {
		// This error occurs if EXIF data is present but no GPS tags are found
		return 0, 0, fmt.Errorf("no GPS data found in EXIF or error getting coordinates: %w", err)
	}

	return lat, long, nil
}

// reverseGeocode is a placeholder for actual reverse geocoding logic.
func reverseGeocode(latitude, longitude float64) string {
	// In a real application, you would call a geocoding API here
	// (e.g., Google Maps Geocoding API, OpenStreetMap Nominatim, etc.)
	return fmt.Sprintf("Reverse geocoding not implemented in this example for %.6f, %.6f", latitude, longitude)
}
