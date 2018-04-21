package dynamostore

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

// NotFoundError is a hleper method to determine if an
// encountered error is due to a 404
func NotFoundError(err error) bool {
	return errors.Cause(err).Error() == dynamodb.ErrCodeResourceNotFoundException
}

// TableNotFoundError is a hleper method to determine if an
// encountered error is due to a 404
func TableNotFoundError(err error) bool {
	errRoot := strings.Split(errors.Cause(err).Error(), ":")[0]
	fmt.Printf("\n\n\nerrRoot: %s\n\n\n", errRoot)
	return errors.Cause(err).Error() == dynamodb.ErrCodeTableNotFoundException ||
		errRoot == dynamodb.ErrCodeResourceNotFoundException
}
