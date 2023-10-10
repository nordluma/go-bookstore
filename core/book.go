package core

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/nordluma/go-bookstore/data"
	"github.com/nordluma/go-bookstore/util"
	"github.com/nordluma/go-bookstore/values"
)

var (
	// Create a new book
	CreateBook = createBook

	// Return a single book
	GetBook = getBook

	// Returns a list of books
	GetAllBooks = getAllBooks

	UpdateBook = updateBook
	DeleteBook = deleteBook

	// Borrows a book if available or returns the book if it's borrowed
	BorrowOrReturnBook = borrowOrReturnBook
)

func createBook(
	ctx context.Context,
	requestBody io.Reader,
) (response interface{}, err error) {
	type createBookRequest struct {
		BookName    string
		AuthorName  string
		Publisher   string
		Description string
	}

	request := &createBookRequest{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		cause := "Failed to decode JSON"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.BookName = strings.TrimSpace(request.BookName)
	if request.BookName == "" {
		cause := "Trying to create a book with empty name"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.AuthorName = strings.TrimSpace(request.AuthorName)
	if request.AuthorName == "" {
		cause := "Trying to create a book with empty author name"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.Publisher = strings.TrimSpace(request.Publisher)
	if request.Publisher == "" {
		cause := "Trying to create a book with empty publisher name"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	response, err = data.CreateBook(
		ctx,
		request.BookName,
		request.AuthorName,
		request.Publisher,
		util.NewNullableString(request.Description),
	)
	if err != nil {
		cause := "Failed to create book"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	return
}

func getBook(
	ctx context.Context,
	bookId string,
) (response interface{}, err error) {
	if bookId == "" {
		cause := "Invalid value for bookId parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	book, err := data.GetBook(ctx, bookId)
	if err != nil {
		cause := "Failed to bget book"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	if book == nil {
		cause := "Book not found"
		err = util.NewError(
			cause,
			util.ErrorCodeEntityNotFound,
			util.ErrResourceNotFound,
			err,
		)
		return
	}

	response = book
	return
}

func getAllBooks(
	ctx context.Context,
	searchTerm string,
	rowOffset, rowLimit, userRole int,
) (response interface{}, err error) {
	if rowOffset < 0 {
		cause := "Invalid value for row offset parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	if rowLimit < 0 || rowLimit > values.MaxRowLimit {
		cause := "Invalid value for row limit parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	if rowLimit == 0 {
		rowLimit = values.MaxRowLimit
	}

	var books interface{}
	if userRole == values.UserRoleMember {
		books, err = data.GetAllBooksForMember(
			ctx,
			searchTerm,
			rowOffset,
			rowLimit,
		)
	} else {
		books, err = data.GetAllBooksForLibrarian(ctx, searchTerm, rowOffset, rowLimit)
	}

	if err != nil {
		cause := "Failed to get all books"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	type metaData struct {
		SearchTerm string `json:",omitempty"`
		RowOffset  int    `json:",omitempty"`
		RowLimit   int
	}

	meta := &metaData{
		SearchTerm: searchTerm,
		RowOffset:  rowOffset,
		RowLimit:   rowLimit,
	}

	type getAllResponse struct {
		Data interface{} `json:"data"`
		Meta interface{} `json:"meta"`
	}

	response = &getAllResponse{
		Data: books,
		Meta: meta,
	}

	return
}

func updateBook(
	ctx context.Context,
	requestBody io.Reader,
) (response interface{}, err error) {
	type updateBookRequest struct {
		BookId      string
		BookName    string
		AuthorName  string
		Publisher   string
		Description string
	}

	request := &updateBookRequest{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		cause := "Failed to decode JSON"
		err = util.NewError(
			cause,
			util.ErrorCodeInvalidJSONBody,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.BookId = strings.TrimSpace(request.BookId)
	if request.BookId == "" {
		cause := "Invalid value for book id parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.BookName = strings.TrimSpace(request.BookName)
	if request.BookName == "" {
		cause := "Invalid value for book name parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.AuthorName = strings.TrimSpace(request.AuthorName)
	if request.AuthorName == "" {
		cause := "Invalid value for author name parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.Publisher = strings.TrimSpace(request.Publisher)
	if request.Publisher == "" {
		cause := "Invalid value for publisher name parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	updatedAt, err := data.UpdateBook(
		ctx,
		request.BookId,
		request.BookName,
		request.AuthorName,
		request.Publisher,
		util.NewNullableString(request.Description),
	)
	if err != nil {
		cause := "Failed to update book"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	if updatedAt == time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) {
		cause := "Book not found"
		err = util.NewError(
			cause,
			util.ErrorCodeEntityNotFound,
			util.ErrResourceNotFound,
			err,
		)
		return
	}

	type updateBookResponse struct {
		UpdatedAt time.Time
	}

	response = &updateBookResponse{
		UpdatedAt: updatedAt,
	}

	return
}

func deleteBook(ctx context.Context, bookId string) (err error) {
	bookId = strings.TrimSpace(bookId)
	if bookId == "" {
		cause := "Invalid value for book id parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	rowsAffected, err := data.DeleteBook(ctx, bookId)
	if err != nil {
		cause := "Failed to delete book"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	if rowsAffected == 0 {
		cause := "Book not found"
		err = util.NewError(
			cause,
			util.ErrorCodeEntityNotFound,
			util.ErrResourceNotFound,
			err,
		)
		return
	}

	return
}

func borrowOrReturnBook(
	ctx context.Context,
	token string,
	requestBody io.Reader,
) (err error) {
	type borrowOrReturnBookRequest struct {
		BookId string
	}

	request := &borrowOrReturnBookRequest{}
	err = json.NewDecoder(requestBody).Decode(request)
	if err != nil {
		cause := "Failed to decode JSON"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	request.BookId = strings.TrimSpace(request.BookId)
	if request.BookId == "" {
		cause := "Invalid value for book id parameter"
		err = util.NewError(
			cause,
			util.ErrorCodeValidation,
			util.ErrBadRequest,
			err,
		)
		return
	}

	status, err := data.GetBookStatus(ctx, request.BookId)
	if err != nil {
		cause := "Failed to get book status"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	if status == values.UserRoleUnknown {
		cause := "Book not found"
		err = util.NewError(
			cause,
			util.ErrorCodeEntityNotFound,
			util.ErrResourceNotFound,
			err,
		)
		return
	}

	userId, err := data.GetUserId(ctx, token)
	if err != nil {
		cause := "Failed to get user id"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	newStatus := values.BookStatusAvailable
	if status == values.BookStatusAvailable {
		newStatus = values.BookStatusBorrowed
	} else {
		borrowerId := ""
		borrowerId, err = data.GetBorrowerId(ctx, request.BookId)
		if err != nil {
			cause := "Failed to get borrower id"
			err = util.NewError(cause, util.ErrorCodeInternal, util.ErrInternal, err)
			return
		}

		if borrowerId != userId {
			cause := "Book not available"
			err = util.NewError(cause, util.ErrorCodeEntityNotFound, util.ErrResourceNotFound, err)
			return
		}

		userId = ""
		newStatus = values.BookStatusAvailable
	}

	err = data.ChangeBookStatus(
		ctx,
		request.BookId,
		newStatus,
		util.NewNullableString(userId),
	)
	if err != nil {
		cause := "Failed to change book status"
		err = util.NewError(
			cause,
			util.ErrorCodeInternal,
			util.ErrInternal,
			err,
		)
		return
	}

	return
}
