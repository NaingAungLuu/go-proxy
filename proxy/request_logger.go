package proxy

import (
	"net/http"
	"time"
)

type RequestLog struct {
	Request http.Request
	Time    time.Duration
}

type RequestLogger interface {
	Log(log RequestLog)
}

func NewRequestLog(request http.Request, time time.Duration) RequestLog {
	return RequestLog{
		Request: request,
		Time:    time,
	}
}
