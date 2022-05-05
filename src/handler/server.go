package handler

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"layeh.com/radius"
	"log"
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

	if err := rs.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (rs *RadiusServer) Shutdown() error {
	return rs.server.Shutdown(context.Background())
}
