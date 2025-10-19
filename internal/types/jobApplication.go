package types

import (
	"encoding/json"
	"errors"
	"slices"
)

// Type for standardising application status values.
type ApplicationStatus byte

const (
	Active          ApplicationStatus = iota // Pending action from applicant.
	PendingResponse                          // Pending response from employer.
	Rejected
	Offer
)

const MaxStatus ApplicationStatus = Offer

func (jr ApplicationStatus) String() string {
	switch jr {
	case Active:
		return "Active"
	case PendingResponse:
		return "Pending Response"
	case Rejected:
		return "Rejected"
	case Offer:
		return "Offer"
	default:
		return "err"
	}
}

// Type for standardising jobs title strings.
type JobRole string

const (
	SoftwareEngineer JobRole = "Software Engineer"
)

func GetSupportedJobRoles() []JobRole {
	return []JobRole{SoftwareEngineer}
}

// Job application details.
type JobApplication struct {
	company string
	role    JobRole
	status  ApplicationStatus
	notes   []string
}

func (job_application *JobApplication) MarshalJSON() ([]byte, error) {
	company_json, err := json.Marshal(job_application.company)
	if err != nil {
		return nil, errors.Join(errors.New("could not marshal company"), err)
	}
	role_json, err := json.Marshal(job_application.role)
	if err != nil {
		return nil, errors.Join(errors.New("could not marshal role"), err)
	}
	status_json, err := json.Marshal(job_application.status.String())
	if err != nil {
		return nil, errors.Join(errors.New("could not marshal status"), err)
	}
	notes_json, err := json.Marshal(job_application.notes)
	if err != nil {
		return nil, errors.Join(errors.New("could not marshal notes"), err)
	}

	// I could construct this as a string then convert afterwards to make it cleaner,
	//	but this version is impervious to whether I change JobApplication's field types,
	// 	and also makes adding new fields comically easy.
	// 	The impact is minimal anyway so I don't care.
	result := []byte(`{"company":`)
	result = append(result, company_json...)
	result = append(result, []byte(`, "role":`)...)
	result = append(result, role_json...)
	result = append(result, []byte(`, "status":`)...)
	result = append(result, status_json...)
	result = append(result, []byte(`, "notes":`)...)
	result = append(result, notes_json...)
	result = append(result, []byte(`}`)...)

	return result, nil
}

func (JobApplication *JobApplication) UnmarshalJSON(data []byte) error {
	aux := struct {
		Company string            `json:"company"`
		Role    JobRole           `json:"role"`
		Status  ApplicationStatus `json:"status"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if !slices.Contains(GetSupportedJobRoles(), aux.Role) {
		return errors.New("`role` is not supported by type JobRole")
	}

	if aux.Status > MaxStatus {
		return errors.New("`status` is not supported by type ApplicationStatus")
	}

	JobApplication.company = aux.Company
	JobApplication.role = aux.Role
	JobApplication.status = aux.Status

	return nil
}

func NewJobApplication(company string, role JobRole, status ApplicationStatus, notes []string) (*JobApplication, error) {
	if !slices.Contains(GetSupportedJobRoles(), role) {
		return nil, errors.New("`role` is not supported by type JobRole")
	}

	if status > MaxStatus {
		return nil, errors.New("`status` is not supported by type ApplicationStatus")
	}

	return &JobApplication{
		company: company,
		role:    role,
		status:  status,
		notes:   notes,
	}, nil
}

func (job_application *JobApplication) GetCompany() string {
	return job_application.company
}

func (job_application *JobApplication) GetRole() JobRole {
	return job_application.role
}

func (job_application *JobApplication) GetStatus() ApplicationStatus {
	return job_application.status
}

func (job_application *JobApplication) UpdateStatus(status ApplicationStatus) error {
	if status > MaxStatus {
		return errors.New("`status` is not supported by type ApplicationStatus")
	}

	job_application.status = status

	return nil
}

func (job_application *JobApplication) GetNotes() []string {
	return job_application.notes
}

func (job_application *JobApplication) AddNote(note string) {
	job_application.notes = append(job_application.notes, note)
}

func (job_application *JobApplication) NumNotes() int {
	return len(job_application.notes)
}

func (job_application *JobApplication) RemoveNote(index int) error {
	if index >= len(job_application.notes) {
		return errors.New("index out of range")
	} else if index < 0 {
		return errors.New("index negative")
	}

	job_application.notes = slices.Delete(job_application.notes, index, index+1)
	return nil
}

func (job_application *JobApplication) String() string {
	notes_string := "["

	if len(job_application.notes) > 1 {
		notes_string += "\"" + job_application.notes[0] + "\""

		for i := 1; i < len(job_application.notes); i++ {
			notes_string += " \"" + job_application.notes[i] + "\""
		}
	} else if len(job_application.notes) == 1 {
		notes_string += "\"" + job_application.notes[0] + "\""
	}
	notes_string += "]"

	return "&{company: " + job_application.company + ", role: " + string(job_application.role) + ", status: " + job_application.status.String() + ", notes: " + notes_string + "}"
}
