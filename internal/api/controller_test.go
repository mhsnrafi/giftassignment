package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mhsnrafi/giftassignmentservice/internal/api"
	cust_error "github.com/mhsnrafi/giftassignmentservice/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGiftAssignmentManager struct {
	mock.Mock
}

func (m *MockGiftAssignmentManager) ReturnGift(ctx context.Context, employeeID string) error {
	args := m.Called(ctx, employeeID)
	return args.Error(0)
}

func (m *MockGiftAssignmentManager) AssignGifts(ctx context.Context, employeeID string) (string, error) {
	args := m.Called(ctx, employeeID)
	return args.String(0), args.Error(1)
}

func TestAssignGifts(t *testing.T) {
	testCases := []struct {
		name             string
		employeeID       string
		mockGiftResponse string
		mockError        error
		expectedResponse string
		expectedStatus   int
	}{
		{
			name:             "Valid request",
			employeeID:       "employee-1",
			mockGiftResponse: "gift-1",
			mockError:        nil,
			expectedResponse: `{"message":"Gift successfully assigned","gift":"gift-1"}`,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Employee already has a gift assigned",
			employeeID:       "employee-1",
			mockGiftResponse: "",
			mockError:        cust_error.ErrGiftAlreadyAssigned,
			expectedResponse: `{"error":"Conflict","message":"gift is already assigned to this employee"}`,
			expectedStatus:   http.StatusConflict,
		},
		{
			name:             "Server error on gift assignment",
			employeeID:       "employee-1",
			mockGiftResponse: "",
			mockError:        errors.New("Server error"),
			expectedResponse: `{"error":"Internal Server Error","message":"Internal Server Error: failed to assign gift"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:             "Request with empty employeeID",
			employeeID:       "",
			mockGiftResponse: "",
			mockError:        nil,
			expectedResponse: `{"error":"Bad Request","message":"Invalid Employee ID: ID cannot be empty"}`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "AssignGifts successful but returns an empty gift",
			employeeID:       "employee-2",
			mockGiftResponse: "",
			mockError:        nil,
			expectedResponse: `{"message":"Gift successfully assigned","gift":""}`,
			expectedStatus:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockManager := new(MockGiftAssignmentManager)
			mockManager.On("AssignGifts", mock.Anything, tc.employeeID).Return(tc.mockGiftResponse, tc.mockError)

			controller := api.NewController(mockManager)
			router := mux.NewRouter()
			router.HandleFunc("/assign-gift/{id}", controller.AssignGifts)

			req, err := http.NewRequest(http.MethodPost, "/assign-gift/"+tc.employeeID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code == http.StatusNotFound {
				assert.Equal(t, http.StatusNotFound, rr.Code)
				return
			}

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.JSONEq(t, tc.expectedResponse, rr.Body.String())
		})
	}
}

func TestReturnGift(t *testing.T) {
	testCases := []struct {
		name             string
		employeeID       string
		mockError        error
		expectedResponse string
		expectedStatus   int
	}{
		{
			name:             "Valid request",
			employeeID:       "employee-1",
			mockError:        nil,
			expectedResponse: `{"message":"Gift successfully returned"}`,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Gift is not returnable",
			employeeID:       "employee-1",
			mockError:        cust_error.ErrGiftNotReturnable,
			expectedResponse: `{"error":"Conflict","message":"gift is not returnable"}`,
			expectedStatus:   http.StatusConflict,
		},
		{
			name:             "Server error on gift return",
			employeeID:       "employee-1",
			mockError:        errors.New("Server error"),
			expectedResponse: `{"error":"Internal Server Error","message":"Internal Server Error: failed to return gift"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockManager := new(MockGiftAssignmentManager)
			mockManager.On("ReturnGift", mock.Anything, tc.employeeID).Return(tc.mockError)

			controller := api.NewController(mockManager)
			router := mux.NewRouter()
			router.HandleFunc("/return-gift/{id}", controller.ReturnGift)

			req, err := http.NewRequest(http.MethodPost, "/return-gift/"+tc.employeeID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.JSONEq(t, tc.expectedResponse, rr.Body.String())
		})
	}
}
