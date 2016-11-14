package engine

import (
	"encoding/json"
	"net/http"
)

type NightcErr struct {
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
	HttpCode int         `json:"-"`
}

func (e *NightcErr) Error() (string) {
	es, _ := json.Marshal(e)
	return string(es)
}

func NewNightcErr(httpCode, code int, data interface{}) *NightcErr {
	return &NightcErr{
		HttpCode: httpCode,
		Code:     code,
		Data:     data,
	}
}

var (
	ParamError  = NewNightcErr(http.StatusForbidden, 1, "param err")
	ServerError = NewNightcErr(http.StatusInternalServerError, 1, "server err")
	RecodeNotFound = NewNightcErr(http.StatusNotFound, 1, "recode not found")
)
