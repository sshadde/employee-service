package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/employee-service/employee"
	"example.com/employee-service/transport"
)

func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/employees").Handler(kithttp.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/employees/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PATCH").Path("/employees/{id}/position").Handler(kithttp.NewServer(
		endpoints.ChangePosition,
		decodeChangePositionRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Employee); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeChangePositionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var body struct {
		Position string `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transport.ChangePositionRequest{ID: vars["id"], Position: body.Position}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	case errors.Is(err, employee.ErrEmployeeNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
