package controller

import (
	"testing"

	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/usecase/library"
	"go.uber.org/zap"
)

func Test_implementation_GetAuthorBooks(t *testing.T) {
	t.Parallel()
	type fields struct {
		logger        *zap.Logger
		booksUseCase  library.BooksUseCase
		authorUseCase library.AuthorUseCase
	}
	type args struct {
		req    *generated.GetAuthorBooksRequest
		stream generated.Library_GetAuthorBooksServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
				req:    &generated.GetAuthorBooksRequest{AuthorId: ""},
				stream: nil,
			},
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
			if err := i.GetAuthorBooks(tt.args.req, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("GetAuthorBooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
