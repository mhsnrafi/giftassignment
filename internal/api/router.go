package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

const (
	assignGift = "/assign-gift/{id}"
	returnGift = "/employees/{employeeId}/return-gift"
)

var limiter = rate.NewLimiter(10, 100)

// limitMiddleware is a middleware that limits the number of requests per second
// to 10 requests per second
// If the limit is exceeded, a 429 Too Many Requests response is returned
func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewAPI returns a new API router
// The router is configured with the API controller
// and the rate limiting middleware
func NewAPI(apiController Controller) http.Handler {
	router := mux.NewRouter()

	// Add rate limiting middleware to all endpoints
	router.Use(limitMiddleware)

	router.HandleFunc(assignGift, apiController.AssignGifts).Methods(http.MethodPost)
	router.HandleFunc(returnGift, apiController.AssignGifts).Methods(http.MethodPost)

	return router
}
