package example

import (
	"fmt"

	"github.com/d-tsuji/example/gen/models"
	"github.com/d-tsuji/example/gen/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p operations.GetUsersParams) middleware.Responder {
	ctx := p.HTTPRequest.Context()
	users, err := scanUsers(ctx)
	if err != nil {
		return operations.NewGetUsersInternalServerError().WithPayload(&models.Error{
			Message: fmt.Sprintf("scan users error: %v", err),
		})
	}
	var resp models.Users
	for _, u := range users {
		u := u
		resp = append(resp, &models.User{
			UserID: &u.UserID,
			Name:   &u.UserName,
		})
	}
	return operations.NewGetUsersOK().WithPayload(resp)
}

func PostUsers(p operations.PostUsersParams) middleware.Responder {
	ctx := p.HTTPRequest.Context()
	u := &User{
		UserID:   *p.PostUsers.UserID,
		UserName: *p.PostUsers.Name,
	}
	if err := createUser(ctx, u); err != nil {
		return operations.NewPostUsersInternalServerError().WithPayload(&models.Error{
			Message: fmt.Sprintf("create user %v, error: %v", *u, err),
		})
	}
	return operations.NewPostUsersOK().WithPayload(p.PostUsers)
}
