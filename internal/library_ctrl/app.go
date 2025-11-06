package library_ctrl

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
							flags.Authors,
							flags.ISBN,
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
			{
				Name:  "list",
				Usage: "List operations",
				Commands: []*cli.Command{
					{
						Name:   "books",
						Usage:  "List all books",
						Action: listBooks,
					},
				},
			},
			{
				Name:  "update",
				Usage: "Update operations",
				Commands: []*cli.Command{
					{
						Name:  "book",
						Usage: "Update an existing book",
						Flags: []cli.Flag{
							flags.BookID,
							flags.Title,
							flags.Authors,
							flags.ISBN,
						},
						Action: updateBook,
					},
				},
			},
			{
				Name:  "delete",
				Usage: "Delete operations",
				Commands: []*cli.Command{
					{
						Name:      "book",
						Usage:     "Delete a book by ID",
						ArgsUsage: "[book_id]",
						Flags: []cli.Flag{
							flags.BookID,
						},
						Action: deleteBook,
					},
				},
			},
		},
	}
}

func addBook(ctx context.Context, cmd *cli.Command) error {
	serverURL := cmd.String(flags.Flag_server)
	title := cmd.String(flags.Flag_title)
	authorsSlice := cmd.StringSlice(flags.Flag_authors)
	isbn := cmd.String(flags.Flag_isbn)

	client := library_v1connect.NewLibraryServiceClient(http.DefaultClient, serverURL)

	// Parse authors from "given_name family_name" format
	var authors []*library_v1.Author
	for _, authorStr := range authorsSlice {
		parts := strings.Fields(authorStr)
		if len(parts) < 2 {
			return fmt.Errorf("author must be in format 'given_name family_name', got: %s", authorStr)
		}
		authors = append(authors, &library_v1.Author{
			GivenName:  parts[0],
			FamilyName: strings.Join(parts[1:], " "), // Handle multiple family names
		})
	}

	req := connect.NewRequest(&library_v1.Book{
		Isbn:    isbn,
		Title:   title,
		Authors: authors,
	})

	resp, err := client.AddBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add book: %w", err)
	}

	fmt.Printf("Book added successfully with ID: %s\n", resp.Msg.BookId)

	return nil
}

func listBooks(ctx context.Context, cmd *cli.Command) error {
	serverURL := cmd.String(flags.Flag_server)
	client := library_v1connect.NewLibraryServiceClient(http.DefaultClient, serverURL)

	req := connect.NewRequest(&library_v1.ListBooksRequest{})

	resp, err := client.ListBooks(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list books: %w", err)
	}

	if len(resp.Msg.Books) == 0 {
		fmt.Println("No books found.")
		return nil
	}

	fmt.Printf("Found %d books:\n", len(resp.Msg.Books))
	for i, book := range resp.Msg.Books {
		fmt.Printf("\n%d. Book ID: %s\n", i+1, book.BookId)
		fmt.Printf("   ISBN: %s\n", book.Isbn)
		fmt.Printf("   Title: %s\n", book.Title)
		for j, author := range book.Authors {
			fmt.Printf("   Author %d: %s %s\n", j+1, author.GivenName, author.FamilyName)
		}
	}

	return nil
}

func updateBook(ctx context.Context, cmd *cli.Command) error {
	serverURL := cmd.String(flags.Flag_server)
	bookID := cmd.String(flags.Flag_book_id)
	title := cmd.String(flags.Flag_title)
	authorsSlice := cmd.StringSlice(flags.Flag_authors)
	isbn := cmd.String(flags.Flag_isbn)

	if bookID == "" {
		return fmt.Errorf("book ID is required")
	}

	client := library_v1connect.NewLibraryServiceClient(http.DefaultClient, serverURL)

	// Parse authors from "given_name family_name" format
	var authors []*library_v1.Author
	for _, authorStr := range authorsSlice {
		parts := strings.Fields(authorStr)
		if len(parts) < 2 {
			return fmt.Errorf("author must be in format 'given_name family_name', got: %s", authorStr)
		}
		authors = append(authors, &library_v1.Author{
			GivenName:  parts[0],
			FamilyName: strings.Join(parts[1:], " "), // Handle multiple family names
		})
	}

	req := connect.NewRequest(&library_v1.Book{
		BookId:  bookID,
		Isbn:    isbn,
		Title:   title,
		Authors: authors,
	})

	resp, err := client.UpdateBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	fmt.Printf("Book updated successfully:\n")
	fmt.Printf("  Book ID: %s\n", resp.Msg.BookId)
	fmt.Printf("  ISBN: %s\n", resp.Msg.Isbn)
	fmt.Printf("  Title: %s\n", resp.Msg.Title)
	for _, author := range resp.Msg.Authors {
		fmt.Printf("  Author: %s %s\n", author.GivenName, author.FamilyName)
	}

	return nil
}

func deleteBook(ctx context.Context, cmd *cli.Command) error {
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

	req := connect.NewRequest(&library_v1.DeleteBookRequest{
		BookId: bookID,
	})

	resp, err := client.DeleteBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	if resp.Msg.Success {
		fmt.Printf("Book with ID %s deleted successfully.\n", bookID)
	} else {
		fmt.Printf("Failed to delete book with ID %s.\n", bookID)
	}

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
		BookId: bookID,
	})

	resp, err := client.GetBook(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	fmt.Printf("Book found:\n")
	fmt.Printf("  Book ID: %s\n", resp.Msg.BookId)
	fmt.Printf("  ISBN: %s\n", resp.Msg.Isbn)
	fmt.Printf("  Title: %s\n", resp.Msg.Title)

	for _, author := range resp.Msg.Authors {
		fmt.Printf("  Author: %s %s\n", author.GivenName, author.FamilyName)
	}

	return nil
}
