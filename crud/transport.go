package crud

import (
	"context"
	"fmt"
	net_http "net/http"

	utils_err "github.com/bhuvankumar123/klg/utils/err"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/transport/http"
)

var (
	errBadRequest     = errors.New("bad request")
	errInternalServer = errors.New("internal server verror")
)

// dto for create request, read from form data
type createRequestData struct {
	key   string
	value string
}

// decoder, reads the http request and gets the required
// request property for the service
func createDecoder(
	cx context.Context, req *net_http.Request,
) (interface{}, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, errors.Wrap(
			errBadRequest, "failed to parse form data: "+err.Error(),
		)
	}

	crd := &createRequestData{
		key:   req.Form.Get("key"),
		value: req.Form.Get("value"),
	}

	fmt.Println("Key:", crd.key, "value:", crd.value)

	return crd, err
}

// we don't need any encoder, as `Create()` doesn't return any value
// if there is no error, we should return 204
func createNoContentEncoder(
	cx context.Context,
	rw net_http.ResponseWriter,
	object interface{},
) (err error) {
	if object != nil {
		// use the default encoder
		return http.NewDefaultJSONEncoder()(cx, rw, object)
	}

	rw.WriteHeader(net_http.StatusNoContent)
	return nil
}

// endpoint handles the call to service
func createEndpoint(svc Service) endpoint.Endpoint {
	return func(cx context.Context, req interface{}) (res interface{}, err error) {
		rq, ok := req.(*createRequestData)
		if !ok {
			return nil, errors.Wrap(
				errInternalServer, "failed to cast object",
			)
		}

		err = svc.Create(cx, rq.key, rq.value)
		return nil, err
	}
}

func NewCreateHandler(service Service) http.Handler {
	return http.Handler(createEndpoint(service))
}

func NewCreateHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(createDecoder),
		http.HandlerWithEncoder(createNoContentEncoder),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

// same as get, delete has very similar handling of request and response.
// only difference is we return 204 instead of any json data

func deleteDecoder(
	cx context.Context, req *net_http.Request,
) (interface{}, error) {
	var (
		params = http.Parameters(req)
		key    = params.ByName("key")
	)

	if key == "" {
		return nil, errors.Wrap(
			errBadRequest, "key missing from url params",
		)
	}

	return key, nil
}

func deleteNoContentEncoder(
	cx context.Context,
	rw net_http.ResponseWriter,
	object interface{},
) (err error) {
	if object != nil {
		return http.NewDefaultJSONEncoder()(cx, rw, object)
	}

	rw.WriteHeader(net_http.StatusNoContent)
	return nil
}

func deleteEndpoint(service Service) endpoint.Endpoint {
	return func(cx context.Context, req interface{}) (res interface{}, err error) {
		rq, ok := req.(string)
		if !ok {
			return nil, errors.Wrap(
				errInternalServer, "failed to cast object",
			)
		}

		err = service.Delete(cx, rq)
		return nil, err
	}
}

func NewDeleteHandler(service Service) http.Handler {
	return http.Handler(deleteEndpoint(service))
}

func NewDeleteHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(deleteDecoder),
		http.HandlerWithEncoder(deleteNoContentEncoder),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

// dto for get is simple string i.e. key to which we will return value

func getDecoder(
	cx context.Context, req *net_http.Request,
) (interface{}, error) {
	var (
		params = http.Parameters(req)
		key    = params.ByName("key")
	)

	if key == "" {
		return nil, errors.Wrap(
			errBadRequest, "key missing from url params",
		)
	}

	return key, nil
}

// we will use default JSON Encoder
// 	http.NewDefaultJSONEncoder()

func getEndpoint(svc Service) endpoint.Endpoint {
	return func(cx context.Context, req interface{}) (res interface{}, err error) {
		rq, ok := req.(string)
		if !ok {
			return nil, errors.Wrap(
				errInternalServer, "failed to cast object",
			)
		}

		return svc.Get(cx, rq)
	}
}

func NewGetHandler(service Service) http.Handler {
	return http.Handler(getEndpoint(service))
}

func NewGetHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(getDecoder),
		http.HandlerWithEncoder(http.NewDefaultJSONEncoder()),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

// unlike any of the APIs that we have, this one doesn't take any input

func listEndpoint(svc Service) endpoint.Endpoint {
	return func(cx context.Context, req interface{}) (res interface{}, err error) {
		return svc.List(cx)
	}
}

func NewListHandler(service Service) http.Handler {
	return http.Handler(listEndpoint(service))
}

func NewListHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(http.NopRequestDecoder()),
		http.HandlerWithEncoder(http.NewDefaultJSONEncoder()),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

func errorWriter(
	cx context.Context, er *utils_err.Error, w net_http.ResponseWriter,
) {
	bt, err := er.JSON()
	if err != nil {
		net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(int(er.Code()))
	w.Write(bt)
}

func errEncoder(
	cx context.Context,
	err error,
	w net_http.ResponseWriter,
) {
	cause := errors.Cause(err)
	switch cause {
	case errBadRequest:
		errorWriter(
			cx,
			utils_err.NewError(err, net_http.StatusBadRequest, "bad request"),
			w,
		)
	case ErrNotFound:
		errorWriter(
			cx, utils_err.NewError(err, net_http.StatusNotFound, "not found"),
			w,
		)
	case ErrEmptyKey:
		errorWriter(
			cx, utils_err.NewError(err, net_http.StatusBadRequest, "Bad Request, key Empty"),
			w,
		)
	case errInternalServer:
		fallthrough
	default:
		errorWriter(
			cx,
			utils_err.NewError(err, net_http.StatusInternalServerError, "internal server error"),
			w,
		)
	}
}
