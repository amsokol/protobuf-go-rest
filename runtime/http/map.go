package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/amsokol/protobuf-rest/runtime"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request)

type Paths map[runtime.Path]Handler

type Methods map[string]Paths

type Map struct {
	Methods Methods // HTTP method -> path map
}

func (m *Map) Add(method string, template string, handler Handler) error {
	pp, ok := m.Methods[method]
	if !ok {
		pp = make(Paths)
		m.Methods[method] = pp
	}

	p, err := runtime.NewPath(template)
	if err != nil {
		return fmt.Errorf("add path template for '%s': %w", method, err)
	}

	pp[p] = handler

	return nil
}

func (m *Map) Match(method string, urlPath string) (Handler, runtime.Values) {
	pp, ok := m.Methods[method]
	if !ok {
		return nil, nil
	}

	sp := strings.Split(strings.Trim(urlPath, "/"), "/")

	for p, h := range pp {
		v := make(runtime.Values)
		if ok := p.Match(sp, v); ok {
			return h, v
		}
	}

	return nil, nil
}

func NewMap() Map {
	return Map{Methods: make(Methods)}
}
