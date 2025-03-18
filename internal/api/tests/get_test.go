package tests

import (
	"context"
	"fmt"
	_ "fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/katyafirstova/auth_service/internal/api/user"
	"github.com/katyafirstova/auth_service/internal/model"
	"github.com/katyafirstova/auth_service/internal/service"
	serviceMocks "github.com/katyafirstova/auth_service/internal/service/mocks"
	desc "github.com/katyafirstova/auth_service/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid      = gofakeit.UUID()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = gofakeit.IntRange(0, 2)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("Service error")

		req = &desc.GetRequest{
			Uuid: uuid,
		}

		modelReq = model.User{
			UUID:      uuid,
			Name:      name,
			Email:     email,
			Role:      model.Role(role),
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}

		res = &desc.GetResponse{
			Uuid:      uuid,
			Name:      name,
			Email:     email,
			Role:      desc.Role(role),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx, req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(modelReq, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx, req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(modelReq, serviceErr)
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

			newID, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
