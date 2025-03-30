package app

import (
	"fmt"
	"reflect"

	"github.com/unbxd/go-base/kit/transport/http"
)

type Binder interface {
	Bind(tt *http.Transport, opts ...http.HandlerOption)
}

func WithHTTPBinder(binder Binder) Option {
	fmt.Println(">> Initialising -- ", reflect.TypeOf(binder))
	return func(a *App) (err error) {
		a.binders = append(a.binders, binder)
		return err
	}
}
