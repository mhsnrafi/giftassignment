package giftassisgmentmanager

import (
	"github.com/mhsnrafi/giftassignmentservice/internal/storage"
)

type Storage interface {
	storage.StorageClient
}

type EmployeeGiftManagementClient struct {
	storageClient storage.StorageClient
	locker        storage.DistributedLocker
}
