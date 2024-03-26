package publisher

import "context"

type Dispatcher interface {
	Publish(corelationId string, topic string, payload interface{})
	Stop(context context.Context)
}
