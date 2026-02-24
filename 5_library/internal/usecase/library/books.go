package library

import (
	"context"
    "errors"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
)

func (l *libraryImpl) RegisterBook(ctx context.Context, name string, authorIDs []string) (entity.Book, error) {
	return l.booksRepository.CreateBook(ctx, entity.Book{
		ID:        uuid.New().String(),
		Name:      name,
		AuthorIDs: authorIDs,
	})
}

func (l *libraryImpl) GetBook(ctx context.Context, bookID string) (entity.Book, error) {
	return l.booksRepository.GetBook(ctx, bookID)
}

func (l *libraryImpl) UpdateBook(ctx context.Context, bookID string, name string, authorIDs []string) error {
	return l.booksRepository.UpdateBook(ctx, bookID, name, authorIDs)
}

func (l *libraryImpl) ReserveBook(ctx context.Context, bookID string, bookedBy string) (entity.Book, error) {
	book, err := l.booksRepository.GetBook(ctx, bookID)
	if err != nil {
		return entity.Book{}, err
	}
	if book.Booked {
		return entity.Book{}, errors.New("book already reserved")
	}
	return l.booksRepository.ReserveBook(ctx, bookID, bookedBy)
}

func (l *libraryImpl) ReleaseBook(ctx context.Context, bookID string) (entity.Book, error) {
	book, err := l.booksRepository.GetBook(ctx, bookID)
	if err != nil {
		return entity.Book{}, err
	}
	if !book.Booked {
		return entity.Book{}, errors.New("book is not reserved")
	}
	return l.booksRepository.ReleaseBook(ctx, bookID)
}
