package tests

import (
	"context"
	"fmt"
	_ "fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/katyafirstova/auth_service/internal/api/user"
	"github.com/katyafirstova/auth_service/internal/model"
	"github.com/katyafirstova/auth_service/internal/service"
	serviceMocks "github.com/katyafirstova/auth_service/internal/service/mocks"
	desc "github.com/katyafirstova/auth_service/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid     = gofakeit.UUID()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, false, false, false, 32)
		role     = gofakeit.IntRange(0, 2)

		serviceErr = fmt.Errorf("Service error")
		//passwordsErr = fmt.Errorf("Passwords are not equal")

		modelReq = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            desc.Role(role),
		}

		req = &model.CreateUser{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            model.Role(role),
		}

		res = &desc.CreateResponse{
			Uuid: uuid,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx, req: modelReq,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, *req).Return(uuid, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: modelReq,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, *req).Return("", serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(noteServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
