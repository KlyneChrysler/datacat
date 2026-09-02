package dynamo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

// DecisionStore implements ports.DecisionStore against DynamoDB.
// Both operations are O(1) point reads/writes on the partition key.
type DecisionStore struct {
	client *dynamodb.Client
	table  string
	ttl    time.Duration
}

var _ ports.DecisionStore = (*DecisionStore)(nil)

func NewDecisionStore(client *dynamodb.Client, table string, ttl time.Duration) *DecisionStore {
	return &DecisionStore{client: client, table: table, ttl: ttl}
}

func (s *DecisionStore) Save(ctx context.Context, d domain.Decision) error {
	item, err := attributevalue.MarshalMap(encodeDecision(d, time.Now(), s.ttl))
	if err != nil {
		return fmt.Errorf("marshal decision: %w", err)
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamo put %s: %w", s.table, err)
	}
	return nil
}

func (s *DecisionStore) FindBySession(ctx context.Context, sessionID string) (domain.Decision, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key:       sessionKey(sessionID),
	})
	if err != nil {
		return domain.Decision{}, fmt.Errorf("dynamo get %s: %w", s.table, err)
	}
	if len(out.Item) == 0 {
		return domain.Decision{}, domain.ErrDecisionNotFound
	}
	return s.unmarshalLive(out.Item)
}

func (s *DecisionStore) unmarshalLive(item map[string]types.AttributeValue) (domain.Decision, error) {
	var rec decisionRecord
	if err := attributevalue.UnmarshalMap(item, &rec); err != nil {
		return domain.Decision{}, fmt.Errorf("unmarshal decision: %w", err)
	}
	if expired(rec, time.Now()) {
		return domain.Decision{}, domain.ErrDecisionNotFound
	}
	return decodeDecision(rec), nil
}

func sessionKey(sessionID string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"session_id": &types.AttributeValueMemberS{Value: sessionID},
	}
}
