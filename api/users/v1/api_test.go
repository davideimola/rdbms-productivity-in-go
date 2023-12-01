package v1_test

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	users "rdbms-go.davideimola.dev/api/users/v1"
	dbchema "rdbms-go.davideimola.dev/db"
	v1 "rdbms-go.davideimola.dev/gen/proto/davideimola/users/v1"
	c "rdbms-go.davideimola.dev/gen/proto/davideimola/users/v1/usersv1connect"
	"rdbms-go.davideimola.dev/internal/queries"
	"rdbms-go.davideimola.dev/internal/testutils"
	"testing"
)

var (
	usersCli c.UsersServiceClient
)

func TestMain(m *testing.M) {
	connString, purge, err := testutils.InitPostgres()
	if err != nil {
		logrus.Fatalf("Could not init database: %s", err)
	}

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	if err := dbchema.Migrate(fmt.Sprintf("%s?sslmode=disable", connString)); err != nil {
		logrus.Fatalf("Could not migrate database: %s", err)
	}

	q := queries.New(db)

	mux := http.NewServeMux()
	mux.Handle(c.NewUsersServiceHandler(
		users.NewUsersService(q),
	))
	server := httptest.NewServer(mux)
	server.EnableHTTP2 = true
	defer server.Close()

	usersCli = c.NewUsersServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := purge(); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestListEmptyUsers(t *testing.T) {
	ctx := context.Background()
	req := connect.NewRequest(&v1.ListUsersRequest{
		Offset: 0,
		Limit:  10,
	})

	resp, err := usersCli.ListUsers(ctx, req)
	if err != nil {
		t.Fatalf("Could not list users: %s", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, int32(0), resp.Msg.Totat)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 20; i++ {
		createTestUser(t, fmt.Sprintf("test-%d", i))
	}

	ctx := context.Background()
	req := connect.NewRequest(&v1.ListUsersRequest{
		Offset: 0,
		Limit:  10,
	})

	resp, err := usersCli.ListUsers(ctx, req)
	if err != nil {
		t.Fatalf("Could not list users: %s", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, int32(20), resp.Msg.Totat)
	assert.Equal(t, 10, len(resp.Msg.Users))
}

func TestCreateUser(t *testing.T) {
	user := createTestUser(t, "test")
	assert.Equal(t, "test", user.Name)
}

func createTestUser(t *testing.T, name string) *v1.User {
	t.Helper()

	ctx := context.Background()
	req := connect.NewRequest(&v1.CreateUserRequest{
		Name: name,
	})

	resp, err := usersCli.CreateUser(ctx, req)
	if err != nil {
		logrus.Fatalf("Could not create user: %s", err)
	}

	t.Cleanup(func() {
		ctx := context.Background()
		req := connect.NewRequest(&v1.DeleteUserRequest{
			Id: resp.Msg.User.Id,
		})

		_, err := usersCli.DeleteUser(ctx, req)
		if err != nil {
			logrus.Fatalf("Could not delete user: %s", err)
		}
	})

	return resp.Msg.User
}
