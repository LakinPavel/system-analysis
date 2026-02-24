//nolint:errcheck //why
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/internal/entity"
)

var _ BooksRepository = (*postgresRepository)(nil)
var _ AuthorRepository = (*postgresRepository)(nil)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

func (p *postgresRepository) CreateAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	tx, err := p.db.Begin(ctx)

	if err != nil {
		return entity.Author{}, err
	}

	defer tx.Rollback(ctx)

	const queryAuthor = `
INSERT INTO author (id, name)
VALUES ($1,$2)
`

	result := entity.Author{
		ID:   author.ID,
		Name: author.Name,
	}

	_, err = tx.Exec(ctx, queryAuthor, author.ID, author.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.Author{}, entity.ErrAuthorAlreadyExists
		}
		return entity.Author{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return entity.Author{}, err
	}

	return result, nil
}

func (p *postgresRepository) GetAuthor(ctx context.Context, authorID string) (entity.Author, error) {
	const query = `
    SELECT id, name
    FROM author
    WHERE id = $1
    `

	var author entity.Author
	err := p.db.QueryRow(ctx, query, authorID).
		Scan(&author.ID, &author.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Author{}, entity.ErrAuthorNotFound
	}

	if err != nil {
		return entity.Author{}, err
	}

	return author, nil
}

func (p *postgresRepository) CreateBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	tx, err := p.db.Begin(ctx)

	if err != nil {
		return entity.Book{}, err
	}

	defer tx.Rollback(ctx)

	const queryBook = `
INSERT INTO book (id, name, booked, booked_by, reservation_start, reservation_end)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING created_at, updated_at
`

	result := entity.Book{
		ID:        book.ID,
		Name:      book.Name,
		AuthorIDs: book.AuthorIDs,
		Booked:    false,
		BookedBy:  nil,
        ReservationStart: nil,
        ReservationEnd:   nil,
	}

	err = tx.QueryRow(ctx, queryBook, book.ID, book.Name, false, nil, nil, nil).Scan(&result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.Book{}, entity.ErrBookAlreadyExists
		}
		return entity.Book{}, err
	}
	const queryAuthorBooks = `
INSERT INTO author_book
(author_id, book_id)
VALUES ($1, $2)
`

	for _, authorID := range book.AuthorIDs {
		const checkAuthorQuery = `
SELECT id 
FROM author 
WHERE id = $1
`
		var temp string
		err = tx.QueryRow(ctx, checkAuthorQuery, authorID).Scan(&temp)
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Book{}, entity.ErrAuthorNotFound
		}
		if err != nil {
			return entity.Book{}, err
		}

		_, err = tx.Exec(ctx, queryAuthorBooks, authorID, book.ID)

		if err != nil {
			return entity.Book{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return entity.Book{}, err
	}

	return result, nil
}

func (p *postgresRepository) GetBook(ctx context.Context, bookID string) (entity.Book, error) {
	const query = `
SELECT id, name, created_at, updated_at, booked, booked_by, reservation_start, reservation_end
FROM book
WHERE id = $1
`

	var book entity.Book
	var bookedBy sql.NullString
    var reservationStart, reservationEnd sql.NullTime
	err := p.db.QueryRow(ctx, query, bookID).
		Scan(&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.Booked&bookedBy, &reservationStart, &reservationEnd,)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Book{}, entity.ErrBookNotFound
	}

	if err != nil {
		return entity.Book{}, err
	}

	if bookedBy.Valid {
        book.BookedBy = &bookedBy.String
    }
    if reservationStart.Valid {
        book.ReservationStart = &reservationStart.Time
    }
    if reservationEnd.Valid {
        book.ReservationEnd = &reservationEnd.Time
    }

	const queryAuthors = `
SELECT author_id
FROM author_book
WHERE book_id = $1
`

	rows, err := p.db.Query(ctx, queryAuthors, bookID)

	if err != nil {
		return entity.Book{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var authorID string

		if err := rows.Scan(&authorID); err != nil {
			return entity.Book{}, err
		}

		book.AuthorIDs = append(book.AuthorIDs, authorID)
	}

	return book, nil
}

func (p *postgresRepository) UpdateBook(ctx context.Context, bookID string, newName string, newAuthorIDs []string) error {
	const checkBookQuery = `
SELECT id
FROM book
WHERE id = $1
`
	var temp string
	err := p.db.QueryRow(ctx, checkBookQuery, bookID).Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.ErrBookNotFound
	}
	if err != nil {
		return err
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const queryUpdateBook = `
    UPDATE book 
    SET name = $1, updated_at = now()
    WHERE id = $2
    `
	_, err = tx.Exec(ctx, queryUpdateBook, newName, bookID)
	if err != nil {
		return err
	}

	const queryDeleteAuthors = `
    DELETE FROM author_book 
    WHERE book_id = $1
    `
	_, err = tx.Exec(ctx, queryDeleteAuthors, bookID)
	if err != nil {
		return err
	}

	const queryInsertAuthors = `
    INSERT INTO author_book (author_id, book_id)
    VALUES ($1, $2)
    `
	for _, authorID := range newAuthorIDs {
		const checkAuthorQuery = `SELECT id
FROM author
WHERE id = $1
`
		var authorTemp string
		err = tx.QueryRow(ctx, checkAuthorQuery, authorID).Scan(&authorTemp)
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrAuthorNotFound
		}
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, queryInsertAuthors, authorID, bookID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) ChangeAuthorInfo(ctx context.Context, authorID string, authorName string) error {
	const query = `
    UPDATE author 
    SET name = $1
    WHERE id = $2
    `

	result, err := p.db.Exec(ctx, query, authorName, authorID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrAuthorNotFound
	}

	return nil
}

