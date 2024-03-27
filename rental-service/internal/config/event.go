package config

import (
	"context"
	"shared/pkg/event/monitor"
)

func StartEventHandler() {
	eventHandler.(monitor.EventHandler).Start(context.Background())
}

func StopEventHandler() {
	eventHandler.(monitor.EventHandler).Stop(context.Background())
}
