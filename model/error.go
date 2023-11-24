package model

import (
	"fmt"
	"log"
)

type ErrorMessage struct {
	Status int64
	Msg    string
}

var (
	NetworkError = &ErrorMessage{
		Status: 0,
		Msg:    "网络错误",
	}
	ParamError = &ErrorMessage{
		Status: 0,
		Msg:    "参数错误",
	}
)

func ErrorHandle(err error, index int) *ErrorMessage {
	fmt.Println(err)
	log.Println(err)
	switch index {
	case 0:
		return &ErrorMessage{Msg: fmt.Sprint(err)}
	case 1:
		return &ErrorMessage{Msg: fmt.Sprint(err)}
	}
	return nil
}
