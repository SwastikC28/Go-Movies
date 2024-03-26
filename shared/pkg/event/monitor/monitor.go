package monitor

// EventInfo (used to emit and receive events)
// EventHandler Interface
// GenericEventHandler to implement EventHandler
// AddRoute method to bind event to a handler

type EventInfo struct {
	Payload interface{}
}
