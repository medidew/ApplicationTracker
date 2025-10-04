package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/medidew/ApplicationTracker/types"
)

func (a *App) ListApplications(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query(context.Background(), "select company from applications")
	if err != nil {
		http.Error(w, "DB query failed: " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	companies := []string{}

	for rows.Next() {
		var company string
		err = rows.Scan(&company)
		if err != nil {
			http.Error(w, "extracting query info failed: " + err.Error(), http.StatusInternalServerError)
			return
		}
		companies = append(companies, company)
	}

	res, err := json.Marshal(companies)
	if err != nil {
		http.Error(w, "failed to marshal: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "failed to write response: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) GetApplication(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")

	var role types.JobRole
	var status types.ApplicationStatus
	var notes []string
	err := a.DB.QueryRow(context.Background(), "select role, status, notes from applications where company=$1", companyID).Scan(&role, &status, &notes)
	if err != nil {
		http.Error(w, "DB query failed: " + err.Error(), http.StatusInternalServerError)
		return
	}

	ja, err := types.NewJobApplication(companyID, role, status)
	if err != nil {
		http.Error(w, "JobApplication constructor error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(notes); i++ {
		ja.AddNote(notes[i])
	}

	res, err := json.Marshal(ja)
	if err != nil {
		http.Error(w, "failed to marshal: " + err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "failed to write reponse: " + err.Error(), http.StatusInternalServerError)
		return
	}
	
}