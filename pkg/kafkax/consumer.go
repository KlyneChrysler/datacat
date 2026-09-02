package kafkax

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Handler processes one record. Handlers must be idempotent: delivery is
// at-least-once (offsets commit only after the poll batch is handled).
type Handler func(ctx context.Context, key, value []byte) error

type Consumer struct {
	client *kgo.Client
	log    *slog.Logger
}

func NewConsumer(brokers, group, topic string, log *slog.Logger) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(splitBrokers(brokers)...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer for %s: %w", topic, err)
	}
	return &Consumer{client: client, log: log}, nil
}

// Consume blocks, invoking handle for each record until ctx is cancelled.
// A record whose handler fails is logged and skipped — never silently
// dropped, never allowed to wedge the partition.
func (c *Consumer) Consume(ctx context.Context, handle Handler) error {
	defer c.client.Close()
	for {
		if err := c.pollOnce(ctx, handle); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

func (c *Consumer) pollOnce(ctx context.Context, handle Handler) error {
	fetches := c.client.PollFetches(ctx)
	if ctx.Err() != nil {
		return nil
	}
	c.logFetchErrors(ctx, fetches)
	fetches.EachRecord(func(r *kgo.Record) { c.handleRecord(ctx, r, handle) })
	return c.commit(ctx)
}

func (c *Consumer) handleRecord(ctx context.Context, r *kgo.Record, handle Handler) {
	if err := handle(ctx, r.Key, r.Value); err != nil {
		c.log.ErrorContext(ctx, "record handling failed; skipping",
			"err", err, "topic", r.Topic, "partition", r.Partition, "offset", r.Offset)
	}
}

func (c *Consumer) commit(ctx context.Context) error {
	if err := c.client.CommitUncommittedOffsets(ctx); err != nil && ctx.Err() == nil {
		return fmt.Errorf("commit offsets: %w", err)
	}
	return nil
}

func (c *Consumer) logFetchErrors(ctx context.Context, fetches kgo.Fetches) {
	for _, fe := range fetches.Errors() {
		c.log.ErrorContext(ctx, "fetch error",
			"err", fe.Err, "topic", fe.Topic, "partition", fe.Partition)
	}
}
