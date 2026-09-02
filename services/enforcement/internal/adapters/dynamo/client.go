// Package dynamo holds the DynamoDB decision store adapter.
package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient builds the client, a non empty endpoint targets DynamoDB Local.
func NewClient(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	if endpoint == "" {
		return dynamodb.NewFromConfig(awsCfg), nil
	}

	return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) { o.BaseEndpoint = aws.String(endpoint) }), nil
}
