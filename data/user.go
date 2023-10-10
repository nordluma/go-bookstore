package data

import "context"

var (
	// Find user with provided username and password and return user's token
	LoginUser = loginUser

	// Return user's role if the user has a token, otherwise return an empty
	// string.
	AuthorizeUser = authorizeUser

	// Return userId from provided token
	GetUserId = getUserId
)

func loginUser(
	ctx context.Context,
	username, password string,
) (response string, err error) {
	query := `
        SELECT token
        FROM library_user
        WHERE username = $1 
        AND user_password = crypt($2, user_password)`

	return executeQueryWithStringResponse(ctx, query, username, password)
}

func authorizeUser(
	ctx context.Context,
	token string,
) (response int64, err error) {
	query := `
        SELECT user_role
        FROM library_user
        WHERE token = $1`

	return executeQueryWithInt64Response(ctx, query, token)
}

func getUserId(
	ctx context.Context,
	token string,
) (response string, err error) {
	query := `
        SELECT user_id
        FROM library_user
        WHERE token = $1`

	return executeQueryWithStringResponse(ctx, query, token)
}
