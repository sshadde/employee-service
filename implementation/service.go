package implementation

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"example.com/employee-service/employee"
)

type service struct {
	repository employee.Repository
	logger     log.Logger
}

func NewService(repo employee.Repository, logger log.Logger) employee.Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, e employee.Employee) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	e.ID = id
	e.HiredOn = time.Now().Unix()

	if err := s.repository.CreateEmployee(ctx, e); err != nil {
		level.Error(logger).Log("err", err)
		return "", employee.ErrCmdRepository
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (employee.Employee, error) {
	logger := log.With(s.logger, "method", "GetByID")

	e, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if errors.Is(err, employee.ErrEmployeeNotFound) {
			return e, employee.ErrEmployeeNotFound
		}
		return e, employee.ErrQueryRepository
	}
	return e, nil
}

func (s *service) ChangePosition(ctx context.Context, id string, position string) error {
	logger := log.With(s.logger, "method", "ChangePosition")

	if err := s.repository.ChangeEmployeePosition(ctx, id, position); err != nil {
		level.Error(logger).Log("err", err)
		if errors.Is(err, employee.ErrEmployeeNotFound) {
			return employee.ErrEmployeeNotFound
		}
		return employee.ErrCmdRepository
	}
	return nil
}
