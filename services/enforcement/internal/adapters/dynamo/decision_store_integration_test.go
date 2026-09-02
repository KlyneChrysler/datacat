//go:build integration

// Integration tests need docker compose up -d dynamodb, run with -tags integration.
package dynamo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

const itTable = "it-datacat-decisions"

func itStore(t *testing.T) *DecisionStore {
	t.Helper()
	t.Setenv("AWS_ACCESS_KEY_ID", "local")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "local")
	t.Setenv("AWS_REGION", "us-east-1")

	ctx := context.Background()
	client, err := NewClient(ctx, "http://localhost:8000")
	if err != nil {
		t.Fatal(err)
	}
	createTableIfMissing(t, ctx, client)
	return NewDecisionStore(client, itTable, time.Hour)
}

func createTableIfMissing(t *testing.T, ctx context.Context, client *dynamodb.Client) {
	t.Helper()
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(itTable),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("session_id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("session_id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	var inUse *types.ResourceInUseException
	if err != nil && !errors.As(err, &inUse) {
		t.Fatalf("create table: %v", err)
	}
}

func TestSaveAndFindBySession(t *testing.T) {
	store := itStore(t)
	decision := policy.Decision{SessionID: "it-s1", Class: policy.Abusive, Action: policy.Block}

	if err := store.Save(context.Background(), decision); err != nil {
		t.Fatal(err)
	}
	found, err := store.FindBySession(context.Background(), "it-s1")

	if err != nil {
		t.Fatal(err)
	}
	if found != decision {
		t.Errorf("found = %+v, want %+v", found, decision)
	}
}

func TestFindUnknownSessionReturnsNotFound(t *testing.T) {
	store := itStore(t)

	_, err := store.FindBySession(context.Background(), "it-nobody")

	if !errors.Is(err, policy.ErrDecisionNotFound) {
		t.Fatalf("err = %v, want ErrDecisionNotFound", err)
	}
}

func TestSaveOverwritesPreviousDecision(t *testing.T) {
	store := itStore(t)
	first := policy.Decision{SessionID: "it-s2", Class: policy.Unverified, Action: policy.Challenge}
	second := policy.Decision{SessionID: "it-s2", Class: policy.Abusive, Action: policy.Block}

	if err := store.Save(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	found, err := store.FindBySession(context.Background(), "it-s2")

	if err != nil {
		t.Fatal(err)
	}
	if found.Action != policy.Block {
		t.Errorf("action = %s, want block (latest wins)", found.Action)
	}
}
