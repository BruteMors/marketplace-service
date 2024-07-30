package producer

import (
	"fmt"
	"time"

	"github.com/BruteMors/marketplace-service/libs/kafka"
	"github.com/IBM/sarama"
)

type SyncProducer struct {
	syncProducer sarama.SyncProducer
}

func NewSyncProducer(conf kafka.Config, opts ...Option) (*SyncProducer, error) {
	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return &SyncProducer{syncProducer: syncProducer}, nil
}

func (p *SyncProducer) SendMessage(
	topicName string,
	key []byte,
	message []byte,
	headers map[string]string,
) (partition int32, offset int64, err error) {
	kafkaHeaders := make([]sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, sarama.RecordHeader{Key: []byte(k), Value: []byte(v)})
	}

	msg := &sarama.ProducerMessage{
		Topic:     topicName,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(message),
		Headers:   kafkaHeaders,
		Timestamp: time.Now().UTC(),
	}

	partition, offset, err = p.syncProducer.SendMessage(msg)
	if err != nil {
		return partition, offset, fmt.Errorf("SendMessage failed: %w", err)
	}

	return partition, offset, nil
}

func (p *SyncProducer) Close() error {
	if p.syncProducer == nil {
		return nil
	}
	err := p.syncProducer.Close()
	if err != nil {
		return fmt.Errorf("failed to close sync producer: %w", err)
	}

	return nil
}
