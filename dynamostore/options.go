package dynamostore

import (
	"github.com/aws/aws-sdk-go-v2/aws"
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

type createGSIReq struct {
	gsiUpdates []dynamodb.GlobalSecondaryIndexUpdate
	tableName  string
}

// NewGSIOption takes a set of GSI updates and creates a TableOption
// that can be used in the lazy creation of tables
func NewGSIOption(gsiUpdates []dynamodb.GlobalSecondaryIndexUpdate, tableName string) TableOption {
	return &createGSIReq{
		gsiUpdates: gsiUpdates,
		tableName:  tableName,
	}
}

// Name describes the TTL table option
func (c *createGSIReq) Name() string {
	return "GSI"
}

// Send allows the GSI option to be applied against dynamo
func (c *createGSIReq) Send(db *dynamodb.DynamoDB) error {
	input := &dynamodb.UpdateTableInput{
		TableName:                   aws.String(c.tableName),
		GlobalSecondaryIndexUpdates: c.gsiUpdates,
	}
	req := db.UpdateTableRequest(input)
	_, err := req.Send()
	return err
}
