module wfRadius

go 1.14

require (
	github.com/google/wire v0.5.0
	github.com/gorilla/websocket v1.4.2
	github.com/imroc/req v0.3.0
	github.com/jpillora/overseer v1.1.6
	github.com/mitchellh/mapstructure v1.4.2
	github.com/spf13/viper v1.9.0
	github.com/wgjtyu/goutil v0.0.0-00010101000000-000000000000
	github.com/wgjtyu/logMansion v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.19.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.16
	layeh.com/radius v0.0.0-20201203135236-838e26d0c9be
	moul.io/zapgorm2 v1.1.0
)

replace github.com/wgjtyu/logMansion => ../logMansion

replace github.com/wgjtyu/goutil => ../goutil
