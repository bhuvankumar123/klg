package crud

import (
	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/transport/http"
	"github.com/unbxd/go-base/utils/log"
)

type Binder struct {
	service Service
}

func (b *Binder) Bind(ht *http.Transport, opts ...http.HandlerOption) {
	// Post Call to create log
	ht.POST(
		"/v1.0/logs",
		NewCreateHandler(b.service),
		append(opts, NewCreateHandlerOption()...)...,
	)

	// Get Call to fetch
	ht.GET(
		"/v1.0/logs/:id",
		NewGetHandler(b.service),
		append(opts, NewGetHandlerOption()...)...,
	)

	// Get call
	ht.GET(
		"/v1.0/logs",
		NewListHandler(b.service),
		append(opts, NewListHandlerOption()...)...,
	)

	ht.DELETE(
		"/v1.0/logs",
		NewDeleteHandler(b.service),
		append(opts, NewDeleteHandlerOption()...)...,
	)
}

func (b *Binder) Service() Service { return b.service }

func NewHTTPBinder(logger log.Logger, mongoURI, database string) (*Binder, error) {
	service, err := NewMongoService(mongoURI, database)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize MongoDB service")
	}

	return &Binder{service}, nil
}
