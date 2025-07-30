package main

import (
	"bytes"
	"strings"
	"testing"

	// REQUIRED for //go:embed directives
	_ "embed"
)

// --- Embedded Test Data ---
// IMPORTANT: These files MUST exist in your 'testdata/' directory
// and match their descriptions for tests to pass correctly.

//  1. A real JPEG with actual GPS data. VERIFY WITH EXIFTOOL.
//     Example: taken directly from a smartphone with location enabled.
//
//go:embed testdata/image_with_gps_actual.jpeg
var imageWithGPSData []byte

//  2. A real JPEG WITHOUT GPS data. VERIFY WITH EXIFTOOL.
//     Example: a screenshot saved as JPEG, or a stripped JPEG.
//
//go:embed testdata/image_without_gps_actual.jpeg
var imageWithoutGPSData []byte

//  3. A text file that is NOT a JPEG, but has some content.
//     Example: a simple text file with "This is not an image." inside.
//
//go:embed testdata/not_a_jpeg_with_content.txt
var notAJPEGData []byte

//  4. An empty file (0 bytes).
//     Create with `touch testdata/empty_file.txt`
//
//go:embed testdata/empty_file.txt
var emptyFileData []byte

// abs helper function (used in tests for float comparison)
// Keep this here if it's not in main.go
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// --- Test Functions ---

func TestGetGPSFromImage(t *testing.T) {
	// IMPORTANT: REPLACE THESE WITH THE *ACTUAL* GPS COORDINATES
	// from your testdata/image_with_gps_actual.jpeg obtained via 'exiftool -n -s -c "%+.6f" -GPSLatitude -GPSLongitude testdata/image_with_gps_actual.jpeg'
	const (
		actualImageLat = 34.052235   // Example: Use your actual GPS Latitude
		actualImageLon = -118.243683 // Example: Use your actual GPS Longitude
	)

	tests := []struct {
		name          string
		imageData     []byte
		expectedLat   float64
		expectedLon   float64
		expectErr     bool
		errorContains string
	}{
		{
			name:          "Image with valid EXIF GPS",
			imageData:     imageWithGPSData,
			expectedLat:   actualImageLat,
			expectedLon:   actualImageLon,
			expectErr:     false,
			errorContains: "", // No error expected if the file truly has GPS
		},
		{
			name:          "Image without EXIF GPS",
			imageData:     imageWithoutGPSData,
			expectedLat:   0,
			expectedLon:   0,
			expectErr:     true,
			errorContains: "no GPS data found", // Correct error for valid JPEG but no GPS tags
		},
		{
			name:          "Corrupted/Invalid Image Data (not a JPEG, has content)",
			imageData:     notAJPEGData,
			expectedLat:   0,
			expectedLon:   0,
			expectErr:     true,
			errorContains: "failed to decode EXIF data", // General decode error expected
		},
		{
			name:        "Empty Image Data (embedded empty file)",
			imageData:   emptyFileData, // Use the truly empty embedded file
			expectedLat: 0,
			expectedLon: 0,
			expectErr:   true,
			// Match the exact error string goexif gives for completely empty input
			errorContains: "failed to decode EXIF data: exif: error reading 4 byte header, got 0, EOF",
		},
		{
			name:        "Empty Byte Slice Directly (runtime empty slice)",
			imageData:   []byte{}, // Directly provide an empty slice
			expectedLat: 0,
			expectedLon: 0,
			expectErr:   true,
			// Match the exact error string goexif gives for completely empty input
			errorContains: "failed to decode EXIF data: exif: error reading 4 byte header, got 0, EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.imageData)
			lat, lon, err := getGPSFromImage(reader) // Calls the function from main.go

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected an error, but got none for %s", tt.name)
				}
				// Check if the actual error message contains the expected substring
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("for %s, expected error to contain '%s', but got '%v'", tt.name, tt.errorContains, err)
				}
			} else { // Expect no error
				if err != nil {
					t.Errorf("did not expect an error, but got: %v for %s", err, tt.name)
				}
				const epsilon = 0.000001
				if !(abs(lat-tt.expectedLat) < epsilon && abs(lon-tt.expectedLon) < epsilon) {
					t.Errorf("for %s, expected (lat, lon) = (%.6f, %.6f), got (%.6f, %.6f)",
						tt.name, tt.expectedLat, tt.expectedLon, lat, lon)
				}
			}
		})
	}
}

func TestReverseGeocode(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected string
	}{
		{
			name:     "Known Coordinates",
			lat:      34.052235,
			lon:      -118.243683,
			expected: "Reverse geocoding not implemented in this example for 34.052235, -118.243683",
		},
		{
			name:     "Zero Coordinates",
			lat:      0.0,
			lon:      0.0,
			expected: "Reverse geocoding not implemented in this example for 0.000000, 0.000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reverseGeocode(tt.lat, tt.lon) // Calls the function from main.go
			if result != tt.expected {
				t.Errorf("reverseGeocode(%.6f, %.6f) = '%s'; want '%s'", tt.lat, tt.lon, result, tt.expected)
			}
		})
	}
}
