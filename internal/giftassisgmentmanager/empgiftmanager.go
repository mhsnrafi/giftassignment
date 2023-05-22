package giftassisgmentmanager

import (
	"context"
	custom_error "github.com/mhsnrafi/giftassignmentservice/internal/errors"
	"github.com/mhsnrafi/giftassignmentservice/internal/storage"
	"github.com/spf13/cast"
	"log"
	"time"
)

type Operation func(ctx context.Context, tx storage.Transaction) (string, error)

func NewEmployeeGiftManagementClient(storage storage.StorageClient, locker storage.DistributedLocker) *EmployeeGiftManagementClient {
	return &EmployeeGiftManagementClient{
		storageClient: storage,
		locker:        locker,
	}
}

func (egm *EmployeeGiftManagementClient) AssignGifts(ctx context.Context, employeeID string) (string, error) {
	operation := func(ctx context.Context, tx storage.Transaction) (string, error) {
		// Check if employee exists and doesn't have a gift yet
		employee, err := egm.getValidEmployeeWithoutGift(ctx, employeeID)
		if err != nil {
			return "", err
		}

		// Try to find a gift that matches the employee's interests
		giftName, err := egm.assignGiftBasedOnInterest(ctx, employee)
		if err != nil {
			return "", err
		}

		// If a gift was assigned based on interests, return it
		if giftName != "" {
			return giftName, nil
		}

		// If no gift was found based on interests, assign a random gift
		return egm.assignRandomGift(ctx, employeeID)
	}

	return egm.withLockAndTransaction(ctx, employeeID, operation)
}

func (egm *EmployeeGiftManagementClient) ReturnGift(ctx context.Context, employeeID string) error {
	operation := func(ctx context.Context, tx storage.Transaction) (string, error) {
		// Verify if employee is valid and has a gift to return
		employee, err := egm.getValidEmployeeWithGift(ctx, employeeID)
		if err != nil {
			return "", err
		}

		// Process gift return
		err = egm.processGiftReturn(ctx, employeeID, employee.Gift.String)
		if err != nil {
			return "", err
		}

		// Return empty string as operation function expects a string and an error
		return "", nil
	}

	_, err := egm.withLockAndTransaction(ctx, employeeID, operation)
	return err
}

// getValidEmployeeWithoutGift checks if the employee is valid and does not have a gift already.
func (egm *EmployeeGiftManagementClient) getValidEmployeeWithoutGift(ctx context.Context, employeeID string) (*storage.Employee, error) {
	// Verify that the employeeID is valid
	valid, err := egm.storageClient.EmployeeRepo.ValidateEmployeeID(ctx, employeeID)
	if err != nil || !valid {
		log.Printf("Invalid Employee ID: %s", employeeID)
		return nil, custom_error.ErrInvalidEmployeeID
	}

	// Get the employee's interests and check if they already have a gift
	employee, err := egm.storageClient.EmployeeRepo.GetEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	// If the employee already has a gift that hasn't been returned, return an error
	if employee.Gift.Valid && employee.Gift.String != "" && !employee.Returned {
		return nil, custom_error.ErrGiftAlreadyAssigned
	}

	return &employee, nil
}

// getValidEmployeeWithGift checks if the employee is valid and has a gift to return.
func (egm *EmployeeGiftManagementClient) getValidEmployeeWithGift(ctx context.Context, employeeID string) (*storage.Employee, error) {
	// Verify that the employeeID is valid
	valid, err := egm.storageClient.EmployeeRepo.ValidateEmployeeID(ctx, employeeID)
	if err != nil || !valid {
		log.Printf("Invalid Employee ID: %s", employeeID)
		return nil, custom_error.ErrInvalidEmployeeID
	}

	// Get the employee's gift
	employee, err := egm.storageClient.EmployeeRepo.GetEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	// If the employee doesn't have a gift or has already returned a gift, return an error
	if !employee.Gift.Valid || employee.Gift.String == "" || employee.Returned {
		return nil, custom_error.ErrGiftNotReturnable
	}

	return &employee, nil
}

