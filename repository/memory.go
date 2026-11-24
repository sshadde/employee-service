package repository

import (
	"context"
	"sync"

	"example.com/employee-service/employee"
)

type memoryRepo struct {
	mu        sync.RWMutex
	employees map[string]employee.Employee
}

func NewMemoryRepository() employee.Repository {
	return &memoryRepo{
		employees: make(map[string]employee.Employee),
	}
}

func (r *memoryRepo) CreateEmployee(ctx context.Context, e employee.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.employees[e.ID] = e
	return nil
}

func (r *memoryRepo) GetEmployeeByID(ctx context.Context, id string) (employee.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.employees[id]
	if !ok {
		return employee.Employee{}, employee.ErrEmployeeNotFound
	}
	return e, nil
}

func (r *memoryRepo) ChangeEmployeePosition(ctx context.Context, id string, position string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.employees[id]
	if !ok {
		return employee.ErrEmployeeNotFound
	}
	e.Position = position
	r.employees[id] = e
	return nil
}
