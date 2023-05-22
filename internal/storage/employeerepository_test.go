package storage

import (
	"context"
	"github.com/spf13/cast"
	"testing"

	"github.com/mhsnrafi/giftassignmentservice/internal/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeRepository_GetEmployee_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	employeeRepository := NewEmployeeRepository(testEnv.DB)

	// Assuming you have a valid employeeID
	employeeID := "1"

	// Act
	employee, err := employeeRepository.GetEmployee(context.Background(), employeeID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, cast.ToInt(employeeID), employee.ID)
}

func TestEmployeeRepository_UpdateEmployeeGift_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	employeeRepository := NewEmployeeRepository(testEnv.DB)

	// Assuming you have a valid employeeID and giftName
	employeeID := "1"
	giftName := "book"

	// Act
	err = employeeRepository.UpdateEmployeeGift(context.Background(), giftName, employeeID)

	// Assert
	assert.NoError(t, err)

	// Check if the update was successful
	employee, err := employeeRepository.GetEmployee(context.Background(), employeeID)
	assert.NoError(t, err)
	assert.Equal(t, giftName, cast.ToString(employee.Gift.String))
}

func TestEmployeeRepository_ValidateEmployeeID_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	employeeRepository := NewEmployeeRepository(testEnv.DB)

	// Assuming you have a valid employeeID
	employeeID := "1"

	// Act
	isValid, err := employeeRepository.ValidateEmployeeID(context.Background(), employeeID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isValid)
}
