package swsutil

import (
	"database/sql"
	"log"
	"time"
)

// RollbackTx attepts to rollback a transaction.
func RollbackTx(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Println(err)
	}
}

// GetNow gets the current UTC timestamp.
func GetNow() time.Time {
	return time.Now().UTC()
}
