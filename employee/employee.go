package employee

import "context"

type Employee struct {
	ID       string  `json:"id,omitempty"`
	FullName string  `json:"full_name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
	HiredOn  int64   `json:"hired_on,omitempty"`
}

type Service interface {
	Create(ctx context.Context, employee Employee) (string, error)
	GetByID(ctx context.Context, id string) (Employee, error)
	ChangePosition(ctx context.Context, id string, position string) error
}

type Repository interface {
	CreateEmployee(ctx context.Context, employee Employee) error
	GetEmployeeByID(ctx context.Context, id string) (Employee, error)
	ChangeEmployeePosition(ctx context.Context, id string, position string) error
}
