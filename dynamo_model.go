package example

type User struct {
	UserID   string `dynamo:"user_id,hash"`
	UserName string `dynamo:"user_name"`
}
