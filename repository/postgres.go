package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/employee-service/employee"
)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, dsn string) (employee.Repository, func(), error) {
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 30; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				break
			}
			pool.Close()
		}
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	if err != nil {
		return nil, nil, err
	}

	const schema = `CREATE TABLE IF NOT EXISTS employees (
		id TEXT PRIMARY KEY,
		full_name TEXT NOT NULL,
		job_position TEXT NOT NULL,
		salary DOUBLE PRECISION NOT NULL,
		hired_on BIGINT NOT NULL
	)`
	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("create schema: %w", err)
	}

	return &postgresRepo{pool: pool}, pool.Close, nil
}

func (r *postgresRepo) CreateEmployee(ctx context.Context, e employee.Employee) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO employees (id, full_name, job_position, salary, hired_on) VALUES ($1, $2, $3, $4, $5)`,
		e.ID, e.FullName, e.Position, e.Salary, e.HiredOn)
	return err
}

func (r *postgresRepo) GetEmployeeByID(ctx context.Context, id string) (employee.Employee, error) {
	var e employee.Employee
	err := r.pool.QueryRow(ctx,
		`SELECT id, full_name, job_position, salary, hired_on FROM employees WHERE id = $1`, id).
		Scan(&e.ID, &e.FullName, &e.Position, &e.Salary, &e.HiredOn)
	if errors.Is(err, pgx.ErrNoRows) {
		return employee.Employee{}, employee.ErrEmployeeNotFound
	}
	if err != nil {
		return employee.Employee{}, err
	}
	return e, nil
}

func (r *postgresRepo) ChangeEmployeePosition(ctx context.Context, id string, position string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE employees SET job_position = $1 WHERE id = $2`, position, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return employee.ErrEmployeeNotFound
	}
	return nil
}
