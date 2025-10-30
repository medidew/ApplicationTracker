package store

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medidew/ApplicationTracker/internal/auth"
)

type Store interface {
	ListApplications(username string) ([]*JobApplication, error)
	GetApplication(username string, companyID string) (*JobApplication, error)
	CreateApplication(username string, application *JobApplication) error
	DeleteApplication(username string, companyID string) error
	UpdateApplicationStatus(username string, companyID string, status ApplicationStatus) error
	AddApplicationNote(username string, companyID string, note string) error
	RemoveApplicationNote(username string, companyID string, noteIndex int) error
	ListApplicationNotes(username string, companyID string) ([]string, error)

	CreateUser(email string, username string, argon2auth *auth.Argon2Auth, hashedPassword []byte) error
	GetUserHashedPassword(username string) ([]byte, error)
	GetUserArgon2Auth(username string) (*auth.Argon2Auth, error)
}

type DB struct {
	Pool	*pgxpool.Pool
}

func (db *DB) ListApplications(username string) ([]*JobApplication, error) {
	rows, err := db.Pool.Query(context.Background(), "select company, role, status, notes from applications where username=$1", username)
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

func (db *DB) GetApplication(username string, companyID string) (*JobApplication, error) {
	var role JobRole
	var status ApplicationStatus
	var notes []string
	err := db.Pool.QueryRow(context.Background(), "select role, status, notes from applications where company=$1 and username=$2", companyID, username).Scan(&role, &status, &notes)
	if err != nil {
		return nil, err
	}

	job_application, err := NewJobApplication(companyID, role, status, notes)
	if err != nil {
		return nil, err
	}

	return job_application, nil
}

func (db *DB) CreateApplication(username string, application *JobApplication) error {
	_, err := db.Pool.Exec(context.Background(), "insert into applications (company, role, status, notes, username) values ($1, $2, $3, $4, $5)",
		application.GetCompany(),
		application.GetRole(),
		application.GetStatus(),
		application.GetNotes(),
		username,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteApplication(username string, companyID string) error {
	_, err := db.Pool.Exec(context.Background(), "delete from applications where company=$1 and username=$2", companyID, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateApplicationStatus(username string, companyID string, status ApplicationStatus) error {
	if status > MaxStatus {
		return errors.New("invalid status value")
	}

	_, err := db.Pool.Exec(context.Background(), "update applications set status=$1 where company=$2 and username=$3",
		status,
		companyID,
		username,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddApplicationNote(username string, companyID string, note string) error {
	_, err := db.Pool.Exec(context.Background(), "update applications set notes = array_append(notes, $1) where company=$2 and username=$3",
		note,
		companyID,
		username,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveApplicationNote(username string, companyID string, noteIndex int) error {
	_, err := db.Pool.Exec(context.Background(), "update applications set notes = array_remove(notes, notes[$1::int]) where company=$2 and username=$3",
		noteIndex,
		companyID,
		username,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ListApplicationNotes(username string, companyID string) ([]string, error) {
	var notes []string
	err := db.Pool.QueryRow(context.Background(), "select notes from applications where company=$1 and username=$3", companyID, username).Scan(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (db *DB) CreateUser(email string, username string, argon2auth *auth.Argon2Auth, hashedPassword []byte) error {
	mem := int(argon2auth.Argon2Memory)
	time := int(argon2auth.Argon2Time)
	threads := int(argon2auth.Argon2Threads)
	salt := base64.RawStdEncoding.EncodeToString(argon2auth.Salt)

	_, err := db.Pool.Exec(context.Background(), "insert into users (email, username, argon2_memory, argon2_time, argon2_threads, hashed_password, salt) values ($1, $2, $3, $4, $5, $6, $7)",
		email,
		username,
		mem,
		time,
		threads,
		base64.RawStdEncoding.EncodeToString(hashedPassword),
		salt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserHashedPassword(username string) ([]byte, error) {
	var hashed_password string
	err := db.Pool.QueryRow(context.Background(), "select hashed_password from users where username=$1", username).Scan(&hashed_password)
	if err != nil {
		return nil, err
	}

	decoded_password, err := base64.RawStdEncoding.DecodeString(hashed_password)
	if err != nil {
		return nil, err
	}

	return decoded_password, nil
}

func (db *DB) GetUserArgon2Auth(username string) (*auth.Argon2Auth, error) {
	var mem int
	var time int
	var threads int
	var salt string

	err := db.Pool.QueryRow(context.Background(), "select argon2_memory, argon2_time, argon2_threads, salt from users where username=$1", username).Scan(&mem, &time, &threads, &salt)
	if err != nil {
		return nil, err
	}

	decoded_salt, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	argon2auth := &auth.Argon2Auth{
		Argon2Memory:  uint32(mem),
		Argon2Time:    uint32(time),
		Argon2Threads: uint8(threads),
		Salt:          decoded_salt,
	}

	return argon2auth, nil
}