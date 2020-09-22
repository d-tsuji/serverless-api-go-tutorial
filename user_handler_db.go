package example

import (
	"context"
	"errors"

	"github.com/guregu/dynamo"
)

func scanUsers(ctx context.Context) ([]User, error) {
	var resp []User
	table := gdb.Table(usersTable)
	if err := table.Scan().AllWithContext(ctx, &resp); err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return nil, nil
		}
		return resp, err
	}
	return resp, nil
}

func createUser(ctx context.Context, u *User) error {
	if err := gdb.Table(usersTable).Put(u).RunWithContext(ctx); err != nil {
		return err
	}
	return nil
}
