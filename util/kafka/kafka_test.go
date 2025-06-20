package kafka

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestWriteKafka(t *testing.T) {
	WriteKafka(ctx)
}

func TestReadKafka(t *testing.T) {
	go ListenSignal()
	ReadKafka(ctx)
}
