package controller

import (
	"context"

	generated "github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *implementation) ReleaseBook(ctx context.Context, req *generated.ReleaseBookRequest) (*generated.ReleaseBookResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	book, err := i.booksUseCase.ReleaseBook(ctx, req.GetId())
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &generated.ReleaseBookResponse{
		Book: &generated.Book{
			Id:        book.ID,
			Name:      book.Name,
			AuthorId:  book.AuthorIDs,
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
			Booked:    book.Booked,
		},
	}, nil
}