package mock

import (
	"fmt"

	"github.com/rlongo/ictf-gradings-backend/api"
)

func (db *MockStorageService) AllBeltTests() (api.BeltTests, error) {
	return db.BeltTestsDB, nil
}

func (db *MockStorageService) GetBeltTest(testID int64) (*api.BeltTest, error) {
	for _, belttest := range db.BeltTestsDB {
		if belttest.ID == testID {
			return belttest, nil
		}
	}

	return nil, fmt.Errorf("BeltTest not found: %d", testID)
}

func (db *MockStorageService) CreateBeltTest(test api.BeltTest) (int64, error) {
	test.ID = (int64)(len(db.BeltTestsDB))
	db.BeltTestsDB = append(db.BeltTestsDB, &test)

	return test.ID, nil
}
