package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/medidew/ApplicationTracker/types"
)

func (app *App) ListApplications(response_writer http.ResponseWriter, request *http.Request) {
	rows, err := app.DB.Query(context.Background(), "select company from applications")
	if err != nil {
		http.Error(response_writer, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	companies := []string{}

	for rows.Next() {
		var company string
		err = rows.Scan(&company)
		if err != nil {
			http.Error(response_writer, "extracting query info failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		companies = append(companies, company)
	}

	response, err := json.Marshal(companies)
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

	var role types.JobRole
	var status types.ApplicationStatus
	var notes []string
	err := app.DB.QueryRow(context.Background(), "select role, status, notes from applications where company=$1", companyID).Scan(&role, &status, &notes)
	if err != nil {
		http.Error(response_writer, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	job_application, err := types.NewJobApplication(companyID, role, status)
	if err != nil {
		http.Error(response_writer, "JobApplication constructor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(notes); i++ {
		job_application.AddNote(notes[i])
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
