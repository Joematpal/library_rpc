package flags

import (
	"github.com/urfave/cli/v3"
)

// Flag names as constants
const ()

type Flag = string

const (
	Flag_authors = "authors"
	Flag_book_id = "book-id"
	Flag_isbn    = "isbn"
	Flag_server  = "server"
	Flag_title   = "title"
)

// Global flags
var (
	ServerURL = &cli.StringFlag{
		Name:    "server",
		Aliases: []string{"s"},
		Value:   "http://localhost:8080",
		Usage:   "Library server URL",
	}
)

// Add command flags
var (
	Title = &cli.StringFlag{
		Name:     Flag_title,
		Aliases:  []string{"t"},
		Usage:    "Book title",
		Required: true,
	}

	Authors = &cli.StringSliceFlag{
		Name:     Flag_authors,
		Aliases:  []string{"a", "author"},
		Usage:    "Book authors in format 'given_name family_name' (can specify multiple)",
		Required: true,
	}

	ISBN = &cli.StringFlag{
		Name:     Flag_isbn,
		Usage:    "Book ISBN",
		Required: true,
	}

	BookID = &cli.StringFlag{
		Name:     "book-id",
		Aliases:  []string{"i"},
		Usage:    "Book ID",
		Required: true,
	}
)
