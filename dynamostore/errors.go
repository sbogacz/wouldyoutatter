package dynamostore

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// NotFoundError is a hleper method to determine if an
// encountered error is due to a 404
func NotFoundError(err error) bool {
	return err.Error() == dynamodb.ErrCodeResourceNotFoundException
}
