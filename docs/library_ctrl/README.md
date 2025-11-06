# Library Control CLI Documentation

The `library_ctrl` CLI provides a command-line interface for managing books in the library RPC system.

## Prerequisites

1. Start the library server:
   ```bash
   go run ./cmd/library
   ```
   The server will run on `http://localhost:8080` by default.

2. Build or run the CLI client:
   ```bash
   go run ./cmd/library_ctrl [command]
   ```

## Available Commands

### Global Options
- `--server, -s`: Library server URL (default: "http://localhost:8080")

### Add Book
Add a new book to the library with ISBN and multiple authors.

```bash
go run ./cmd/library_ctrl add book \
  --title "Book Title" \
  --isbn "978-XXXXXXXXX" \
  --authors "Given Family" \
  --authors "Another Author"
```

**Options:**
- `--title, -t`: Book title (required)
- `--authors, -a, --author`: Book authors in format 'given_name family_name' (required, can specify multiple)
- `--isbn`: Book ISBN (required)

**Example:**
```bash
go run ./cmd/library_ctrl add book \
  --title "The Go Programming Language" \
  --isbn "978-0134190440" \
  --authors "Alan Donovan" \
  --authors "Brian Kernighan"
```

### Get Book
Retrieve a book by its ID.

```bash
go run ./cmd/library_ctrl get book [book_id]
# OR
go run ./cmd/library_ctrl get book --book-id [book_id]
```

**Options:**
- `--book-id, -i`: Book ID (required if not provided as argument)

**Example:**
```bash
go run ./cmd/library_ctrl get book cm2r8abc0000001234567890
```

### List Books
List all books in the library.

```bash
go run ./cmd/library_ctrl list books
```

**Example Output:**
```
Found 2 books:

1. Book ID: cm2r8abc0000001234567890
   ISBN: 978-0134190440
   Title: The Go Programming Language
   Author 1: Alan Donovan
   Author 2: Brian Kernighan

2. Book ID: cm2r8def0000001234567891
   ISBN: 978-0132350884
   Title: Clean Code
   Author 1: Robert Martin
```

### Update Book
Update an existing book's information.

```bash
go run ./cmd/library_ctrl update book \
  --book-id [book_id] \
  --title "Updated Title" \
  --isbn "978-XXXXXXXXX" \
  --authors "Updated Author"
```

**Options:**
- `--book-id`: Book ID to update (required)
- `--title, -t`: Updated book title (required)
- `--authors, -a, --author`: Updated book authors in format 'given_name family_name' (required, can specify multiple)
- `--isbn`: Updated book ISBN (required)

**Example:**
```bash
go run ./cmd/library_ctrl update book \
  --book-id cm2r8def0000001234567891 \
  --title "Clean Code: Updated Edition" \
  --isbn "978-0132350884" \
  --authors "Robert C Martin"
```

### Delete Book
Delete a book by its ID.

```bash
go run ./cmd/library_ctrl delete book [book_id]
# OR
go run ./cmd/library_ctrl delete book --book-id [book_id]
```

**Options:**
- `--book-id, -i`: Book ID (required if not provided as argument)

**Example:**
```bash
go run ./cmd/library_ctrl delete book cm2r8def0000001234567891
```

## Complete Testing Workflow

Here's a complete workflow to test all functionality:

### 1. Start the Server
```bash
go run ./cmd/library
```

### 2. Add Some Books
```bash
# Add first book with multiple authors
go run ./cmd/library_ctrl add book \
  --title "The Go Programming Language" \
  --isbn "978-0134190440" \
  --authors "Alan Donovan" \
  --authors "Brian Kernighan"

# Add second book with single author
go run ./cmd/library_ctrl add book \
  --title "Clean Code" \
  --isbn "978-0132350884" \
  --authors "Robert Martin"

# Add third book
go run ./cmd/library_ctrl add book \
  --title "Design Patterns" \
  --isbn "978-0201633612" \
  --authors "Erich Gamma" \
  --authors "Richard Helm" \
  --authors "Ralph Johnson" \
  --authors "John Vlissides"
```

### 3. List All Books
```bash
go run ./cmd/library_ctrl list books
```

### 4. Get Specific Book
```bash
# Use a book_id from the add responses
go run ./cmd/library_ctrl get book <book_id_from_add_response>
```

### 5. Update a Book
```bash
# Update the Clean Code book
go run ./cmd/library_ctrl update book \
  --book-id <book_id_from_clean_code> \
  --title "Clean Code: A Handbook of Agile Software Craftsmanship" \
  --isbn "978-0132350884" \
  --authors "Robert C Martin"
```

### 6. Delete a Book
```bash
# Delete one of the books
go run ./cmd/library_ctrl delete book <book_id_to_delete>
```

### 7. Verify Changes
```bash
# List books again to see the changes
go run ./cmd/library_ctrl list books
```

## Author Format

Authors must be specified in the format `"given_name family_name"`. For authors with multiple family names, everything after the first space is treated as the family name.

**Examples:**
- `"John Doe"` → given_name: "John", family_name: "Doe"
- `"Mary Jane Smith"` → given_name: "Mary", family_name: "Jane Smith"
- `"Jean-Claude Van Damme"` → given_name: "Jean-Claude", family_name: "Van Damme"

## Error Handling

The CLI provides clear error messages for common issues:
- Missing required fields
- Invalid author format
- Book not found
- Server connection issues

## Book ID Generation

Books are automatically assigned unique internal IDs using the xid library. These IDs are separate from the ISBN and are used for internal tracking and operations.

## Data Validation

The system validates:
- Title: Required, 1-200 characters
- ISBN: Required, 10-17 characters with specific format
- Authors: At least one author required, each with given and family names
- Book ID: Required for update/delete operations