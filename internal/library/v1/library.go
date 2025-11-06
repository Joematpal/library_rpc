package library_v1

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect"
	library_v1 "github.com/joematpal/library_rpc/pkg/library/v1"
	"github.com/joematpal/library_rpc/pkg/library/v1/library_v1connect"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LibraryService struct {
	mu    sync.RWMutex
	books map[string]*library_v1.Book
	library_v1connect.UnimplementedLibraryServiceHandler
}

// generateID generates a new xid
func generateID() string {
	return xid.New().String()
}

func NewLibraryService() (*LibraryService, error) {
	return &LibraryService{
		books: map[string]*library_v1.Book{},
	}, nil
}

func (s *LibraryService) AddBook(ctx context.Context, req *connect.Request[library_v1.Book]) (*connect.Response[library_v1.Book], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := req.Msg

	// Generate internal ID if not provided
	if book.BookId == "" {
		book.BookId = generateID()
	}

	// Set timestamps
	now := timestamppb.New(time.Now())
	book.CreatedAt = now
	book.UpdatedAt = now

	// Validate required fields
	if book.Title == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("title is required"))
	}
	if book.Isbn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("ISBN is required"))
	}
	if len(book.Authors) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("at least one author is required"))
	}

	// Check if book with this ID already exists
	if _, exists := s.books[book.BookId]; exists {
		return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("book with ID %s already exists", book.BookId))
	}

	// Store the book
	s.books[book.BookId] = book

	return connect.NewResponse(book), nil
}

func (s *LibraryService) GetBook(ctx context.Context, req *connect.Request[library_v1.GetBookRequest]) (*connect.Response[library_v1.Book], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bookID := req.Msg.BookId
	if bookID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("book ID is required"))
	}

	book, exists := s.books[bookID]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("book with ID %s not found", bookID))
	}

	return connect.NewResponse(book), nil
}

func (s *LibraryService) ListBooks(ctx context.Context, req *connect.Request[library_v1.ListBooksRequest]) (*connect.Response[library_v1.BooksList], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]*library_v1.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}

	return connect.NewResponse(&library_v1.BooksList{
		Books: books,
	}), nil
}

func (s *LibraryService) UpdateBook(ctx context.Context, req *connect.Request[library_v1.Book]) (*connect.Response[library_v1.Book], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := req.Msg
	if book.BookId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("book ID is required"))
	}

	// Check if book exists
	existingBook, exists := s.books[book.BookId]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("book with ID %s not found", book.BookId))
	}

	// Validate required fields
	if book.Title == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("title is required"))
	}
	if book.Isbn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("ISBN is required"))
	}
	if len(book.Authors) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("at least one author is required"))
	}

	// Preserve creation time, update modification time
	book.CreatedAt = existingBook.CreatedAt
	book.UpdatedAt = timestamppb.New(time.Now())

	// Update the book
	s.books[book.BookId] = book

	return connect.NewResponse(book), nil
}

func (s *LibraryService) DeleteBook(ctx context.Context, req *connect.Request[library_v1.DeleteBookRequest]) (*connect.Response[library_v1.DeleteBookResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bookID := req.Msg.BookId
	if bookID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("book ID is required"))
	}

	// Check if book exists
	if _, exists := s.books[bookID]; !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("book with ID %s not found", bookID))
	}

	// Delete the book
	delete(s.books, bookID)

	return connect.NewResponse(&library_v1.DeleteBookResponse{
		Success: true,
	}), nil
}
