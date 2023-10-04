package values

const (
	UserRoleUnknown   = 0
	UserRoleMember    = 1
	UserRoleLibrarian = 2

	// The amount of rows that can be fetched from database
	MaxRowLimit = 50

	BookStatusUnknown   = 0
	BookStatusAvailable = 1
	BookStatusBorrowed  = 2
)
