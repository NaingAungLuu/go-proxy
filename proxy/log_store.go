package proxy

import "net/http"

type LogStore interface {
	Save(request *http.Request, response *http.Response)
}

type OnMemoryLogStore struct {
	Logs []string
}

func (s *OnMemoryLogStore) Save(request *http.Request, response *http.Response) {
	s.Logs = append(s.Logs, request.Method+" ", request.URL.Path)
}

func NewLogStore() LogStore {
	return &OnMemoryLogStore{}
}
