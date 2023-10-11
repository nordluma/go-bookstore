package database

import "database/sql"

var (
	// Return a null string if the parameter is empty,
	// otherwise returns a valid string.
	NewNullableString = newNullableString

	// Return empty string if the value is null,
	// oterwise returns the value of the string.
	GetNullStringValue = getNullStringValue
)

// NullString represents nullable string.
// This is used when fetching data from SQL queries
type NullString struct {
	sql.NullString
}

func newNullableString(x string) NullString {
	if x == "" {
		return NullString{}
	}

	return NullString{sql.NullString{String: x, Valid: true}}
}

func getNullStringValue(x NullString) string {
	if !x.Valid {
		return ""
	}

	return x.String
}
