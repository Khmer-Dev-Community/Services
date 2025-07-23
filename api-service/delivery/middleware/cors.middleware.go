package middleware

import (
	"log"
	"net/http"
)

// CORSMiddleware sets CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow the following methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,user-agent")

		// Log CORS handling
		log.Printf("CORS handling for %s %s", r.Method, r.URL.Path)

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			log.Println("Handled preflight request")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
