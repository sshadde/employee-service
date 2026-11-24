package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/employee-service/employee"
)

type CreateRequest struct {
	Employee employee.Employee
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

func (r CreateResponse) Failed() error { return r.Err }

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Employee employee.Employee `json:"employee"`
	Err      error             `json:"-"`
}

func (r GetByIDResponse) Failed() error { return r.Err }

type ChangePositionRequest struct {
	ID       string
	Position string
}

type ChangePositionResponse struct {
	Err error `json:"-"`
}

func (r ChangePositionResponse) Failed() error { return r.Err }

type Endpoints struct {
	Create         endpoint.Endpoint
	GetByID        endpoint.Endpoint
	ChangePosition endpoint.Endpoint
}

func MakeEndpoints(s employee.Service) Endpoints {
	return Endpoints{
		Create:         makeCreateEndpoint(s),
		GetByID:        makeGetByIDEndpoint(s),
		ChangePosition: makeChangePositionEndpoint(s),
	}
}

func makeCreateEndpoint(s employee.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Employee)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s employee.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		e, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Employee: e, Err: err}, nil
	}
}

func makeChangePositionEndpoint(s employee.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangePositionRequest)
		err := s.ChangePosition(ctx, req.ID, req.Position)
		return ChangePositionResponse{Err: err}, nil
	}
}
