package contender

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

// Contender is the model for the tattoo options
type Contender struct {
	Name        string
	Description string
	SVG         []byte
	Wins        int
	Losses      int
	Score       int
}

// Matchup is the model for the head-to-head records
// between contenders
type Matchup struct {
	Contender1     string
	Contender2     string
	Contender1Wins int
	Contender2Wins int
}

// ToAttributeMap encodes the values of a contender into the map format
// that dynamo expects
func (c Contender) ToAttributeMap() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Name":        stringToAttributeValue(c.Name),
		"Description": stringToAttributeValue(c.Description),
		"SVG":         bytesToAttributeValue(c.SVG),
		"Wins":        intToAttributeValue(c.Wins),
		"Losses":      intToAttributeValue(c.Wins),
		"Score":       intToAttributeValue(c.Wins),
	}
}

// FromAttributeMap tries to decode a Contender from a dynamo response
func FromAttributeMap(aMap map[string]dynamodb.AttributeValue) (*Contender, error) {
	wins, err := getInt(aMap["Wins"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Wins attribute")
	}
	losses, err := getInt(aMap["Losses"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Losses attribute")
	}
	score, err := getInt(aMap["Score"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Score attribute")
	}

	return &Contender{
		Name:        getString(aMap["Name"]),
		Description: getString(aMap["Description"]),
		SVG:         getBytes(aMap["SVG"]),
		Wins:        wins,
		Losses:      losses,
		Score:       score,
	}, nil
}

func stringToAttributeValue(s string) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{S: aws.String(s)}
}

func intToAttributeValue(n int) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", n))}
}

func bytesToAttributeValue(b []byte) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{B: b}
}

func getString(a dynamodb.AttributeValue) string {
	return *a.S
}

func getInt(a dynamodb.AttributeValue) (int, error) {
	return strconv.Atoi(*a.N)
}

func getBytes(a dynamodb.AttributeValue) []byte {
	return a.B
}
