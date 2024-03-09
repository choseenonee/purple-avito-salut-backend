package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"template/internal/model/entities"
	"testing"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserRepository_Get(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%routers' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := InitUserRepo(db)

	tests := []struct {
		name    string
		mock    func()
		want    entities.User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := mock.NewRows([]string{"users.id", "users.name"})
				rows = rows.AddRow(1, "vadim")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT users.id, users.name FROM users WHERE users.id = $1;`)).WithArgs(1).WillReturnRows(rows)
			},
			want: entities.User{
				UserBase: entities.UserBase{
					Name: "vadim",
				},
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "Error row",
			mock: func() {
				rows := sqlxmock.NewRows([]string{"users.id"}).AddRow(1).RowError(0, errors.New("failed to get"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT users.id, users.name FROM users WHERE users.id = $1;`)).WithArgs(1).WillReturnRows(rows)
			},
			want:    entities.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.Get(context.Background(), 1)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
