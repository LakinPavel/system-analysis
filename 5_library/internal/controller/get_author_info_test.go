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

func Test_implementation_GetAuthorInfo(t *testing.T) {
	t.Parallel()
	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		ctx context.Context
		req *generated.GetAuthorInfoRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generated.GetAuthorInfoResponse
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
				req: &generated.GetAuthorInfoRequest{},
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
				req: &generated.GetAuthorInfoRequest{Id: "9c6a3e4a-137f-449d-8666-4bce1fb2d117"},
			},
			want:    &generated.GetAuthorInfoResponse{Id: "9c6a3e4a-137f-449d-8666-4bce1fb2d117", Name: "testName"},
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
				req: &generated.GetAuthorInfoRequest{},
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
			got, err := i.GetAuthorInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
