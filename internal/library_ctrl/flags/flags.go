package flags

import (
	"github.com/urfave/cli/v3"
)

// Flag names as constants
const ()

type Flag = string

const (
	Flag_author  = "author"
	Flag_book_id = "book-id"
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

	Author = &cli.StringFlag{
		Name:     Flag_author,
		Aliases:  []string{"a"},
		Usage:    "Book author",
		Required: true,
	}

	BookID = &cli.StringFlag{
		Name:     "book-id",
		Aliases:  []string{"i"},
		Usage:    "Book ID",
		Required: true,
	}
)
