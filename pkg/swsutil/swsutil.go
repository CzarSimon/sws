package swsutil

import (
	"bytes"
	"database/sql"
	"log"
	"os/exec"
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

// runShellCommand executes a command against the os.
func RunShellCommand(main string, args ...string) (string, error) {
	cmd := exec.Command(main, args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errOut.String(), err
	}
	return out.String(), nil
}
