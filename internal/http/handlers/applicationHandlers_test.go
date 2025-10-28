package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/medidew/ApplicationTracker/internal/store"
)

func setupTestRouter(*App) http.Handler {
	router := SetupRouter(setupTestApp())
	return router
}

func setupTestApp() *App {
	fake_application_one, err := store.NewJobApplication("Fake Company", store.SoftwareEngineer, store.Active, []string{"Note one.", "Note two."})
	if err != nil {
		panic("Failed to create fake job application for testing: " + err.Error())
	}
	
	fake_application_two, err := store.NewJobApplication("Another Fake Company", store.SoftwareEngineer, store.PendingResponse, []string{"Initial note."})
	if err != nil {
		panic("Failed to create fake job application for testing: " + err.Error())
	}

	db := store.NewFakeStore([]*store.JobApplication{fake_application_one, fake_application_two})

	return &App{
		DB:     db,
		Logger: zap.NewNop(),
	}
}

func TestListApplications(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodGet, "/applications", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	app.ListApplications(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestGetApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	// Test fetching an existing application
	request := httptest.NewRequest(http.MethodGet, "/applications/Fake%20Company", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestGetInvalidApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	// Test fetching a non-existing application
	request := httptest.NewRequest(http.MethodGet, "/applications/Non%20Existent%20Company", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestCreateApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	new_application := `{
		"company": "New Company",
		"role": "Software Engineer",
		"status": 0,
		"notes": ["Exciting opportunity."]
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications",  io.NopCloser(strings.NewReader(new_application)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, response.StatusCode)
	}
}

func TestCreateDuplicateApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	duplicate_application := `{
		"company": "Fake Company",
		"role": "Software Engineer",
		"status": 0,
		"notes": ["Another note."]
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications",  io.NopCloser(strings.NewReader(duplicate_application)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)
	
	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestDeleteApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d", http.StatusNoContent, response.StatusCode)
	}
}

func TestDeleteInvalidApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Non%20Existent%20Company", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestUpdateApplicationStatus(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	status_update := `{
		"status": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Fake%20Company", io.NopCloser(strings.NewReader(status_update)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d", http.StatusNoContent, response.StatusCode)
	}
}

func TestUpdateInvalidApplicationStatus(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	status_update := `{
		"status": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Non%20Existent%20Company", io.NopCloser(strings.NewReader(status_update)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestUpdateApplicationStatusInvalidPayload(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	invalid_status_update := `{
		"state": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Fake%20Company", io.NopCloser(strings.NewReader(invalid_status_update)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, response.StatusCode)
	}
}

func TestListApplicationNotes(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodGet, "/applications/Fake%20Company/notes", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestAddApplicationNote(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	note_addition := `{
		"note": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Fake%20Company/notes", io.NopCloser(strings.NewReader(note_addition)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, response.StatusCode)
	}
}

func TestAddApplicationNoteInvalidApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	note_addition := `{
		"note": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Non%20Existent%20Company/notes", io.NopCloser(strings.NewReader(note_addition)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestAddApplicationNoteInvalidPayload(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	invalid_note_addition := `{
		"text": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Fake%20Company/notes", io.NopCloser(strings.NewReader(invalid_note_addition)))
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, response.StatusCode)
	}
}

func TestRemoveApplicationNote(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/0", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d", http.StatusNoContent, response.StatusCode)
	}
}

func TestRemoveApplicationNoteInvalidApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Non%20Existent%20Company/notes/0", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestRemoveApplicationNoteOutOfRange(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/10", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}

func TestRemoveApplicationNoteNegativeIndex(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/-1", nil)
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}