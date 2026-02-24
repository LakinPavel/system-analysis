package controller

import (
	generated "github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) GetAuthorBooks(req *generated.GetAuthorBooksRequest, stream generated.Library_GetAuthorBooksServer) error {
	if err := req.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	books, err := i.authorUseCase.GetAuthorBooks(stream.Context(), req.GetAuthorId())
	if err != nil {
		return i.convertErr(err)
	}

	for _, book := range books {
		protoBook := &generated.Book{
			Id:       book.ID,
			Name:     book.Name,
			AuthorId: book.AuthorIDs,
			Booked:    book.Booked,
			BookedBy:  safeString(book.BookedBy),
			ReservationStart: safeTimestamp(book.ReservationStart),
			ReservationEnd:   safeTimestamp(book.ReservationEnd), 
		}

		if err := stream.Send(protoBook); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
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
