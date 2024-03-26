package config

import "context"

func StartEventHandler() {
	eventHandler.Start(context.Background())
}

func StopEventHandler() {
	eventHandler.Stop(context.Background())
}
