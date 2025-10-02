package types
import (
	"testing"
)

func TestJobApplicationNew(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	ja, err := NewJobApplication(company, role, status)
	if err != nil {
		t.Errorf("Failed to create NewJobApplication: %v", err)
	} else if ja.company != company {
		t.Errorf("NewJobApplication has incorrect company name")
	} else if ja.role != role {
		t.Errorf("NewJobApplication has incorrect role")
	} else if ja.status != status {
		t.Errorf("NewJobApplication has incorrect status")
	} else if len(ja.notes) > 0 {
		t.Errorf("NewJobApplication has >0 notes")
	}
}

func TestJobApplicationNewFakeRole(t *testing.T) {
	company := "Medidew Inc."
	var role JobRole = "Not a real job"
	status := Active

	_, err := NewJobApplication(company, role, status)
	if err == nil {
		t.Errorf("Failed to throw error on non-JobRole string")
	}
}

func TestJobApplicationNewInvalidStatus(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	var status ApplicationStatus = 255 // max uint8 val

	_, err := NewJobApplication(company, role, status)
	if err == nil {
		t.Errorf("Failed to throw error on invalid status value")
	}
}

func TestUpdateStatus(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	ja, err := NewJobApplication(company, role, status)
	if err != nil {
		t.Errorf("Failed to create NewJobApplication: %v", err)
	}

	err = ja.UpdateStatus(PendingResponse)
	if err != nil {
		t.Errorf("Failed to updates NewJobApplication status: %v", err)
	}
}

func TestUpdateStatusInvalid(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	ja, err := NewJobApplication(company, role, status)
	if err != nil {
		t.Errorf("Failed to create NewJobApplication: %v", err)
	}

	err = ja.UpdateStatus(255) // max uint8 val
	if err == nil {
		t.Errorf("Failed to throw error on invalid status update")
	}
}

func TestNumNotes(t *testing.T) {
	company := "Medidew Inc."
	role := SoftwareEngineer
	status := Active

	ja, err := NewJobApplication(company, role, status)
	if err != nil {
		t.Errorf("Failed to create NewJobApplication: %v", err)
	}

	ja.AddNote("test")
	ja.AddNote("tester \n abc")
	numNotes := ja.NumNotes()
	if numNotes != 2 {
		t.Errorf("Failed to count notes correctly: %v instead of 2", numNotes)
	}
}