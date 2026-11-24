package grpctransport

import (
	"context"
	"errors"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example.com/employee-service/employee"
	"example.com/employee-service/pb"
	"example.com/employee-service/transport"
)

type grpcServer struct {
	pb.UnimplementedEmployeeServiceServer
	create         kitgrpc.Handler
	getByID        kitgrpc.Handler
	changePosition kitgrpc.Handler
}

func NewServer(endpoints transport.Endpoints) pb.EmployeeServiceServer {
	return &grpcServer{
		create: kitgrpc.NewServer(
			endpoints.Create,
			decodeCreateRequest,
			encodeCreateResponse,
		),
		getByID: kitgrpc.NewServer(
			endpoints.GetByID,
			decodeGetByIDRequest,
			encodeGetByIDResponse,
		),
		changePosition: kitgrpc.NewServer(
			endpoints.ChangePosition,
			decodeChangePositionRequest,
			encodeChangePositionResponse,
		),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) ChangePosition(ctx context.Context, req *pb.ChangePositionRequest) (*pb.ChangePositionResponse, error) {
	_, resp, err := s.changePosition.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ChangePositionResponse), nil
}

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateRequest)
	e := req.GetEmployee()
	return transport.CreateRequest{
		Employee: employee.Employee{
			FullName: e.GetFullName(),
			Position: e.GetPosition(),
			Salary:   e.GetSalary(),
		},
	}, nil
}

func decodeGetByIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.GetId()}, nil
}

func decodeChangePositionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ChangePositionRequest)
	return transport.ChangePositionRequest{ID: req.GetId(), Position: req.GetPosition()}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.CreateResponse)
	if res.Err != nil {
		return nil, toGRPCError(res.Err)
	}
	return &pb.CreateResponse{Id: res.ID}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.GetByIDResponse)
	if res.Err != nil {
		return nil, toGRPCError(res.Err)
	}
	return &pb.GetByIDResponse{
		Employee: &pb.Employee{
			Id:       res.Employee.ID,
			FullName: res.Employee.FullName,
			Position: res.Employee.Position,
			Salary:   res.Employee.Salary,
			HiredOn:  res.Employee.HiredOn,
		},
	}, nil
}

func encodeChangePositionResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.ChangePositionResponse)
	if res.Err != nil {
		return nil, toGRPCError(res.Err)
	}
	return &pb.ChangePositionResponse{}, nil
}

func toGRPCError(err error) error {
	if errors.Is(err, employee.ErrEmployeeNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
