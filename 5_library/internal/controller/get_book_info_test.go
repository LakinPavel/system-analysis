package controller

import (
	"context"
	"reflect"
	"testing"

	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/entity"
	"github.com/project/library/internal/usecase/library"
	myMocks "github.com/project/library/mocks"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_implementation_GetBookInfo(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		ctx context.Context
		req *generated.GetBookInfoRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generated.GetBookInfoResponse
		wantErr bool
	}{
		{
			name: "IncorrectValidate",
			fields: fields{
				logger:        zap.NewNop(),
				booksUseCase:  nil,
				authorUseCase: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &generated.GetBookInfoRequest{},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "CorrectBooksUseCase",
			fields: fields{
				logger:        zap.NewNop(),
				booksUseCase:  myMocks.CorrectBookMock(t),
				authorUseCase: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &generated.GetBookInfoRequest{Id: "9e370e77-9e80-4b39-bcb2-8fe8789e25f8"},
			},
			want:    &generated.GetBookInfoResponse{Book: &generated.Book{Id: "", Name: "", AuthorId: nil, CreatedAt: timestamppb.New(entity.Book{}.CreatedAt), UpdatedAt: timestamppb.New(entity.Book{}.CreatedAt)}},
			wantErr: false,
		},

		{
			name: "IncorrectBooksUseCase",
			fields: fields{
				logger:        zap.NewNop(),
				booksUseCase:  myMocks.IncorrectBookMock(t),
				authorUseCase: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &generated.GetBookInfoRequest{Id: "9e370e77-9e80-4b39-bcb2-8fe8789e25f8"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := &implementation{
				logger:        tt.fields.logger,
				booksUseCase:  tt.fields.booksUseCase,
				authorUseCase: tt.fields.authorUseCase,
			}
			got, err := i.GetBookInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBookInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
