package psql

import (
	"strings"

	"github.com/rlongo/ictf-gradings-backend/api"
)

const AdminDelimiter = "|"

func (db *PSQLStorageService) AllBeltTests() (api.BeltTests, error) {
	query := "SELECT id, test_name, test_date, dojang, admins FROM belt_tests"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	belttests := make(api.BeltTests, 0)
	for rows.Next() {
		t := new(api.BeltTest)
		var admins string

		err := rows.Scan(&t.ID, &t.Name, &t.Date, &t.Location, &admins)
		if err != nil {
			return nil, err
		}

		t.Admins = strings.Split(admins, AdminDelimiter)

		belttests = append(belttests, t)
	}

	return belttests, nil
}

func (db *PSQLStorageService) GetBeltTest(testID int) (*api.BeltTest, error) {
	t := new(api.BeltTest)
	var admins string

	query := "SELECT id, test_name, test_date, dojang, admins FROM belt_tests WHERE id=$1"
	row := db.QueryRow(query, testID)

	err := row.Scan(&t.ID, &t.Name, &t.Date, &t.Location, &admins)
	if err != nil {
		return nil, err
	}

	t.Admins = strings.Split(admins, AdminDelimiter)

	return t, nil
}

func (db *PSQLStorageService) CreateBeltTest(test api.BeltTest) (int64, error) {
	var id int64 = 0

	query := `
		INSERT INTO belt_tests (test_name, test_date, dojang, admins)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	admins := strings.Join(test.Admins, AdminDelimiter)

	err := db.QueryRow(query, test.Name, test.Date, test.Location, admins).Scan(&id)

	return id, err
}
