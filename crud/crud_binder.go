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
	// post call to create crud entry
	// service.Create
	ht.POST(
		"/v1.0/crud",
		NewCreateHandler(b.service),
		append(opts, NewCreateHandlerOption()...)...,
	)

	ht.GET(
		"/v1.0/crud/:key",
		NewGetHandler(b.service),
		append(opts, NewGetHandlerOption()...)...,
	)

	ht.DELETE(
		"/v1.0/crud/:key",
		NewDeleteHandler(b.service),
		append(opts, NewDeleteHandlerOption()...)...,
	)

	ht.GET(
		"/v1.0/crud",
		NewListHandler(b.service),
		append(opts, NewListHandlerOption()...)...,
	)
}

func (b *Binder) Service() Service { return b.service }

func NewHTTPBinder(logger log.Logger, sampleKey, sampleValue string) (*Binder, error) {
	ss, err := NewCrudService(
		sampleKey, sampleValue,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise CRUD service")
	}

	return &Binder{ss}, err
}
