package storage

import (
	"database/sql"
	"github.com/lib/pq"
)

type StorageClient struct {
	EmployeeRepo *EmployeeRepository
	GiftRepo     *GiftRepository
}

type Transaction interface {
	Commit() error
	Rollback() error
}

func NewStorageClient(db *sql.DB) StorageClient {
	return StorageClient{
		EmployeeRepo: NewEmployeeRepository(db),
		GiftRepo:     NewGiftRepository(db),
	}
}

type Employee struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Interests pq.StringArray `json:"interests"`
	Gift      sql.NullString `json:"gift,omitempty"`
	Returned  bool           `json:"returned,omitempty"`
}

type Gift struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Categories pq.StringArray `json:"categories"`
	IsTaken    bool           `json:"isTaken,omitempty"`
	Returned   bool           `json:"returned,omitempty"`
}
