package data

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
