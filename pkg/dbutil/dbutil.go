package dbutil

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// VerifyDuplicatedConstraintKey checks if the error is a unique constraint violation
//
// It expects the error to be of type *pq.Error and checks for the specific
// PostgreSQL error code for unique constraint violation (23505)
func VerifyDuplicatedConstraintKey(target error) error {
	var pqErr *pq.Error
	// 23505 is the code for unique contraint violation
	if errors.As(target, &pqErr) && pqErr.Code == "23505" {
		// Compile the regular expression to find the pattern "Key (field)="
		re := regexp.MustCompile(`(?i)Key\s*\(\s*(.*?)\s*\)\s*=`)
		matches := re.FindStringSubmatch(pqErr.Detail)
		field := matches[1]

		return fault.NewConflict(fmt.Sprintf("field %s is already taken", field))
	}

	return nil
}

func ExecTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error rolling back tx: %w", err)
		}
		return fmt.Errorf("something went wrong with tx: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting tx: %w", err)
	}

	return nil
}
