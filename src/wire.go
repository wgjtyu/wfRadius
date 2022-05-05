package src

import (
	"github.com/google/wire"
	"wfRadius/src/handler"
	"wfRadius/src/root"
	"wfRadius/ws"
)

var Set = wire.NewSet(
	root.Set,
	ws.Set,
	handler.NewRadiusServer,
)
