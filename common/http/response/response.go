package commonhttpresp

import commonerr "waizlytest/common/errors"

type Response struct {
	// this field should be filled when 4xx and 5xx status returned
	Error *commonerr.Error `json:"error"`

	// any object from service should be on this field
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}, err *commonerr.Error) *Response {
	return &Response{
		Data:  data,
		Error: err,
	}
}
