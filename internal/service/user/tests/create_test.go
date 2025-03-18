package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/katyafirstova/auth_service/internal/model"
	"github.com/katyafirstova/auth_service/internal/repository"
	serviceMocks "github.com/katyafirstova/auth_service/internal/service/mocks"
	"github.com/katyafirstova/auth_service/internal/service/user"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req model.CreateUser
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid     = gofakeit.UUID()
		name     = gofakeit.Word()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, false, false, false, 32)
		role     = gofakeit.IntRange(0, 2)

		repoErr = fmt.Errorf("repo error")

		req = &model.CreateUser{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            model.Role(role),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         string
		err          error
		userRepoMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx, req: *req,
			},
			want: uuid,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, *req).Return(uuid, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: *req,
			},
			want: "",
			err:  repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, *req).Return("", repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.userRepoMock(mc)
			service := user.NewMockService(noteServiceMock)

			newID, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
