package middleware

import (
	"github.com/Francisundermoon/commonChainApi/db"
	"github.com/Francisundermoon/commonChainApi/model"
)

type MultiRecord struct {
	Id          int64
	FromAddress string
}

func RecoverData() {
	db := db.InnerDbInit()
	db.Model(&model.ChargeRecord{}).Where("id = ?", 1).Update("is_open", 0)
	d, _ := db.DB()
	d.Close()
}
