package v1

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	v1 "rdbms-go.davideimola.dev/gen/proto/davideimola/users/v1"
	"rdbms-go.davideimola.dev/gen/proto/davideimola/users/v1/usersv1connect"
	"rdbms-go.davideimola.dev/internal/queries"
)

type srv struct {
	q *queries.Queries
}

func NewUsersService(q *queries.Queries) usersv1connect.UsersServiceHandler {
	return &srv{
		q: q,
	}
}

func newUser(user queries.User) *v1.User {
	return &v1.User{
		Id:   user.ID,
		Name: user.Name,
	}
}

func (s srv) CreateUser(ctx context.Context, req *connect_go.Request[v1.CreateUserRequest]) (*connect_go.Response[v1.CreateUserResponse], error) {
	user, err := s.q.CreateUser(ctx, req.Msg.Name)
	if err != nil {
		return nil, err
	}

	return connect_go.NewResponse(&v1.CreateUserResponse{
		User: newUser(user),
	}), nil
}

func (s srv) ListUsers(ctx context.Context, req *connect_go.Request[v1.ListUsersRequest]) (*connect_go.Response[v1.ListUsersResponse], error) {
	users, err := s.q.ListUsers(ctx, queries.ListUsersParams{
		Offset: req.Msg.Offset,
		Limit:  req.Msg.Limit,
	})
	if err != nil {
		return nil, err
	}

	tot, err := s.q.CountUsers(ctx)

	res := make([]*v1.User, len(users))
	for i, user := range users {
		res[i] = newUser(user)
	}

	return connect_go.NewResponse(&v1.ListUsersResponse{
		Users: res,
		Totat: int32(tot),
	}), nil
}

func (s srv) DeleteUser(ctx context.Context, req *connect_go.Request[v1.DeleteUserRequest]) (*connect_go.Response[v1.DeleteUserResponse], error) {
	user, err := s.q.DeleteUser(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect_go.NewResponse(&v1.DeleteUserResponse{
		User: newUser(user),
	}), nil
}

func (s srv) GetUser(ctx context.Context, c *connect_go.Request[v1.GetUserRequest]) (*connect_go.Response[v1.GetUserResponse], error) {
	user, err := s.q.GetUser(ctx, c.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect_go.NewResponse(&v1.GetUserResponse{
		User: newUser(user),
	}), nil
}
