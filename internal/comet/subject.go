package comet

import (
	"context"
	"git.internal.yunify.com/qxp/molecule/internal/listener"
)

func (comet *comet) AddObserve(ob... listener.Observer) {
	comet.observers = append(comet.observers, ob...)
}

func (comet *comet) RemoveObserve(ob listener.Observer) {
	for i, s := range comet.observers {
		if s == ob {
			comet.observers = append(comet.observers[:i], comet.observers[i+1:]...)
		}
	}
}

func (comet *comet) Notify(ctx context.Context,param listener.Param) {
	for _, s := range comet.observers {
		s.Update(ctx,param)
	}
}

func (p *processServer) AddObserve(ob... listener.Observer) {
	p.observers = append(p.observers, ob...)
}

func (p *processServer) Notify(ctx context.Context,param listener.Param) {
	for _, s := range p.observers {
		s.Update(ctx,param)
	}
}

// DataModifyTrigger DataModifyTrigger event
type DataModifyTrigger struct {
	TableID string          `json:"tableID"`
	Entity  *[]interface{}  `json:"entity"`
	Method  string          `json:"method"`
	UserID  string          `json:"userID"`
	Topic   string          `json:"topic"`
	Event   string          `json:"event"`
}
// Eval 暂空
func (d *DataModifyTrigger)Eval() {}
