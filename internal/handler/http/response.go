package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/s-gurman/user-segmentation/pkg/logger"
)

// HTTP 200 response with some value
type successResponse struct {
	Value interface{} `json:"result"`
}

// HTTP 4xx/5xx response with detailed error
type failedResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"error"`
	err  error  `json:"-"`
}

// Logs performed actions and writes HTTP response.
func writeAndLogValue(w http.ResponseWriter, resp successResponse, l logger.Logger, log string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	l.Info(log)

	code := http.StatusOK
	data, err := json.Marshal(resp)
	if err != nil {
		l.Errorf("httpapi - json encode err: %s", err)
		data = []byte(`{"error":"internal error","code":500}`)
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	if _, err = w.Write(data); err != nil {
		l.Errorf("httpapi - response write err: %s", err)
	}
}

// Logs input error and writes HTTP response.
func writeAndLogError(w http.ResponseWriter, resp failedResponse, l logger.Logger, log string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	l.Errorw(
		log,
		"code", resp.Code,
		"error", resp.err,
	)

	data, err := json.Marshal(resp)
	if err != nil {
		l.Errorf("httpapi - json encode err: %s", err)
		data = []byte(`{"error":"internal error","code":500}`)
		resp.Code = http.StatusInternalServerError
	}

	w.WriteHeader(resp.Code)
	if _, err = w.Write(data); err != nil {
		l.Errorf("httpapi - response write err: %s", err)
	}
}
