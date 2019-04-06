package psql

import (
	_ "github.com/lib/pq"
	"database/sql"
	"strings"
	"io/ioutil"
	"fmt"
	"os"
	"log"
)

var migrationFile = fmt.Sprintf("%s/src/github.com/rlongo/itcf-gradings-backend/storage/psql/schema.sql",
									os.Getenv("GOPATH"))

type PSQLStorageService struct {
    *sql.DB
}

func Open(storageConnectionString string) (*PSQLStorageService, error) {
	db, err := sql.Open("postgres", storageConnectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err!=nil {
		return nil, err
	}

	// Run DB Migration
	log.Printf("Preparing Database")
	if file, err := ioutil.ReadFile(migrationFile); err==nil {
		requests := strings.Split(string(file), ";\n")
		for _, request := range requests {
			if _, err := db.Exec(request); err!= nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}
	
	return &PSQLStorageService{db}, nil
}

func Close(db *PSQLStorageService) error {
	return db.Close()
}