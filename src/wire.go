package src

import (
	"github.com/google/wire"
	"wfRadius/src/handler"
	"wfRadius/src/root"
	"wfRadius/src/wifilog"
	"wfRadius/src/ws"
)

var Set = wire.NewSet(
	root.Set,
	ws.Set,
	wifilog.NewUploader,
	handler.NewRadiusServer,
)
