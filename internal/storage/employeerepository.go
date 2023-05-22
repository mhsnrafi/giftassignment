package storage

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db}
}

func (gr *EmployeeRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return gr.db.BeginTx(ctx, nil)
}

func (gr *EmployeeRepository) GetEmployee(ctx context.Context, employeeID string) (Employee, error) {
	var employee Employee
	err := gr.db.QueryRowContext(ctx, "SELECT id, name, interests, gift, returned FROM employees WHERE id = $1", employeeID).Scan(&employee.ID, &employee.Name, &employee.Interests, &employee.Gift, &employee.Returned)
	return employee, err
}

func (gr *EmployeeRepository) UpdateEmployeeGift(ctx context.Context, giftName string, employeeID string) error {
	_, err := gr.db.ExecContext(ctx, "UPDATE employees SET gift = $1, returned = false WHERE id = $2", giftName, employeeID)
	return err
}

func (gr *EmployeeRepository) ValidateEmployeeID(ctx context.Context, employeeID string) (bool, error) {
	// Check the employee ID is not empty
	if employeeID == "" {
		return false, errors.New("employee ID cannot be empty")
	}

	// Check the employee ID exists in the database
	var count int
	err := gr.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM employees WHERE id = $1", employeeID).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, errors.New("employee ID not found in the database")
	}

	// If all checks passed
	return true, nil
}

func (repo *EmployeeRepository) UpdateEmployeeGiftReturnStatus(ctx context.Context, employeeID string, returned bool) error {
	// Create a prepared statement for executing an SQL update
	stmt, err := repo.db.PrepareContext(ctx, "UPDATE employees SET returned = ? WHERE id = ?")
	if err != nil {
		return err
	}

	// Execute the SQL update
	_, err = stmt.ExecContext(ctx, returned, employeeID)
	if err != nil {
		return err
	}

	return nil
}
