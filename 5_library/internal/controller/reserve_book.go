package controller

import (
	"context"

	generated "github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *implementation) ReserveBook(ctx context.Context, req *generated.ReserveBookRequest) (*generated.ReserveBookResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	book, err := i.booksUseCase.ReserveBook(ctx, req.GetId(), req.GetBookedBy())
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &generated.ReserveBookResponse{
		Book: &generated.Book{
			Id:        book.ID,
			Name:      book.Name,
			AuthorId:  book.AuthorIDs,
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
			Booked:    book.Booked,
			BookedBy:  safeString(book.BookedBy),
			ReservationStart: safeTimestamp(book.ReservationStart),
			ReservationEnd:   safeTimestamp(book.ReservationEnd),
		},
	}, nil
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}