package store

import (
	"errors"
	"fmt"

	"github.com/medidew/ApplicationTracker/internal/auth"
)

type FakeStore struct {
	Applications []*JobApplication
}

func NewFakeStore(applications []*JobApplication) *FakeStore {
	return &FakeStore{
		Applications: applications,
	}
}

func (fs *FakeStore) ListApplications() ([]*JobApplication, error) {
	return fs.Applications, nil
}

func (fs *FakeStore) GetApplication(companyID string) (*JobApplication, error) {
	for _, application := range fs.Applications {
		if application.GetCompany() == companyID {
			return application, nil
		}
	}

	return nil, errors.New("application not found")
}

func (fs *FakeStore) CreateApplication(application *JobApplication) error {
	for _, existing_application := range fs.Applications {
		if existing_application.GetCompany() == application.company {
			return errors.New("application already exists")
		}
	}

	fs.Applications = append(fs.Applications, application)
	return nil
}

func (fs *FakeStore) DeleteApplication(companyID string) error {
	for i, application := range fs.Applications {
		if application.GetCompany() == companyID {
			fs.Applications = append(fs.Applications[:i], fs.Applications[i+1:]...)
			return nil
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) UpdateApplicationStatus(companyID string, status ApplicationStatus) error {
	fmt.Printf("companyID: %v\n", companyID)
	fmt.Printf("status: %v\n", status)
	for _, application := range fs.Applications {
		fmt.Printf("application: %v\n", application)
		if application.GetCompany() == companyID {
			return application.UpdateStatus(status)
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) ListApplicationNotes(companyID string) ([]string, error) {
	for _, application := range fs.Applications {
		if application.GetCompany() == companyID {
			return application.GetNotes(), nil
		}
	}

	return nil, errors.New("application not found")
}

func (fs *FakeStore) AddApplicationNote(companyID string, note string) error {
	for _, application := range fs.Applications {
		if application.GetCompany() == companyID {
			application.AddNote(note)
			return nil
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) RemoveApplicationNote(companyID string, index int) error {
	for _, application := range fs.Applications {
		if application.GetCompany() == companyID {
			return application.RemoveNote(index)
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) CreateUser(email string, username string, argon2auth *auth.Argon2Auth, hashedPassword []byte) error {
	// FakeStore does not yet implement user storage.
	return nil
}

func (fs *FakeStore) GetUserHashedPassword(username string) ([]byte, error) {
	// FakeStore does not yet implement user storage.
	return nil, nil
}

func (fs *FakeStore) GetUserArgon2Auth(username string) (*auth.Argon2Auth, error) {
	// FakeStore does not yet implement user storage.
	return nil, nil
}