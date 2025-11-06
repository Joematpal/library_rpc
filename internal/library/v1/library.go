package library_v1

import (
	"context"
	"errors"
	"sync"

	"connectrpc.com/connect"
	library_v1 "github.com/joematpal/library_rpc/pkg/library/v1"
	"github.com/joematpal/library_rpc/pkg/library/v1/library_v1connect"
)

type LibraryService struct {
	mu    sync.RWMutex
	books map[string]*library_v1.Book
	library_v1connect.UnimplementedLibraryServiceHandler
}

func NewLibraryService() (*LibraryService, error) {
	return &LibraryService{
		books: map[string]*library_v1.Book{},
	}, nil
}

func (s *LibraryService) AddBook(context.Context, *connect.Request[library_v1.AddBookRequest]) (*connect.Response[library_v1.AddBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("library.v1.LibraryService.AddBook is not implemented"))
}

func (s *LibraryService) GetBook(context.Context, *connect.Request[library_v1.GetBookRequest]) (*connect.Response[library_v1.GetBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("library.v1.LibraryService.GetBook is not implemented"))
}
