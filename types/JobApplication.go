package types

import (
	"errors"
	"slices"
)

type ApplicationStatus byte

// Type for standardising jobs title strings.
type JobRole string

const (
	Active          ApplicationStatus = 0 // pending action from applicant
	PendingResponse ApplicationStatus = 1 // pending response from employer
	Rejected        ApplicationStatus = 2
	Offer           ApplicationStatus = 3
)
var maxStatus ApplicationStatus = Offer

const (
	SoftwareEngineer JobRole = "Software Engineer"
)
var roles []JobRole = []JobRole{SoftwareEngineer}

// Job application details.
type JobApplication struct {
	company string
	role    JobRole
	status  ApplicationStatus
	notes   []string
}

func NewJobApplication(company string, role JobRole, status ApplicationStatus) (*JobApplication, error) {
	if !slices.Contains(roles, role) {
		return nil, errors.New("role is not supported by type JobRole")
	}

	if status > maxStatus {
		return nil, errors.New("status is not supported by type ApplicationStatus (are you using the given constants?)")
	}

	return &JobApplication{
		company: company,
		role:    role,
		status:  status,
		notes:   []string{},
	}, nil
}

func(ja *JobApplication) AddNote(note string) {
	ja.notes = append(ja.notes, note)
}

func(ja *JobApplication) UpdateStatus(status ApplicationStatus) error {
	if status > maxStatus {
		return errors.New("status is not supported by type ApplicationStatus (are you using the given constants?)")
	}

	ja.status = status

	return nil
}

func(ja *JobApplication) NumNotes() int {
	return len(ja.notes)
}