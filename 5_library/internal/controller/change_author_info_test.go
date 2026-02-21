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

func Test_implementation_ChangeAuthorInfo(t *testing.T) {
	t.Parallel()
	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		ctx context.Context
		req *generated.ChangeAuthorInfoRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *generated.ChangeAuthorInfoResponse
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
				req: &generated.ChangeAuthorInfoRequest{},
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
				req: &generated.ChangeAuthorInfoRequest{Id: "ed327ca8-6796-4e2a-911b-42b1cb4441da", Name: "testName"},
			},
			want:    &generated.ChangeAuthorInfoResponse{},
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
				req: &generated.ChangeAuthorInfoRequest{Id: "ed327ca8-6796-4e2a-911b-42b1cb4441da", Name: "testName"},
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
			got, err := i.ChangeAuthorInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeAuthorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeAuthorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
