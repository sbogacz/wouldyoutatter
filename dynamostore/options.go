package dynamostore

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type updateTTLReq struct {
	input *dynamodb.UpdateTimeToLiveInput
}

// NewTTLOption takes an UpdateTimeToLiveInput and creates a TableOption
// that can be used in the lazy creation of tables
func NewTTLOption(input *dynamodb.UpdateTimeToLiveInput) TableOption {
	return &updateTTLReq{input: input}
}

// Name describes the TTL table option
func (u *updateTTLReq) Name() string {
	return "TTL"
}

// Send allows the TTL option to be applied against dynamo
func (u *updateTTLReq) Send(db *dynamodb.DynamoDB) error {
	req := db.UpdateTimeToLiveRequest(u.input)
	_, err := req.Send()
	return err
}

/*type createGSIReq struct {
	input *dynamodb.CreateGlobalSecondaryIndexInput
}

// NewGSIOption takes a CreateGlobalSecondaryIndexInput and creates a TableOption
// that can be used in the lazy creation of tables
func NewGSIOption(input *dynamodb.CreateGlobalSecondaryIndexInput) TableOption {
	return &createGSIReq{input: input}
}

// Name describes the TTL table option
func (c *createGSIReq) Name() string {
	return "GSI"
}

// Send allows the GSI option to be applied against dynamo
func (c *createGSIReq) Send(db *dynamodb.DynamoDB) error {
	req := db.CreateGlobalSecondaryIndexRequest(u.input)
	_, err := req.Send()
	return err
}*/
