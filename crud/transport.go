package crud

import (
	"context"
	"encoding/json"
	net_http "net/http"

	utils_err "github.com/bhuvankumar123/klg/utils/err"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/transport/http"
)

var (
	errBadRequest     = errors.New("bad request")
	errInternalServer = errors.New("internal server error")
)

type createLogRequest struct {
	Level    string                 `json:"level"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// decoder, reads the http request and gets the required
// request property for the service
func createDecoder(
	ctx context.Context, req *net_http.Request,
) (interface{}, error) {
	var request createLogRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(errBadRequest, "failed to decode request")
	}

	if request.Level == "" || request.Message == "" {
		return nil, errors.Wrap(errBadRequest, "level and message are required")
	}

	// Validate log level
	if err := ValidateLogLevel(request.Level); err != nil {
		return nil, err
	}

	return request, nil
}

// endpoint handles the call to service
func createEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		rq, ok := req.(createLogRequest)
		if !ok {
			return nil, errors.Wrap(errInternalServer, "failed to cast request")
		}

		err = svc.Create(ctx, rq.Level, rq.Message, rq.Metadata)
		if err != nil {
			return nil, err
		}

		// Return success response
		return map[string]interface{}{
			"status":  "success",
			"message": "Log entry created successfully",
			"data": map[string]interface{}{
				"level":    rq.Level,
				"message":  rq.Message,
				"metadata": rq.Metadata,
			},
		}, nil
	}
}

func NewCreateHandler(service Service) http.Handler {
	return http.Handler(createEndpoint(service))
}

func NewCreateHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(createDecoder),
		http.HandlerWithEncoder(http.NewDefaultJSONEncoder()),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

// dto for get is simple string i.e. key to which we will return value

func getDecoder(
	ctx context.Context, req *net_http.Request,
) (interface{}, error) {
	var (
		params = http.Parameters(req)
		id     = params.ByName("id")
	)

	if id == "" {
		return nil, errors.Wrap(errBadRequest, "id missing from url params")
	}

	return id, nil
}

// we will use default JSON Encoder
// 	http.NewDefaultJSONEncoder()

func getEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		id, ok := req.(string)
		if !ok {
			return nil, errors.Wrap(errInternalServer, "failed to cast request")
		}

		return svc.Get(ctx, id)
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

func listDecoder(
	ctx context.Context, req *net_http.Request,
) (interface{}, error) {
	filter := make(map[string]interface{})

	// Get query parameters
	query := req.URL.Query()

	// Add level filter if present
	if level := query.Get("level"); level != "" {
		filter["level"] = level
	}

	// Add message filter if present
	if message := query.Get("message"); message != "" {
		filter["message"] = message
	}

	// Add time range filters if present
	if starttime := query.Get("starttime"); starttime != "" {
		filter["starttime"] = starttime
	}

	if endtime := query.Get("endtime"); endtime != "" {
		filter["endtime"] = endtime
	}

	if recent := query.Get("recent"); recent != "" {
		filter["recent"] = recent
	}

	// Add metadata filters if present
	for key, values := range query {
		if key != "level" && key != "message" && key != "starttime" && key != "endtime" {
			filter["metadata."+key] = values[0]
		}
	}

	return filter, nil
}

func listEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		filter, ok := req.(map[string]interface{})
		if !ok {
			return nil, errors.Wrap(errInternalServer, "failed to cast filter")
		}

		return svc.List(ctx, filter)
	}
}

func NewListHandler(service Service) http.Handler {
	return http.Handler(listEndpoint(service))
}

func NewListHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(listDecoder),
		http.HandlerWithEncoder(http.NewDefaultJSONEncoder()),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

func NewDeleteHandler(service Service) http.Handler {
	return http.Handler(deleteEndpoint(service))
}

func NewDeleteHandlerOption() []http.HandlerOption {
	return []http.HandlerOption{
		http.HandlerWithDecoder(deleteDecoder),
		http.HandlerWithEncoder(http.NewDefaultJSONEncoder()),
		http.HandlerWithErrorEncoder(errEncoder),
	}
}

func deleteDecoder(
	ctx context.Context, req *net_http.Request,
) (interface{}, error) {
	filter := make(map[string]interface{})

	// Get query parameters
	query := req.URL.Query()

	// Add level filter if present
	if id := query.Get("id"); id != "" {
		filter["id"] = id
	}

	if before := query.Get("before"); before != "" {
		filter["before"] = before
	}

	if filter["id"] == "" || filter["before"] == "" {
		return nil, errors.Wrap(errBadRequest, "id and before missing from url params")
	}

	if filter["id"] != "" && filter["before"] != "" {
		return nil, errors.Wrap(errBadRequest, "only one of id or before can be provided")
	}

	return filter, nil
}

func deleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		filter, ok := req.(map[string]interface{})
		if !ok {
			return nil, errors.Wrap(errInternalServer, "failed to cast filter")
		}
		err = svc.Delete(ctx, filter)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status":  "success",
			"message": "Log entries deleted successfully",
		}, nil
	}
}

func errorWriter(
	ctx context.Context, er *utils_err.Error, w net_http.ResponseWriter,
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
	ctx context.Context,
	err error,
	w net_http.ResponseWriter,
) {
	cause := errors.Cause(err)
	switch cause {
	case errBadRequest:
		errorWriter(
			ctx,
			utils_err.NewError(err, net_http.StatusBadRequest, "bad request"),
			w,
		)
	case ErrNotFound:
		errorWriter(
			ctx, utils_err.NewError(err, net_http.StatusNotFound, "not found"),
			w,
		)
	case ErrEmptyKey:
		errorWriter(
			ctx, utils_err.NewError(err, net_http.StatusBadRequest, "Bad Request, required fields missing"),
			w,
		)
	case errInternalServer:
		fallthrough
	default:
		errorWriter(
			ctx,
			utils_err.NewError(err, net_http.StatusInternalServerError, "internal server error"),
			w,
		)
	}
}
