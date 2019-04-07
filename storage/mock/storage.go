package mock

import (
	"github.com/rlongo/ictf-gradings-backend/api"
)

type MockStorageService struct {
	BeltTestsDB api.BeltTests
}

func Open(storageConnectionString string) (*MockStorageService, error) {
	return &MockStorageService{}, nil
}
