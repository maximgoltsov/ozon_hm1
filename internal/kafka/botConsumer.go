package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"

	logger "github.com/maximgoltsov/botproject/internal/pkg/logger"
)

type BotConsumer struct {
}

func (c *BotConsumer) Setup(session sarama.ConsumerGroupSession) error {
	logger.Logger.Info("BotConsumer setup", session.Claims())
	return nil
}

func (c *BotConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *BotConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		logger.Logger.Info("%v - %v", claim.Topic(), claim.Partition())
		select {
		case <-session.Context().Done():
			logger.Logger.Info("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				logger.Logger.Info("Data channel closed")
				return nil
			}
			message := &Message{}
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				logger.Logger.Error("cant unmarshal", err)
			}

			logger.Logger.Info("get message [%v] from [%v]", message.Action, message.Service)

			session.MarkMessage(msg, "")
		}
	}
}
