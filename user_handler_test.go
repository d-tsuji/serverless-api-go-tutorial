package example

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/d-tsuji/example/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/d-tsuji/example/gen/restapi/operations"
	"github.com/go-openapi/runtime"
	"github.com/guregu/dynamo"
	"github.com/nsf/jsondiff"
)

func init() {
	dbEndpoint := "http://localhost:4566"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "local",
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Endpoint:   aws.String(dbEndpoint),
			DisableSSL: aws.Bool(true),
		},
	}))
	gdb = dynamo.New(sess)
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		input          func(t *testing.T)
		wantStatusCode int
		want           string
	}{
		{
			name: "複数件のユーザの取得",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).RunWithContext(context.TODO())
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
				inputUsers := []User{{UserID: "001", UserName: "gopher"}, {UserID: "002", UserName: "rubyist"}}
				for _, u := range inputUsers {
					if err := gdb.Table(usersTable).Put(u).RunWithContext(context.TODO()); err != nil {
						t.Errorf("dynamo input user %v: %v", u, err)
					}
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_get_users_1.json",
		},
		{
			name: "ユーザ0件",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).RunWithContext(context.TODO())
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_get_users_2.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input(t)
			t.Cleanup(func() {
				if err := gdb.Table(usersTable).DeleteTable().RunWithContext(context.TODO()); err != nil {
					t.Fatalf("dynamo delete table %s: %v", usersTable, err)
				}
			})

			p := operations.NewGetUsersParams()
			p.HTTPRequest = httptest.NewRequest("GET", "/v1/users", nil)

			resp := GetUsers(p)

			w := httptest.NewRecorder()
			resp.WriteResponse(w, runtime.JSONProducer())

			want, err := ioutil.ReadFile(tt.want)
			if err != nil {
				t.Fatalf("want file read: %v", err)
			}

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("status got %v, but want %v", w.Result().StatusCode, tt.wantStatusCode)
			}

			opt := jsondiff.DefaultConsoleOptions()
			if d, s := jsondiff.Compare(w.Body.Bytes(), want, &opt); d != jsondiff.FullMatch {
				t.Errorf("unmatch, got=%s, want=%s, diff=%s", string(w.Body.Bytes()), string(want), s)
			}
		})
	}
}

func TestPostUsers(t *testing.T) {
	tests := []struct {
		name           string
		input          func(t *testing.T)
		wantStatusCode int
		want           string
	}{
		{
			name: "ユーザの取得成功",
			input: func(t *testing.T) {
				err := gdb.CreateTable(usersTable, User{}).Provision(1, 1).RunWithContext(context.TODO())
				if err != nil {
					t.Errorf("dynamo create table %s: %v", usersTable, err)
				}
			},
			wantStatusCode: 200,
			want:           "./testdata/want_post_users_1.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input(t)
			t.Cleanup(func() {
				if err := gdb.Table(usersTable).DeleteTable().RunWithContext(context.TODO()); err != nil {
					t.Fatalf("dynamo delete table %s: %v", usersTable, err)
				}
			})

			p := operations.NewPostUsersParams()
			userID := "003"
			userName := "ferris"
			p.PostUsers = &models.User{UserID: &userID, Name: &userName}
			p.HTTPRequest = httptest.NewRequest("POST", "/v1/users", nil)

			resp := PostUsers(p)

			w := httptest.NewRecorder()
			resp.WriteResponse(w, runtime.JSONProducer())

			want, err := ioutil.ReadFile(tt.want)
			if err != nil {
				t.Fatalf("want file read: %v", err)
			}

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("status got %v, but want %v", w.Result().StatusCode, tt.wantStatusCode)
			}

			opt := jsondiff.DefaultConsoleOptions()
			if d, s := jsondiff.Compare(w.Body.Bytes(), want, &opt); d != jsondiff.FullMatch {
				t.Errorf("unmatch, got=%s, want=%s, diff=%s", string(w.Body.Bytes()), string(want), s)
			}
		})
	}
}
