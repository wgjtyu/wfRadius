package main

import "go.uber.org/zap"

type badgerLogger struct {
	log *zap.SugaredLogger
}

func (b *badgerLogger) Infof(template string, args ...interface{}) {
	b.log.Infof(template, args...)
}
func (b *badgerLogger) Errorf(template string, args ...interface{}) {
	b.log.Errorf(template, args...)
}

func (b *badgerLogger) Debugf(template string, args ...interface{}) {
	b.log.Debugf(template, args...)
}
func (b *badgerLogger) Warningf(template string, args ...interface{}) {
	b.log.Warnf(template, args...)
}
