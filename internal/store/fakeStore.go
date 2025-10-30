package store

import (
	"errors"

	"github.com/medidew/ApplicationTracker/internal/auth"
)

type FakeStore struct {
	Applications map[string][]*JobApplication
}

func NewFakeStore(applications map[string][]*JobApplication) *FakeStore {
	return &FakeStore{
		Applications: applications,
	}
}

func (fs *FakeStore) ListApplications(username string) ([]*JobApplication, error) {
	return fs.Applications[username], nil
}

func (fs *FakeStore) GetApplication(username string, companyID string) (*JobApplication, error) {
	for _, application := range fs.Applications[username] {
		if application.GetCompany() == companyID {
			return application, nil
		}
	}

	return nil, errors.New("application not found")
}

func (fs *FakeStore) CreateApplication(username string, application *JobApplication) error {
	for _, existing_application := range fs.Applications[username] {
		if existing_application.GetCompany() == application.company {
			return errors.New("application already exists")
		}
	}

	fs.Applications[username] = append(fs.Applications[username], application)
	return nil
}

func (fs *FakeStore) DeleteApplication(username string, companyID string) error {
	for i, application := range fs.Applications[username] {
		if application.GetCompany() == companyID {
			fs.Applications[username] = append(fs.Applications[username][:i], fs.Applications[username][i+1:]...)
			return nil
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) UpdateApplicationStatus(username string, companyID string, status ApplicationStatus) error {
	for _, application := range fs.Applications[username] {
		if application.GetCompany() == companyID {
			return application.UpdateStatus(status)
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) ListApplicationNotes(username string, companyID string) ([]string, error) {
	for _, application := range fs.Applications[username] {
		if application.GetCompany() == companyID {
			return application.GetNotes(), nil
		}
	}

	return nil, errors.New("application not found")
}

func (fs *FakeStore) AddApplicationNote(username string, companyID string, note string) error {
	for _, application := range fs.Applications[username] {
		if application.GetCompany() == companyID {
			application.AddNote(note)
			return nil
		}
	}

	return errors.New("application not found")
}

func (fs *FakeStore) RemoveApplicationNote(username string, companyID string, index int) error {
	for _, application := range fs.Applications[username] {
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