package tests

import (
	"context"
	"fmt"
	_ "fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	_ "google.golang.org/grpc/credentials/google"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/katyafirstova/auth_service/internal/api/user"
	"github.com/katyafirstova/auth_service/internal/converter"
	"github.com/katyafirstova/auth_service/internal/service"
	serviceMocks "github.com/katyafirstova/auth_service/internal/service/mocks"
	desc "github.com/katyafirstova/auth_service/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		uuid  = gofakeit.UUID()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = gofakeit.IntRange(0, 2)

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Uuid:  uuid,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
			Role:  desc.Role(role),
		}

		modelReq = converter.UpdateUserToServiceFromAPI(req)
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx, req: req,
			},
			want: new(emptypb.Empty),
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, uuid, modelReq).Return(nil)
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
				mock.UpdateMock.Expect(ctx, uuid, modelReq).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			newID, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
