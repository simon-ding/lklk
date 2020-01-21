package controller

import "fmt"

type Response struct {
	ErrorCode int         `json:"error_code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
}

func ErrorReturn(msg interface{}) *Response {
	return &Response{
		Msg:       fmt.Sprintf("%v", msg),
		ErrorCode: 1,
	}
}

func SuccessReturn(data ...interface{}) *Response {
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	return &Response{
		Msg:       "success",
		ErrorCode: 0,
		Data:      d,
	}
}
