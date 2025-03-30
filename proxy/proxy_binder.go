package proxy

import (
	net_http "net/http"

	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/endpoint"
	"github.com/unbxd/go-base/kit/transport/http"
	"github.com/unbxd/go-base/kit/transport/http/proxy"
	"github.com/unbxd/go-base/utils/log"
)

type Binder struct {
	logger log.Logger
	fn     endpoint.Endpoint
}

func (b *Binder) Bind(ht *http.Transport, opts ...http.HandlerOption) {
	handler := http.NewHandler(
		http.Handler(b.fn),
		append(
			opts,
			[]http.HandlerOption{
				http.HandlerWithFilter(http.PanicRecovery(b.logger)),
			}...,
		)...,
	)

	ht.Mux().Handler(net_http.MethodGet, "/*", handler)
}

func NewProxyBinder(
	logger log.Logger,
	downstream string,
	options ...proxy.ProxyOption,
) (*Binder, error) {
	fn, err := proxy.NewProxyEndpoint(logger, downstream, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy downstream")
	}

	return &Binder{logger, fn}, nil
}
