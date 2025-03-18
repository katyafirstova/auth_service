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

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req model.UpdateUser
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid  = gofakeit.UUID()
		name  = gofakeit.Word()
		email = gofakeit.Email()
		role  = gofakeit.IntRange(0, 2)

		repoErr = fmt.Errorf("repo error")

		req = &model.UpdateUser{
			Name:  &name,
			Email: &email,
			Role:  model.Role(role),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         error
		err          error
		userRepoMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx, req: *req,
			},
			want: nil,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, uuid, *req).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: *req,
			},
			want: repoErr,
			err:  repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, uuid, *req).Return(repoErr)
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

			err := service.Update(tt.args.ctx, uuid, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
