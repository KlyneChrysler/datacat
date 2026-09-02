// Package kafkax will own the shared Kafka producer/consumer wrappers:
// connection setup, retry with jittered backoff, DLQ routing, and offset
// commit discipline. Implemented in the Kafka phase; the client library is
// resolved with `go get` at that time (candidate: twmb/franz-go), never
// guessed.
package kafkax
