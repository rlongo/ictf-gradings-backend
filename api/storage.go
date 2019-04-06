package api

import (
    "database/sql"
)

type DB struct {
    *sql.DB
}

type Tx struct {
    *sql.Tx
}

type StorageService interface {
    StorageServiceBeltTest
}