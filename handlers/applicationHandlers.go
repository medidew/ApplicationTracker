package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

func (a *App) ListApplications(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query(context.Background(), "select company from applications")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	companies := []string{}

	for rows.Next() {
		var company string
		err = rows.Scan(&company)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		companies = append(companies, company)
	}

	res, err := json.Marshal(companies)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(res)
}