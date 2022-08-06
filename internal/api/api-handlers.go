package api

import (
	"net/http"
	"reflect"
)

type api_error struct {
	ErrorMessage string `json:"errorMessage"`
}

func (a api_error) IsEmpty() bool {
	return reflect.DeepEqual(a, api_error{})
}

var pong = func() (int, interface{}, api_error) {
	body := struct {
		Ping string `json:"ping"`
	}{Ping: "pong"}
	return http.StatusOK, body, api_error{}
}
