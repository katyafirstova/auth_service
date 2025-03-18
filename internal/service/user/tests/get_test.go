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

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid      = gofakeit.UUID()
		name      = gofakeit.Word()
		email     = gofakeit.Email()
		role      = gofakeit.IntRange(0, 2)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		res = &model.User{
			UUID:      uuid,
			Name:      name,
			Email:     email,
			Role:      model.Role(role),
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         model.User
		err          error
		userRepoMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: uuid,
			},
			want: *res,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(*res, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: uuid,
			},
			want: *res,
			err:  repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(*res, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userRepoMock(mc)
			service := user.NewMockService(userServiceMock)

			newID, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
