package example

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	gdb    *dynamo.DB
	region string

	usersTable string
)

func init() {
	region = os.Getenv("AWS_REGION")

	usersTable = os.Getenv("DYNAMO_TABLE_USERS")
	if usersTable == "" {
		log.Fatal("missing env variable: DYNAMO_TABLE_USERS")
	}

	gdb = dynamo.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	})))
}
