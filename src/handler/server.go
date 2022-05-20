package handler

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"layeh.com/radius"
	"sync"
	"wfRadius/src/wifilog"
)

type RadiusServer struct {
	server   *radius.PacketServer
	db       *gorm.DB
	Uploader *wifilog.Uploader
	logger   *zap.Logger
}

func NewRadiusServer(db *gorm.DB, l *zap.Logger, u *wifilog.Uploader) *RadiusServer {
	return &RadiusServer{
		db:       db,
		Uploader: u,
		logger:   l.Named("RadiusServer"),
	}
}

func (rs *RadiusServer) Serve(wg *sync.WaitGroup) {
	defer wg.Done()
	rs.server = &radius.PacketServer{
		Handler:      radius.HandlerFunc(rs.handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	err := rs.server.ListenAndServe()
	if err != nil {
		rs.logger.Debug("Serve出错", zap.Error(err))
	} else {
		rs.logger.Debug("Server结束")
	}
}

func (rs *RadiusServer) Shutdown() error {
	return rs.server.Shutdown(context.Background())
}
