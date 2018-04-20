package dynamostore

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

// NotFoundError is a hleper method to determine if an
// encountered error is due to a 404
func NotFoundError(err error) bool {
	return errors.Cause(err).Error() == dynamodb.ErrCodeResourceNotFoundException
}
