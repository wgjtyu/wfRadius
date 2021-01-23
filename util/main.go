package util

import (
	"bytes"
	"encoding/json"
	"github.com/dgraph-io/badger/v2"
	"go.uber.org/zap"
	"strconv"
	"wfRadius/model"
	"wfRadius/storage"
)

func SaveCodes(codes []model.MWifiCode) error {
	if len(codes) > 0 {
		err := storage.BadgerDB.Update(func(txn *badger.Txn) error {
			for _, c := range codes { // FIXME 保存失败时要rollback
				jc, err := json.Marshal(c)
				if err != nil {
					zap.L().Warn("WifiCodes数据转换成json出错", zap.Error(err))
					return err
				}
				var key bytes.Buffer
				key.WriteString("ID")
				key.WriteString(strconv.FormatUint(c.UserID, 10))
				err = txn.Set(key.Bytes(), jc)
				if err != nil {
					zap.L().Warn("数据库保存WifiCodes数据出错", zap.Error(err))
					return err
				}
			}
			return nil
		})
		return err
	}
	return nil
}
func UpdateCode(code *model.MWifiCode) error {
	err := storage.BadgerDB.Update(func(txn *badger.Txn) error {
		jc, err := json.Marshal(code)
		if err != nil {
			zap.S().Debugf("addCode: WifiCode转换成json出错: %s", err.Error())
			return err
		}
		var key bytes.Buffer
		key.WriteString("ID")
		key.WriteString(strconv.FormatUint(code.UserID, 10))
		err = txn.Set(key.Bytes(), jc)
		if err != nil {
			zap.L().Warn("数据库保存WifiCodes数据出错", zap.String("error", err.Error()))
			return err
		}
		return nil
	})
	return err
}
