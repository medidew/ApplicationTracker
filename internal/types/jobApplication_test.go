package types

import (
	"testing"
)

func TestJobApplicationNew(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active
	notes := []string{"test note", "test note 2"}

	job_application, err := NewJobApplication(company, role, status, notes)
	if err != nil {
		t.Fatalf("Failed to create NewJobApplication: %v", err)
	} else if job_application.company != company {
		t.Fatalf("NewJobApplication has incorrect company name")
	} else if job_application.role != role {
		t.Fatalf("NewJobApplication has incorrect role")
	} else if job_application.status != status {
		t.Fatalf("NewJobApplication has incorrect status")
	} else if len(job_application.notes) != 2 {
		t.Fatalf("NewJobApplication has too many/few notes")
	}
}

func TestJobApplicationNewFakeRole(t *testing.T) {
	company := "Medidew Inc."
	var role JobRole = "Not a real job"
	status := Active
	notes := []string{"test note", "test note 2"}

	_, err := NewJobApplication(company, role, status, notes)
	if err == nil {
		t.Fatalf("Failed to throw error on non-JobRole string")
	}
}

func TestJobApplicationNewInvalidStatus(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	var status ApplicationStatus = 255 // max uint8 val
	notes := []string{"test note", "test note 2"}

	_, err := NewJobApplication(company, role, status, notes)
	if err == nil {
		t.Fatalf("Failed to throw error on invalid status value")
	}
}

func TestUpdateStatus(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active
	notes := []string{"test note", "test note 2"}

	job_application, err := NewJobApplication(company, role, status, notes)
	if err != nil {
		t.Fatalf("Failed to create NewJobApplication: %v", err)
	}

	err = job_application.UpdateStatus(PendingResponse)
	if err != nil {
		t.Fatalf("Failed to updates NewJobApplication status: %v", err)
	}
}

func TestUpdateStatusInvalid(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active
	notes := []string{"test note", "test note 2"}

	job_application, err := NewJobApplication(company, role, status, notes)
	if err != nil {
		t.Fatalf("Failed to create NewJobApplication: %v", err)
	}

	err = job_application.UpdateStatus(255) // max uint8 val
	if err == nil {
		t.Fatalf("Failed to throw error on invalid status update")
	}
}

func TestAddNotes(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	job_application, err := NewJobApplication(company, role, status, []string{})
	if err != nil {
		t.Fatalf("Failed to create NewJobApplication: %v", err)
	}

	job_application.AddNote("test")
	job_application.AddNote("tester abc")
	num_notes := job_application.NumNotes()
	if num_notes != 2 || job_application.notes[0] != "test" || job_application.notes[1] != "tester abc" {
		t.Fatalf("Failed to add notes correctly: %v", job_application)
	}
}

func TestRemoveNotes(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	job_application, err := NewJobApplication(company, role, status, []string{})
	if err != nil {
		t.Fatalf("Failed to create NewJobApplication: %v", err)
	}

	job_application.AddNote("test")
	job_application.AddNote("tester abc")
	job_application.RemoveNote(0)
	if job_application.NumNotes() != 1 || job_application.notes[0] != "tester abc" {
		t.Fatalf("Failed to remove notes correctly: %v", job_application)
	}
}
