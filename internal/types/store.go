package types

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	ListApplications() ([]*JobApplication, error)
	GetApplication(companyID string) (*JobApplication, error)
	CreateApplication(application *JobApplication) error
	DeleteApplication(companyID string) error
	UpdateApplicationStatus(companyID string, status ApplicationStatus) error
	AddApplicationNote(companyID string, note string) error
	RemoveApplicationNote(companyID string, noteIndex int) error
	ListApplicationNotes(companyID string) ([]string, error)
}

type DB struct {
	Pool	*pgxpool.Pool
}

func (db *DB) ListApplications() ([]*JobApplication, error) {
	rows, err := db.Pool.Query(context.Background(), "select company, role, status, notes from applications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := []*JobApplication{}

	for rows.Next() {
		var company string
		var role JobRole
		var status ApplicationStatus
		var notes []string

		err = rows.Scan(&company, &role, &status, &notes)
		if err != nil {
			return nil, err
		}

		job_application, err := NewJobApplication(company, role, status, notes)
		if err != nil {
			return nil, err
		}

		applications = append(applications, job_application)
	}

	return applications, nil
}

func (db *DB) GetApplication(companyID string) (*JobApplication, error) {
	var role JobRole
	var status ApplicationStatus
	var notes []string
	err := db.Pool.QueryRow(context.Background(), "select role, status, notes from applications where company=$1", companyID).Scan(&role, &status, &notes)
	if err != nil {
		return nil, err
	}

	job_application, err := NewJobApplication(companyID, role, status, notes)
	if err != nil {
		return nil, err
	}

	return job_application, nil
}

func (db *DB) CreateApplication(application *JobApplication) error {
	_, err := db.Pool.Exec(context.Background(), "insert into applications (company, role, status, notes) values ($1, $2, $3, $4)",
		application.GetCompany(),
		application.GetRole(),
		application.GetStatus(),
		application.GetNotes(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteApplication(companyID string) error {
	_, err := db.Pool.Exec(context.Background(), "delete from applications where company=$1", companyID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateApplicationStatus(companyID string, status ApplicationStatus) error {
	if status > MaxStatus {
		return errors.New("invalid status value")
	}

	_, err := db.Pool.Exec(context.Background(), "update applications set status=$1 where company=$2",
		status,
		companyID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddApplicationNote(companyID string, note string) error {
	_, err := db.Pool.Exec(context.Background(), "update applications set notes = array_append(notes, $1) where company=$2",
		note,
		companyID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveApplicationNote(companyID string, noteIndex int) error {
	_, err := db.Pool.Exec(context.Background(), "update applications set notes = array_remove(notes, notes[$1::int]) where company=$2",
		noteIndex,
		companyID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ListApplicationNotes(companyID string) ([]string, error) {
	var notes []string
	err := db.Pool.QueryRow(context.Background(), "select notes from applications where company=$1", companyID).Scan(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}