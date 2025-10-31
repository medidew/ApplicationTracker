package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/alexedwards/scs/v2"
	"github.com/medidew/ApplicationTracker/internal/store"
)

func setupTestRouter(app *App) http.Handler {
	router := SetupRouter(app)
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

	database_map := map[string][]*store.JobApplication{
		"testuser": {fake_application_one, fake_application_two},
	}
	database := store.NewFakeStore(database_map)

	session_manager := scs.New()
	session_manager.Cookie.Secure = false

	return &App{
		DB:     database,
		Logger: zap.NewNop(),
		SessionManager: session_manager,
	}
}

func setupSessionContext(app *App, username string) (string, error) {
	ctx, _ := app.SessionManager.Load(context.Background(), "")
	app.SessionManager.Put(ctx, "username", username)
	token, _, err := app.SessionManager.Commit(ctx)
	if err != nil {
		return "", err
	}

	return token, nil
}

func setupAll() (*App, http.Handler, string, error) {
	app := setupTestApp()
	router := setupTestRouter(app)
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		return nil, nil, "", err
	}
	return app, router, token, nil
}

func TestListApplications(t *testing.T) {
	app, router, token, err := setupAll()
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodGet, "/applications", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
	response_recorder := httptest.NewRecorder()

	fmt.Printf("token: %v\n", token)

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestGetApplication(t *testing.T) {
	app := setupTestApp()
	router := setupTestRouter(app)
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	// Test fetching an existing application
	request := httptest.NewRequest(http.MethodGet, "/applications/Fake%20Company", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	// Test fetching a non-existing application
	request := httptest.NewRequest(http.MethodGet, "/applications/Non%20Existent%20Company", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	new_application := `{
		"company": "New Company",
		"role": "Software Engineer",
		"status": 0,
		"notes": ["Exciting opportunity."]
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications",  io.NopCloser(strings.NewReader(new_application)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	duplicate_application := `{
		"company": "Fake Company",
		"role": "Software Engineer",
		"status": 0,
		"notes": ["Another note."]
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications",  io.NopCloser(strings.NewReader(duplicate_application)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Non%20Existent%20Company", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	status_update := `{
		"status": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Fake%20Company", io.NopCloser(strings.NewReader(status_update)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	status_update := `{
		"status": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Non%20Existent%20Company", io.NopCloser(strings.NewReader(status_update)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	invalid_status_update := `{
		"state": 2
	}`

	request := httptest.NewRequest(http.MethodPut, "/applications/Fake%20Company", io.NopCloser(strings.NewReader(invalid_status_update)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodGet, "/applications/Fake%20Company/notes", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	note_addition := `{
		"note": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Fake%20Company/notes", io.NopCloser(strings.NewReader(note_addition)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	note_addition := `{
		"note": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Non%20Existent%20Company/notes", io.NopCloser(strings.NewReader(note_addition)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	invalid_note_addition := `{
		"text": "This is a new note."
	}`

	request := httptest.NewRequest(http.MethodPost, "/applications/Fake%20Company/notes", io.NopCloser(strings.NewReader(invalid_note_addition)))
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/0", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Non%20Existent%20Company/notes/0", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/10", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
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
	token, err := setupSessionContext(app, "testuser")
	if err != nil {
		t.Fatalf("Failed to setup session:%v", err.Error())
	}

	request := httptest.NewRequest(http.MethodDelete, "/applications/Fake%20Company/notes/-1", nil)
	request.AddCookie(&http.Cookie{Name: app.SessionManager.Cookie.Name, Value: token})
	response_recorder := httptest.NewRecorder()

	router.ServeHTTP(response_recorder, request)

	response := response_recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}
}