package comet

import (
	"context"
	"encoding/json"
	"git.internal.yunify.com/qxp/molecule/internal/listener"
	"git.internal.yunify.com/qxp/molecule/internal/mq"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/internal/workpool"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"reflect"
)

const (
	flowTopic = "FlowTrigger"
	commonTopic = "FormData"
	dataModify = "dataModify"
	com = "common"
)

// process 触发流程观察者
type process struct {
	pool        *workpool.WorkPool
	queue       *mq.Queue
}

// NewProcess 创建观察者对象实例
func NewProcess(conf *config.Config, opts ...service.Options) (listener.Observer, error) {
	queue,err := mq.NewMQ(conf)
	if err != nil {
		return nil, err
	}
	p := &process{
		pool:        workpool.New(5, 20, 0),
		queue:       queue,
	}
	return p, nil
}

// Update 更新自身状态
func (p *process) Update(ctx context.Context, param listener.Param) {
	p.pool.Start(ctx).PushTaskFunc(ctx, param, p.Handler)
}

// Event event
type Event struct {
	EventType string        `json:"eventType"`
	EventName string        `json:"eventName"`
	Data      interface{}   `json:"data"`
}

// handler DataModifyTrigger struct
func (p *process) Handler(ctx context.Context, param listener.Param) {
	d := reflect.ValueOf(param).Elem()
	topic := d.FieldByName("Topic").String()
	eventName := d.FieldByName("Event").String()
	if topic == "" {
		topic = commonTopic
	}
	e := &Event{
		EventType: com,
		EventName: eventName,
		Data: d.Interface(),
	}
	paramByte, err := json.Marshal(e)
	if err != nil {
		return
	}
	p.queue.Monitor(ctx)
	p.queue.TriggerASyncProducer(topic,paramByte)
}

