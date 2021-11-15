package listener

import (
	"context"
)

// Observed Observed
type Observed interface {
	Notify(ctx context.Context)
	AddObserve(ob... Observer)
	RemoveObserve(ob Observer)
}

// Observer Observer
type Observer interface {
	Update(context.Context,Param)
}
// Param Param
type Param interface {
	Eval()
}


