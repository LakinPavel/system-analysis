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

func Test_implementation_RegisterAuthor(t *testing.T) {
	t.Parallel()
	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		ctx context.Context
		req *generated.RegisterAuthorRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generated.RegisterAuthorResponse
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
				req: &generated.RegisterAuthorRequest{},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "CorrectAuthorUseCase",
			fields: fields{
				logger:        zap.NewNop(),
				booksUseCase:  nil,
				authorUseCase: myMocks.CorrectAuthorMock(t),
			},
			args: args{
				ctx: context.Background(),
				req: &generated.RegisterAuthorRequest{Name: "testName"},
			},
			want:    &generated.RegisterAuthorResponse{},
			wantErr: false,
		},

		{
			name: "IncorrectAuthorUseCase",
			fields: fields{
				logger:        zap.NewNop(),
				booksUseCase:  nil,
				authorUseCase: myMocks.IncorrectAuthorMock(t),
			},
			args: args{
				ctx: context.Background(),
				req: &generated.RegisterAuthorRequest{Name: "testName"},
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
			got, err := i.RegisterAuthor(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterAuthor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
