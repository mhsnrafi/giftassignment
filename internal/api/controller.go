package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mhsnrafi/giftassignmentservice/internal/errors"
)

// GiftAssignmentManager is the interface for the transaction manager
type GiftAssignmentManager interface {
	AssignGifts(ctx context.Context, employeeID string) (string, error)
	ReturnGift(ctx context.Context, employeeID string) error
}

// Controller is the API controller
type Controller struct {
	GiftAssignmentManager GiftAssignmentManager
}

// NewController constructs a new controller with the given GiftAssignmentManager
func NewController(gam GiftAssignmentManager) Controller {
	return Controller{
		GiftAssignmentManager: gam,
	}
}

// AssignGifts handles the HTTP request for assigning gifts to an employee
func (c *Controller) AssignGifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	employeeID := vars["id"]

	// Attempt to assign a gift using the GiftAssignmentManager
	gift, err := c.GiftAssignmentManager.AssignGifts(ctx, employeeID)
	if err != nil {
		// If the employee already has a gift, return a conflict status
		if err == errors.ErrGiftAlreadyAssigned {
			httpError(w, err.Error(), http.StatusConflict)
			return
		}
		// For all other errors, return an internal server error
		httpError(w, "Internal Server Error: failed to assign gift", http.StatusInternalServerError)
		return
	}

	// If the assignment was successful, return a success response
	response := struct {
		Message string `json:"message"`
		Gift    string `json:"gift"`
	}{
		Message: "Gift successfully assigned",
		Gift:    gift,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// ReturnGift processes the HTTP request to facilitate a gift return for an employee.
func (c *Controller) ReturnGift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	employeeID := vars["id"]

	// Attempt to return a gift using the GiftAssignmentManager
	err := c.GiftAssignmentManager.ReturnGift(ctx, employeeID)
	if err != nil {
		// If the gift is not returnable, return a conflict status
		if err == errors.ErrGiftNotReturnable {
			httpError(w, err.Error(), http.StatusConflict)
			return
		}
		// For all other errors, return an internal server error
		httpError(w, "Internal Server Error: failed to return gift", http.StatusInternalServerError)
		return
	}

	// If the return was successful, return a success response
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Gift successfully returned",
	}
	respondWithJSON(w, http.StatusOK, response)
}

// decodeJSON decodes JSON from a HTTP request into the given value
func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// respondWithJSON writes a JSON response with the given status code and data
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// httpError sends a HTTP error response with the given status code and message
func httpError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
