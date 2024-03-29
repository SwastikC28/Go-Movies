package monitor

import "context"

type EventHandler interface {
	Start(context.Context)
	Stop(context.Context)
}

type EventRouteProvider interface {
	RegisterEvents()
}
