package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/medidew/ApplicationTracker/internal/types"
)

func (app *App) ListApplications(response_writer http.ResponseWriter, request *http.Request) {
	applications, err := app.DB.ListApplications()
	if err != nil {
		http.Error(response_writer, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(applications)
	if err != nil {
		http.Error(response_writer, "failed to marshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusOK)
	_, err = response_writer.Write(response)
	if err != nil {
		http.Error(response_writer, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) GetApplication(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")

	job_application, err := app.DB.GetApplication(companyID)
	if err != nil {
		http.Error(response_writer, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(job_application)
	if err != nil {
		http.Error(response_writer, "failed to marshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = response_writer.Write(response)
	if err != nil {
		http.Error(response_writer, "failed to write reponse: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *App) CreateApplication(response_writer http.ResponseWriter, request *http.Request) {
	request_body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response_writer, "failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	new_application := &types.JobApplication{}

	err = new_application.UnmarshalJSON(request_body)
	if err != nil {
		http.Error(response_writer, "failed to unmarshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = app.DB.CreateApplication(new_application)
	if err != nil {
		http.Error(response_writer, "DB insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusCreated)
}

func (app *App) DeleteApplication(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")

	err := app.DB.DeleteApplication(companyID)
	if err != nil {
		http.Error(response_writer, "DB delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusNoContent)
}

func (app *App) UpdateApplicationStatus(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")

	var status_update struct {
		Status types.ApplicationStatus `json:"status"`
	}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&status_update)
	if err != nil {
		http.Error(response_writer, "failed to unmarshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = app.DB.UpdateApplicationStatus(companyID, status_update.Status)
	if err != nil {
		http.Error(response_writer, "DB update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusNoContent)
}

func (app *App) ListApplicationNotes(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")

	notes, err := app.DB.ListApplicationNotes(companyID)
	if err != nil {
		http.Error(response_writer, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(notes)
	if err != nil {
		http.Error(response_writer, "failed to marshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = response_writer.Write(response)
	if err != nil {
		http.Error(response_writer, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) AddApplicationNote(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")

	var note_addition struct {
		Note string `json:"note"`
	}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	
	err := decoder.Decode(&note_addition)
	if err != nil {
		http.Error(response_writer, "failed to unmarshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = app.DB.AddApplicationNote(companyID, note_addition.Note)
	if err != nil {
		http.Error(response_writer, "DB update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusCreated)
}

func (app *App) RemoveApplicationNote(response_writer http.ResponseWriter, request *http.Request) {
	companyID := chi.URLParam(request, "companyID")
	noteIndex, err := strconv.Atoi(chi.URLParam(request, "noteIndex"))
	if err != nil {
		http.Error(response_writer, "invalid note index: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = app.DB.RemoveApplicationNote(companyID, noteIndex)
	if err != nil {
		http.Error(response_writer, "DB update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusNoContent)
}