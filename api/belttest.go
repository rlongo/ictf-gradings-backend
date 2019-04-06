package api

// The belt test object used by our API
type BeltTest struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Date int64 `json:"date"`
	Location string `json:"location"`
	Admins []string `json:"admins"`
}

type BeltTests []*BeltTest

// Interface to the belt tests
type StorageServiceBeltTest interface {
	AllBeltTests() (BeltTests, error)
	GetBeltTest(int) (*BeltTest, error)
	CreateBeltTest(BeltTest) (int64, error)
}