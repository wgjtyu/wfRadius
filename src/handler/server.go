package handler

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"layeh.com/radius"
)

type RadiusServer struct {
	server *radius.PacketServer
	db     *gorm.DB
	logger *zap.Logger
}

func NewRadiusServer(db *gorm.DB, l *zap.Logger) *RadiusServer {
	return &RadiusServer{
		db:     db,
		logger: l.Named("RadiusServer"),
	}
}

func (rs *RadiusServer) Serve() {
	rs.server = &radius.PacketServer{
		Handler:      radius.HandlerFunc(rs.handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	err := rs.server.ListenAndServe()
	if err != nil {
		rs.logger.Error("Serve出错", zap.Error(err))
	} else {
		rs.logger.Debug("Server结束")
	}
}

func (rs *RadiusServer) Shutdown() error {
	return rs.server.Shutdown(context.Background())
}
