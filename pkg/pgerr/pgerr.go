package pgerr

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

// IsUniqueViolation checks if the error is a PostgreSQL unique constraint violation
// and optionally matches the constraint name.
func IsUniqueViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	if pgErr.Code != uniqueViolation {
		return false
	}
	if constraintName != "" {
		return pgErr.ConstraintName == constraintName
	}
	return true
}
