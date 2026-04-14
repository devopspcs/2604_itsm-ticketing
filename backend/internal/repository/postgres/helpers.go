package postgres

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

func itoa(i int) string {
	return strconv.Itoa(i)
}

func isUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}
