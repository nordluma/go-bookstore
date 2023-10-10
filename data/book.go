package data

import (
	"context"
	"time"

	"github.com/nordluma/go-bookstore/server/dbserver"
	"github.com/nordluma/go-bookstore/util"
	"github.com/nordluma/go-bookstore/values"
)

var (
	// Create a new book
	CreateBook = createBook

	// Retrieve a book
	GetBook = getBook

	// Return a list of books for members
	GetAllBooksForMember = getAllBooksForMember

	// Return a list of books for librarians
	GetAllBooksForLibrarian = getAllBooksForLibrarian

	// Update a book
	UpdateBook = updateBook

	// Delete a book
	DeleteBook = deleteBook

	// Get status of book
	GetBookStatus = getBookStatus

	// Get borrower id of member who has borrowed the book
	GetBorrowerId = getBorrowerId

	// Change the book status
	ChangeBookStatus = changeBookStatus
)

// This struct contains all database columns converted to Go types
type BookEntity struct {
	BookId      string
	BookName    string
	AuthorName  string
	Publisher   string
	Description string `json:",omitempty"`
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	BorrowerId  string `json:",omitempty"`
}

// Struct to decribe a book. This struct is used when quering a single book
type BookDetails struct {
	BookId      string
	BookName    string
	AuthorName  string
	Publisher   string
	Description string `json:",omitempty"`
}

// Struct which is used when librarians queries for all books
type BookInfoLibrarian struct {
	BookId     string
	BookName   string
	AuthorName string
	Publisher  string
	Status     int64
	Borrower   string `json:",omitempty"`
}

// Struct which is used when members queries for all books
type BookInfoMember struct {
	BookId     string
	BookName   string
	AuthorName string
	Publisher  string
}

func createBook(
	ctx context.Context,
	bookName, authorName, publisher string,
	description util.NullString,
) (response *BookEntity, err error) {
	dbRunner := ctx.Value(values.ContextKeyDbRunner).(dbserver.Runner)

	query := `
        INSERT INTO book (
            book_name, author_name, publisher, book_description
        )
        VALUES ($1, $2, $3, $4)
        RETURNING book_id, created_at`

	rows, err := dbRunner.Query(
		ctx,
		query,
		bookName,
		authorName,
		publisher,
		description,
	)
	if err != nil {
		return
	}

	defer rows.Close()

	rr, err := dbserver.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		response = &BookEntity{
			BookId:      rr.ReadByIdxString(0),
			BookName:    bookName,
			AuthorName:  authorName,
			Publisher:   publisher,
			Description: util.GetNullStringValue(description),
			Status:      values.BookStatusAvailable,
			CreatedAt:   rr.ReadByIdxTime(1),
			UpdatedAt:   rr.ReadByIdxTime(1),
			BorrowerId:  "",
		}
	}

	err = rr.Error()

	return
}

func getBook(
	ctx context.Context,
	bookId string,
) (response *BookDetails, err error) {
	dbRunner := ctx.Value(values.ContextKeyDbRunner).(dbserver.Runner)

	query := `
        SELECT
            book_id AS "BookId",
            book_name AS "BookName",
            author_name AS "AuthorName",
            publisher AS "Publisher",
            book_description AS "Description"
        FROM book
        WHERE book_id = $1`

	rows, err := dbRunner.Query(ctx, query, bookId)
	if err != nil {
		return
	}

	defer rows.Close()

	rr, err := dbserver.GetRowReader(rows)
	if err != nil {
		return
	}

	if rr.ScanNext() {
		response = &BookDetails{}
		rr.ReadAllToStruct(response)
	}

	err = rr.Error()

	return
}

func getAllBooksForMember(
	ctx context.Context,
	searchTerm string,
	rowOffset, rowLimit int,
) (response []*BookInfoMember, err error) {
	dbRunner := ctx.Value(values.ContextKeyDbRunner).(dbserver.Runner)

	query := `
        SELECT
            book_id AS "BookId",
            book_name AS "BookName",
            author_name AS "AuthorName",
            publisher AS "Publisher"
        FROM book
        WHERE book_name LIKE '%%' || $1 || '%%'
        AND book_status = $2
        OFFSET $3
        LIMIT $4`

	rows, err := dbRunner.Query(
		ctx,
		query,
		searchTerm,
		values.BookStatusAvailable,
		rowOffset,
		rowLimit,
	)
	if err != nil {
		return
	}

	defer rows.Close()

	rr, err := dbserver.GetRowReader(rows)
	if err != nil {
		return
	}

	response = make([]*BookInfoMember, 0)
	for rr.ScanNext() {
		book := &BookInfoMember{}
		rr.ReadAllToStruct(book)
		response = append(response, book)
	}

	err = rr.Error()

	return
}

func getAllBooksForLibrarian(
	ctx context.Context,
	searchTerm string,
	rowOffset, rowLimit int,
) (response []*BookInfoLibrarian, err error) {
	dbRunner := ctx.Value(values.ContextKeyDbRunner).(dbserver.Runner)

	query := `
        SELECT
            b.book_id AS "BookId",
            b.book_name AS "BookName",
            b.author_name AS "AuthorName",
            b.publisher AS "Publisher",
            b.book_status AS "Status",
            u.full_name AS "Borrower"
        FROM book b
        LEFT JOIN library_user u ON u.user_id = b.borrower_id
        WHERE b.book_name LIKE '%%' || $1 || '%%'
        OFFSET $2
        LIMIT $3`

	rows, err := dbRunner.Query(ctx, query, searchTerm, rowOffset, rowLimit)
	if err != nil {
		return
	}

	defer rows.Close()

	rr, err := dbserver.GetRowReader(rows)
	if err != nil {
		return
	}

	response = make([]*BookInfoLibrarian, 0)
	for rr.ScanNext() {
		book := &BookInfoLibrarian{}
		rr.ReadAllToStruct(book)
		response = append(response, book)
	}

	err = rr.Error()

	return
}

func updateBook(
	ctx context.Context,
	bookId, bookName, authorName, publisher string,
	description util.NullString,
) (response time.Time, err error) {
	query := `
        UPDATE book
        SET
            book_name = $1,
            author_name = $2,
            publisher = $3,
            book_description = $4
        WHERE book_id = $5
        RETURNING updated_at`

	return executeQueryWithTimeResponse(
		ctx,
		query,
		bookName,
		bookId,
		authorName,
		publisher,
		description,
		bookId,
	)
}

func deleteBook(
	ctx context.Context,
	bookId string,
) (response int64, err error) {
	query := `DELETE FROM book WHERE book_id = $1`

	return executeQueryWithRowsAffected(ctx, query, bookId)
}

func getBookStatus(
	ctx context.Context,
	bookId string,
) (response int64, err error) {
	query := `SELECT book_status FROM book WHERE book_id = $1`

	return executeQueryWithInt64Response(ctx, query, bookId)
}

func changeBookStatus(
	ctx context.Context,
	bookId string,
	status int,
	userId util.NullString,
) (err error) {
	dbRunner := ctx.Value(values.ContextKeyDbRunner).(dbserver.Runner)

	query := `
        UPDATE book
        SET
            book_status = $1,
            borrower_id = $2,
        WHERE book_id = $3`

	_, err = dbRunner.Exec(ctx, query, status, userId, bookId)
	return
}

func getBorrowerId(
	ctx context.Context,
	bookId string,
) (response string, err error) {
	query := `SELECT borrower_id FROM book WHERE book_id = $1`

	return executeQueryWithStringResponse(ctx, query, bookId)
}
