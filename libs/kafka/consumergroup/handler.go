package consumergroup

import (
	"log/slog"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*Handler)(nil)

type Msg struct {
	Topic     string            `json:"topic"`
	Partition int32             `json:"partition"`
	Offset    int64             `json:"offset"`
	Key       []byte            `json:"key"`
	Payload   []byte            `json:"payload"`
	Headers   map[string][]byte `json:"headers"`
}

type TopicHandler interface {
	Handle(msg Msg) error
}

type Handler struct {
	topicHandlers map[string]TopicHandler
}

func NewConsumerGroupHandler(
	topicHandlers map[string]TopicHandler,
) *Handler {
	return &Handler{
		topicHandlers: topicHandlers,
	}
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := convertMsg(message)
			if handler, exists := h.topicHandlers[message.Topic]; exists {
				if err := handler.Handle(msg); err != nil {
					slog.Error("Error handling message: %v", err)
				}
			} else {
				slog.Error("No handler for message claimed from topic", slog.String("topic", message.Topic))
				break
			}

			session.MarkMessage(message, "")
			session.Commit()

		case <-session.Context().Done():
			return nil
		}
	}
}

func convertMsg(in *sarama.ConsumerMessage) Msg {
	headers := make(map[string][]byte)
	for _, header := range in.Headers {
		headers[string(header.Key)] = header.Value
	}
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       in.Key,
		Payload:   in.Value,
		Headers:   headers,
	}
}