// assignGiftBasedOnInterest tries to assign a gift to the employee based on their interests.
func (egm *EmployeeGiftManagementClient) assignGiftBasedOnInterest(ctx context.Context, employee *storage.Employee) (string, error) {
	gifts, err := egm.storageClient.GiftRepo.GetAllAvailableGifts(ctx)
	if err != nil {
		return "", err
	}

	// Iterate over each gift
	for _, gift := range gifts {
		for _, interest := range employee.Interests {
			for _, category := range gift.Categories {
				if interest == category {
					// If a match is found, update the gift and employee records within the transaction
					err = egm.storageClient.GiftRepo.UpdateGiftTaken(ctx, gift.Name)
					if err != nil {
						return "", err
					}

					err = egm.storageClient.EmployeeRepo.UpdateEmployeeGift(ctx, gift.Name, cast.ToString(employee.ID))
					if err != nil {
						return "", err
					}

					// If everything went fine, return the gift name
					return gift.Name, nil
				}
			}
		}
	}

	return "", nil
}

// assignRandomGift assigns a random gift to the employee.
func (egm *EmployeeGiftManagementClient) assignRandomGift(ctx context.Context, employeeID string) (string, error) {
	// If no suitable gift is found, assign a random unassigned gift
	randomGift, err := egm.storageClient.GiftRepo.GetRandomAvailableGift(ctx)
	if err != nil {
		return "", err
	}

	err = egm.storageClient.GiftRepo.UpdateGiftTaken(ctx, randomGift.Name)
	if err != nil {
		return "", err
	}

	err = egm.storageClient.EmployeeRepo.UpdateEmployeeGift(ctx, randomGift.Name, employeeID)
	if err != nil {
		return "", err
	}

	return randomGift.Name, nil
}

// processGiftReturn processes the return of a gift by an employee.
func (egm *EmployeeGiftManagementClient) processGiftReturn(ctx context.Context, employeeID string, giftName string) error {
	// Mark the gift as returned in both the gift and employee records
	err := egm.storageClient.GiftRepo.UpdateGiftReturnStatus(ctx, giftName, true)
	if err != nil {
		return err
	}

	err = egm.storageClient.EmployeeRepo.UpdateEmployeeGiftReturnStatus(ctx, employeeID, true)
	if err != nil {
		return err
	}

	return nil
}

func (egm *EmployeeGiftManagementClient) withLockAndTransaction(ctx context.Context, employeeID string, operation Operation) (string, error) {
	var result string
	var err error
	retryCount := 0

	for {
		// Try to acquire the lock
		err = egm.locker.Lock(ctx, employeeID, 10*time.Second)
		if err != nil {
			// If lock could not be acquired and we haven't hit the max retry count, backoff and retry
			if retryCount < 5 { // Assuming max 5 retries
				retryCount++
				log.Printf("Failed to acquire lock for employeeID: %s. Retrying...", employeeID)
				time.Sleep(time.Duration(retryCount) * time.Second) // Increasing backoff
				continue
			} else {
				log.Printf("Failed to acquire lock for employeeID: %s", employeeID)
				return "", err
			}
		}

		// If lock was acquired, break the loop
		break
	}

	// Ensure the lock is released
	defer func() {
		if err := egm.locker.Unlock(ctx, employeeID); err != nil {
			log.Printf("Failed to release lock for employeeID: %s", employeeID)
		}
	}()

	// Start the transaction
	tx, err := egm.storageClient.EmployeeRepo.BeginTransaction(ctx)
	if err != nil {
		log.Printf("Failed to start transaction for employeeID: %s", employeeID)
		return "", err
	}

	// Ensure the transaction is rolled back in the case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction for employeeID: %s", employeeID)
			}
		}
	}()

	// Execute the operation
	result, err = operation(ctx, tx)
	if err != nil {
		log.Printf("Operation failed for employeeID: %s", employeeID)
		return "", err
	}

	// If the operation succeeded, commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction for employeeID: %s", employeeID)
		return "", err
	}

	// If everything is successful, return the result
	return result, nil
}
