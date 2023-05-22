package storage

import (
	"context"
	"testing"

	"github.com/mhsnrafi/giftassignmentservice/internal/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestGiftRepository_GetAllAvailableGifts_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	giftRepository := NewGiftRepository(testEnv.DB)

	// Act
	gifts, err := giftRepository.GetAllAvailableGifts(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, gifts)
}

func TestGiftRepository_UpdateGiftTaken_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	giftRepository := NewGiftRepository(testEnv.DB)

	// Assuming you have a valid giftName
	giftName := "netflix card"

	// Act
	err = giftRepository.UpdateGiftTaken(context.Background(), giftName)

	// Assert
	assert.NoError(t, err)
}

func TestGiftRepository_GetRandomAvailableGift_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	giftRepository := NewGiftRepository(testEnv.DB)

	// Act
	gift, err := giftRepository.GetRandomAvailableGift(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, gift)
}

func TestGiftRepository_UpdateGiftReturnStatus_Success(t *testing.T) {
	// Arrange
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create env: %v", err)
	}
	defer testEnv.Cleanup()

	giftRepository := NewGiftRepository(testEnv.DB)

	// Assuming you have a valid giftName
	giftName := "netflix card"
	returned := true

	// Act
	err = giftRepository.UpdateGiftReturnStatus(context.Background(), giftName, returned)

	// Assert
	assert.NoError(t, err)
}
