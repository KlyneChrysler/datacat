// Package dynamo is the DynamoDB DecisionStore adapter. It replaces the
// in-memory adapter behind the same port, which is what makes enforcement
// horizontally scalable. Storage shape in record.go, conversion in codec.go
// (file taxonomy, standards rule 2).
package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient builds the DynamoDB client. A non-empty endpoint overrides the
// target (DynamoDB Local); empty uses real AWS resolution — region and
// credentials come from the standard AWS environment (factor III).
func NewClient(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	if endpoint == "" {
		return dynamodb.NewFromConfig(awsCfg), nil
	}
	return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	}), nil
}
