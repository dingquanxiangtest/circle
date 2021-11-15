package mq

import (
	"context"
	"git.internal.yunify.com/qxp/misc/kafka"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/Shopify/sarama"
	"sync/atomic"
)

// Queue Queue
type Queue struct {
	p   sarama.AsyncProducer
	sd  int32
}

// NewMQ NewMQ
func NewMQ(conf *config.Config,) (*Queue,error) {
	conf.Kafka.Sarama.Version = sarama.V2_0_0_0
	conf.Kafka.Sarama.Producer.Return.Successes = true
	producer, err := kafka.NewAsyncProducer(conf.Kafka)
	if err != nil {
		return nil,err
	}
	return &Queue{
		p : producer,
		sd: 0,
	},nil
}

// TriggerASyncProducer Flow trigger mq producer
func (q *Queue)TriggerASyncProducer(topic string,msg []byte) {
	input := q.p.Input()
	input <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
}

// Monitor Monitor
func (q *Queue)Monitor(ctx context.Context)  {
	sd := atomic.CompareAndSwapInt32(&q.sd,0,1)
	if sd {
		logger.Logger.Info("开启mq发送返回监视", logger.STDRequestID(ctx))
		go startMonitor(ctx,q.p)
	}
}

func startMonitor(ctx context.Context, producer sarama.AsyncProducer)  {
	for{
		select {
		case <-producer.Successes():
			// 成功(如果不设置config.Producer.Return.Successes = true，不会返回)
			logger.Logger.Info("消息发送mq成功", logger.STDRequestID(ctx))
		case <-producer.Errors():
			// 发送失败，增加重试
			logger.Logger.Error("消息发送mq失败", logger.STDRequestID(ctx))
		case <-ctx.Done():
			return
		}
	}
}