func (p *postgresRepository) GetAuthorBooks(ctx context.Context, authorID string) ([]entity.Book, error) {
	const checkAuthorQuery = `
    SELECT id
    FROM author
    WHERE id = $1
    `
	var temp string
	err := p.db.QueryRow(ctx, checkAuthorQuery, authorID).Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrAuthorNotFound
	}
	if err != nil {
		return nil, err
	}

	const query = `
    SELECT b.id, b.name, b.created_at, b.updated_at, b.booked, b.booked_by, b.reservation_start, b.reservation_end
    FROM book b
    INNER JOIN author_book ab ON b.id = ab.book_id
    WHERE ab.author_id = $1
    `

	rows, err := p.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		var bookedBy sql.NullString
        var reservationStart, reservationEnd sql.NullTime
		if err := rows.Scan(&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.Booked, &bookedBy, &reservationStart, &reservationEnd,); err != nil {
			return nil, err
		}

		const authorsQuery = `
        SELECT author_id 
        FROM author_book 
        WHERE book_id = $1
        `
		authorRows, err := p.db.Query(ctx, authorsQuery, book.ID)
		if err != nil {
			return nil, err
		}

		for authorRows.Next() {
			var authorID string
			if err := authorRows.Scan(&authorID); err != nil {
				authorRows.Close()
				return nil, err
			}
			book.AuthorIDs = append(book.AuthorIDs, authorID)
		}
		authorRows.Close()

		books = append(books, book)
	}

	return books, nil
}

func (p *postgresRepository) ReserveBook(ctx context.Context, bookID string, bookedBy string) (entity.Book, error) {
    tx, err := p.db.Begin(ctx)
    if err != nil {
        return entity.Book{}, err
    }
    defer tx.Rollback(ctx)

    const updateQuery = `
        UPDATE book
        SET booked = true,
            booked_by = $1,
            reservation_start = now(),
            reservation_end = now() + interval '14 days',
            updated_at = now()
        WHERE id = $2
        RETURNING id, name, created_at, updated_at, booked,
                  booked_by, reservation_start, reservation_end
    `

    var book entity.Book
	var bookedByDB sql.NullString
    var reservationStart, reservationEnd sql.NullTime

    err = tx.QueryRow(ctx, updateQuery, bookedBy, bookID).Scan(
        &book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.Booked, &bookedByDB, &reservationStart, &reservationEnd,
    )
    if err != nil {
        return entity.Book{}, err
    }

	if bookedByDB.Valid {
        book.BookedBy = &bookedByDB.String
    }
    if reservationStart.Valid {
        book.ReservationStart = &reservationStart.Time
    }
    if reservationEnd.Valid {
        book.ReservationEnd = &reservationEnd.Time
    }

    const authorsQuery = `
        SELECT author_id
        FROM author_book
        WHERE book_id = $1
    `
    rows, err := tx.Query(ctx, authorsQuery, bookID)
    if err != nil {
        return entity.Book{}, err
    }
    defer rows.Close()

    for rows.Next() {
        var authorID string
        if err := rows.Scan(&authorID); err != nil {
            return entity.Book{}, err
        }
        book.AuthorIDs = append(book.AuthorIDs, authorID)
    }

    if err = tx.Commit(ctx); err != nil {
        return entity.Book{}, err
    }

    return book, nil
}

func (p *postgresRepository) ReleaseBook(ctx context.Context, bookID string) (entity.Book, error) {
    tx, err := p.db.Begin(ctx)
    if err != nil {
        return entity.Book{}, err
    }
    defer tx.Rollback(ctx)

    const updateQuery = `
        UPDATE book
        SET booked = false,
            booked_by = NULL,
            reservation_start = NULL,
            reservation_end = NULL,
            updated_at = now()
        WHERE id = $1
        RETURNING id, name, created_at, updated_at, booked,
                  booked_by, reservation_start, reservation_end
    `

    var book entity.Book
	var bookedByDB sql.NullString
    var reservationStart, reservationEnd sql.NullTime

    err = tx.QueryRow(ctx, updateQuery, bookID).Scan(
        &book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.Booked, &bookedByDB, &reservationStart, &reservationEnd,
    )
    if err != nil {
        return entity.Book{}, err
    }

	if bookedByDB.Valid {
        book.BookedBy = &bookedByDB.String
    }
    if reservationStart.Valid {
        book.ReservationStart = &reservationStart.Time
    }
    if reservationEnd.Valid {
        book.ReservationEnd = &reservationEnd.Time
    }

	const authorsQuery = `
        SELECT author_id
        FROM author_book
        WHERE book_id = $1
    `
    rows, err := tx.Query(ctx, authorsQuery, bookID)
    if err != nil {
        return entity.Book{}, err
    }
    defer rows.Close()

    for rows.Next() {
        var authorID string
        if err := rows.Scan(&authorID); err != nil {
            return entity.Book{}, err
        }
        book.AuthorIDs = append(book.AuthorIDs, authorID)
    }


    if err = tx.Commit(ctx); err != nil {
        return entity.Book{}, err
    }

    return book, nil
}