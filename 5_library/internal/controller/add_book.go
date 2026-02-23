package controller

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) AddBook(ctx context.Context, req *library.AddBookRequest) (*library.AddBookResponse, error) {
	if err := req.ValidateAll(); err != nil {
		// i.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	book, err := i.booksUseCase.RegisterBook(ctx, req.GetName(), req.GetAuthorIds())

	if err != nil {
		return nil, i.convertErr(err)
	}
	// i.logger.Info("done")
	return &library.AddBookResponse{
		Book: &library.Book{
			Id:        book.ID,
			Name:      book.Name,
			AuthorId:  book.AuthorIDs,
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.CreatedAt),
			Booked:    book.Booked,
		},
	}, nil
}
