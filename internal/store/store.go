package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medidew/ApplicationTracker/internal/auth"
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
	CreateUser(email string, username string, argon2auth *auth.Argon2Auth, hashedPassword []byte) error
	GetUserHashedPassword(username string) ([]byte, error)
	GetUserArgon2Auth(username string) (*auth.Argon2Auth, error)
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

func (db *DB) CreateUser(email string, username string, argon2auth *auth.Argon2Auth, hashedPassword []byte) error {
	mem := int(argon2auth.Argon2Memory)
	time := int(argon2auth.Argon2Time)
	threads := int(argon2auth.Argon2Threads)
	salt := argon2auth.Salt

	_, err := db.Pool.Exec(context.Background(), "insert into users (email, username, argon2_memory, argon2_time, argon2_threads, hashed_password, salt) values ($1, $2, $3, $4, $5, $6, $7)",
		email,
		username,
		mem,
		time,
		threads,
		hashedPassword,
		salt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserHashedPassword(username string) ([]byte, error) {
	var hashedPassword []byte
	err := db.Pool.QueryRow(context.Background(), "select hashed_password from users where username=$1", username).Scan(&hashedPassword)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (db *DB) GetUserArgon2Auth(username string) (*auth.Argon2Auth, error) {
	var mem int
	var time int
	var threads int
	var salt []byte

	err := db.Pool.QueryRow(context.Background(), "select argon2_memory, argon2_time, argon2_threads, salt from users where username=$1", username).Scan(&mem, &time, &threads, &salt)
	if err != nil {
		return nil, err
	}

	argon2auth := &auth.Argon2Auth{
		Argon2Memory:  uint32(mem),
		Argon2Time:    uint32(time),
		Argon2Threads: uint8(threads),
		Salt:          salt,
	}

	return argon2auth, nil
}