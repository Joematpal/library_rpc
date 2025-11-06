package library_ctrl

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/joematpal/library_rpc/internal/library_ctrl/flags"
	library_v1 "github.com/joematpal/library_rpc/pkg/library/v1"
	"github.com/joematpal/library_rpc/pkg/library/v1/library_v1connect"
	"github.com/urfave/cli/v3"
)

// NewApp creates a new CLI application
func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "library_ctrl",
		Usage: "Library control CLI",
		Flags: []cli.Flag{
			flags.ServerURL,
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add book operations",
				Commands: []*cli.Command{
					{
						Name:  "book",
						Usage: "Add a new book to the library",
						Flags: []cli.Flag{
							flags.Title,
							flags.Author,
							flags.BookID,
						},
						Action: addBook,
					},
				},
			},
			{
				Name:  "get",
				Usage: "Get operations",
				Commands: []*cli.Command{
					{
						Name:      "book",
						Usage:     "Get a book by ID",
						ArgsUsage: "[book_id]",
						Flags: []cli.Flag{
							flags.BookID,
						},
						Action: getBook,
					},
				},
			},
		},
	}
}

func addBook(ctx context.Context, cmd *cli.Command) error {
	serverURL := cmd.String(flags.Flag_server)
	title := cmd.String(flags.Flag_title)
	author := cmd.String(flags.Flag_author)
	bookID := cmd.String(flags.Flag_book_id)

	client := library_v1connect.NewLibraryServiceClient(http.DefaultClient, serverURL)

	req := connect.NewRequest(&library_v1.AddBookRequest{
		Book: &library_v1.Book{
			BookId: bookID,
			Title:  title,
			Author: author,
		},
	})

	resp, err := client.AddBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add book: %w", err)
	}

	fmt.Printf("Book added successfully with ID: %s\n", resp.Msg.Id)

	return nil
}

func getBook(ctx context.Context, cmd *cli.Command) error {
	serverURL := cmd.String(flags.Flag_server)

	// Get book ID from first argument or flag
	var bookID string
	if cmd.Args().Len() > 0 {
		bookID = cmd.Args().Get(0)
	} else {
		bookID = cmd.String(flags.Flag_book_id)
	}

	if bookID == "" {
		return fmt.Errorf("book ID must be provided either as argument or --book-id flag")
	}

	client := library_v1connect.NewLibraryServiceClient(http.DefaultClient, serverURL)

	req := connect.NewRequest(&library_v1.GetBookRequest{
		Id: bookID,
	})

	resp, err := client.GetBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	fmt.Printf("Book found:\n")
	fmt.Printf("  Book ID: %s\n", resp.Msg.Book.BookId)
	fmt.Printf("  Title: %s\n", resp.Msg.Book.Title)
	fmt.Printf("  Author: %s\n", resp.Msg.Book.Author)

	return nil
}
