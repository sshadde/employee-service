package repository

import (
	"context"
	"errors"
	"testing"

	"example.com/employee-service/employee"
)

func TestCreateAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
	e := employee.Employee{ID: "1", FullName: "Ivan", Position: "Junior", Salary: 100}

	if err := repo.CreateEmployee(ctx, e); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetEmployeeByID(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	if got != e {
		t.Fatalf("got %+v, want %+v", got, e)
	}
}

func TestChangePosition(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
	e := employee.Employee{ID: "1", FullName: "Ivan", Position: "Junior"}

	if err := repo.CreateEmployee(ctx, e); err != nil {
		t.Fatal(err)
	}
	if err := repo.ChangeEmployeePosition(ctx, "1", "Middle"); err != nil {
		t.Fatal(err)
	}
	got, _ := repo.GetEmployeeByID(ctx, "1")
	if got.Position != "Middle" {
		t.Fatalf("position = %q, want Middle", got.Position)
	}
}

func TestNotFound(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.GetEmployeeByID(ctx, "missing")
	if !errors.Is(err, employee.ErrEmployeeNotFound) {
		t.Fatalf("got %v, want ErrEmployeeNotFound", err)
	}
	err = repo.ChangeEmployeePosition(ctx, "missing", "X")
	if !errors.Is(err, employee.ErrEmployeeNotFound) {
		t.Fatalf("got %v, want ErrEmployeeNotFound", err)
	}
}
