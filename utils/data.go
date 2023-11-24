package utils

import (
	"chargeWithdraw/model"
	"strconv"
)

func StrToInt64(str string) (int64, *model.ErrorMessage) {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		model.ErrorHandle(err, 1)
		return 0, &model.ErrorMessage{
			Status: 1,
			Msg:    "参数类型错误",
		}
	}
	return num, nil
}

func StrToFloat64(str string) (float64, *model.ErrorMessage) {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		model.ErrorHandle(err, 1)
		return 0, &model.ErrorMessage{
			Status: 1,
			Msg:    "参数类型错误",
		}
	}
	return num, nil
}
