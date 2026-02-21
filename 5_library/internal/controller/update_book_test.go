package controller

import (
	"context"
	"reflect"
	"testing"

	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/usecase/library"
	myMocks "github.com/project/library/mocks"
	"go.uber.org/zap"
)

func Test_implementation_UpdateBook(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		ctx context.Context
		req *generated.UpdateBookRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generated.UpdateBookResponse
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
				req: &generated.UpdateBookRequest{},
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
				req: &generated.UpdateBookRequest{Id: "9e370e77-9e80-4b39-bcb2-8fe8789e25f8", Name: "test_name", AuthorIds: make([]string, 0)},
			},
			want:    &generated.UpdateBookResponse{},
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
				req: &generated.UpdateBookRequest{Name: "test_name"},
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
			got, err := i.UpdateBook(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}
