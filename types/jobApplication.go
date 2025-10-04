package types

import (
	"encoding/json"
	"errors"
	"slices"
)

type ApplicationStatus byte

// Type for standardising jobs title strings.
type JobRole string

const (
	Active ApplicationStatus = iota // Pending action from applicant.
	PendingResponse // Pending response from employer.
	Rejected
	Offer
)
var maxStatus ApplicationStatus = Offer

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

func (ja *JobApplication) MarshalJSON() ([]byte, error) {
	company_json, err := json.Marshal(ja.company)
	if err != nil { return nil, errors.Join(errors.New("could not marshal company"), err) }
	role_json, err := json.Marshal(ja.role)
	if err != nil { return nil, errors.Join(errors.New("could not marshal role"), err) }
	status_json, err := json.Marshal(ja.status.String())
	if err != nil { return nil, errors.Join(errors.New("could not marshal status"), err) }
	notes_json, err := json.Marshal(ja.notes)
	if err != nil { return nil, errors.Join(errors.New("could not marshal notes"), err) }

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

func NewJobApplication(company string, role JobRole, status ApplicationStatus) (*JobApplication, error) {
	if !slices.Contains(roles, role) {
		return nil, errors.New("`role` is not supported by type JobRole")
	}

	if status > maxStatus {
		return nil, errors.New("`status` is not supported by type ApplicationStatus")
	}

	return &JobApplication{
		company: company,
		role:    role,
		status:  status,
		notes:   []string{},
	}, nil
}

func(ja *JobApplication) GetCompany() string {
	return ja.company
}

func(ja *JobApplication) GetRole() JobRole {
	return ja.role
}

func(ja *JobApplication) GetStatus() ApplicationStatus {
	return ja.status
}

func(ja *JobApplication) UpdateStatus(status ApplicationStatus) error {
	if status > maxStatus {
		return errors.New("`status` is not supported by type ApplicationStatus")
	}

	ja.status = status

	return nil
}

func (ja *JobApplication) GetNotes() []string {
	return ja.notes
}

func(ja *JobApplication) AddNote(note string) {
	ja.notes = append(ja.notes, note)
}

func(ja *JobApplication) NumNotes() int {
	return len(ja.notes)
}

func(ja *JobApplication) RemoveNote(index int) error {
	if index >= len(ja.notes) {
		return errors.New("index out of range")
	} else if index < 0 {
		return errors.New("index negative")
	}

	ja.notes = slices.Delete(ja.notes, index, index + 1)
	return nil
}

func (ja *JobApplication) String() string {
	notes_string := "["

	if len(ja.notes) > 1 {
		notes_string += "\"" + ja.notes[0] + "\""

		for i := 1; i < len(ja.notes); i++ {
			notes_string += " \"" + ja.notes[i] + "\""
		}
	} else if len(ja.notes) == 1 {
		notes_string += "\"" + ja.notes[0] + "\""
	}
	notes_string += "]"

	return "&{company: " + ja.company + ", role: " + string(ja.role) + ", status: " + ja.status.String() + ", notes: " + notes_string + "}"
}